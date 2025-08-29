package iot

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nhan1603/IoTsystem/api/internal/model"
	"github.com/nhan1603/IoTsystem/api/internal/pkg/kafka"
	"github.com/nhan1603/IoTsystem/api/internal/pkg/obsmetrics"
	"github.com/nhan1603/IoTsystem/api/internal/repository"
)

// BenchmarkMetrics tracks performance metrics
type BenchmarkMetrics struct {
	TotalRecords     int64
	ProcessedRecords int64
	FailedRecords    int64
	StartTime        time.Time
	LastProcessed    time.Time
	BatchCount       int64
	TotalLatency     time.Duration
}

// Stop gracefully stops the consumer
func (c *impl) Stop() {
	log.Println("Stopping IoT consumer...")
	close(c.stopChan)
	c.wg.Wait()
	log.Println("IoT consumer stopped")
}

// handleBatch processes a batch of messages from Kafka
func (c *impl) HandleBatch(ctx context.Context, msgs []kafka.ConsumerMessage) error {
	var readings []model.SensorReading
	start := time.Now()
	obsmetrics.BatchSize.Set(float64(len(msgs)))

	for _, msg := range msgs {
		var iotMsg model.IoTDataMessage
		if err := json.Unmarshal(msg.Value, &iotMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			c.incrementFailedRecords()
			continue
		}
		readings = append(readings, model.SensorReading{
			DeviceID:    iotMsg.DeviceID,
			DeviceName:  iotMsg.DeviceName,
			DeviceType:  iotMsg.DeviceType,
			Location:    iotMsg.Location,
			Floor:       iotMsg.Floor,
			Zone:        iotMsg.Zone,
			Temperature: iotMsg.Temperature,
			Humidity:    iotMsg.Humidity,
			CO2:         iotMsg.CO2,
			Timestamp:   iotMsg.Timestamp,
			CreatedAt:   start,
		})

		// Record sensor readings
		obsmetrics.SensorReadings.WithLabelValues(iotMsg.DeviceID, "temperature").Set(iotMsg.Temperature)
		obsmetrics.SensorReadings.WithLabelValues(iotMsg.DeviceID, "humidity").Set(iotMsg.Humidity)
		obsmetrics.SensorReadings.WithLabelValues(iotMsg.DeviceID, "co2").Set(iotMsg.CO2)

		obsmetrics.MessagesProcessed.WithLabelValues(iotMsg.DeviceID, "success").Inc()
	}

	if len(readings) == 0 {
		log.Println("No valid readings to process")
		return nil
	}

	// Insert readings using transaction
	err := c.repo.DoInTx(ctx, func(txRepo repository.Registry) error {
		return txRepo.IoT().BatchInsertReadings(ctx, readings)
	})

	if err != nil {
		log.Printf("Error processing batch: %v", err)
		return fmt.Errorf("failed to process batch: %w", err)
	}

	// Record metrics only after commit
	latency := time.Since(start)
	c.updateMetrics(len(readings), latency)
	c.SaveMetrics(ctx)

	obsmetrics.ProcessingLatency.WithLabelValues("batch").Observe(time.Since(start).Seconds())

	log.Printf("[HandleBatch] Processed batch of %d records in %v", len(readings), latency)
	return nil
}

// updateMetrics updates performance metrics
func (c *impl) updateMetrics(recordCount int, latency time.Duration) {
	c.metricsMutex.Lock()
	defer c.metricsMutex.Unlock()

	c.metrics.ProcessedRecords += int64(recordCount)
	c.metrics.TotalRecords += int64(recordCount)
	c.metrics.BatchCount++
	c.metrics.TotalLatency += latency
	c.metrics.LastProcessed = time.Now()
}

// incrementFailedRecords increments the failed records counter
func (c *impl) incrementFailedRecords() {
	c.metricsMutex.Lock()
	defer c.metricsMutex.Unlock()
	c.metrics.FailedRecords++
}

// addFailedRecords increments the failed records counter
func (c *impl) addFailedRecords(recordsCount int) {
	c.metricsMutex.Lock()
	defer c.metricsMutex.Unlock()
	c.metrics.FailedRecords += int64(recordsCount)
}

// GetMetrics returns current benchmark metrics
func (c *impl) GetMetrics() model.BenchmarkMetrics {
	c.metricsMutex.RLock()
	defer c.metricsMutex.RUnlock()

	var avgLatency float64
	if c.metrics.BatchCount > 0 {
		avgLatency = float64(c.metrics.TotalLatency.Milliseconds()) / float64(c.metrics.BatchCount)
	}

	var throughput float64
	if !c.metrics.StartTime.IsZero() {
		duration := time.Since(c.metrics.StartTime).Seconds()
		if duration > 0 {
			throughput = float64(c.metrics.ProcessedRecords) / duration
		}
	}

	return model.BenchmarkMetrics{
		TotalRecords:     c.metrics.TotalRecords,
		ProcessedRecords: c.metrics.ProcessedRecords,
		FailedRecords:    c.metrics.FailedRecords,
		StartTime:        c.metrics.StartTime,
		EndTime:          c.metrics.LastProcessed,
		AverageLatency:   avgLatency,
		Throughput:       throughput,
		BatchSize:        c.batchSize,
		DatabaseType:     "PostgreSQL",
	}
}

// SaveMetrics saves the current metrics to the database
func (c *impl) SaveMetrics(ctx context.Context) error {
	metrics := c.GetMetrics()
	// validate metrics before passing to repository
	if metrics.TotalRecords <= 0 || metrics.ProcessedRecords < 0 || metrics.FailedRecords < 0 {
		return fmt.Errorf("invalid benchmark metrics data")
	}
	if metrics.StartTime.IsZero() || metrics.EndTime.IsZero() {
		return fmt.Errorf("start and end time must be provided")
	}
	if metrics.EndTime.Before(metrics.StartTime) {
		return fmt.Errorf("end time cannot be before start time")
	}
	if metrics.AverageLatency < 0 || metrics.Throughput < 0 || metrics.BatchSize <= 0 {
		return fmt.Errorf("invalid latency, throughput or batch size")
	}
	// repository call to save metrics
	return c.repo.IoT().SaveBenchmarkMetrics(ctx, metrics)
}
