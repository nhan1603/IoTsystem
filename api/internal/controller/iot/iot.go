package iot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/nhan1603/IoTsystem/api/internal/model"
	"github.com/nhan1603/IoTsystem/api/internal/pkg/env"
	"github.com/nhan1603/IoTsystem/api/internal/pkg/kafka"
	"github.com/nhan1603/IoTsystem/api/internal/repository"
)

// Controller represents the specification of this pkg
type Controller interface {
	GetDevices(ctx context.Context) ([]model.IoTDevice, error)
	GetBenchmarkMetrics(ctx context.Context, limit int) ([]model.BenchmarkMetrics, error)
	GetLatestReadings(ctx context.Context) ([]model.SensorReading, error)
	GetMetrics() model.BenchmarkMetrics
	GetReadings(ctx context.Context, input model.GetReadingsInput) ([]model.SensorReading, error)
	GetReadingsByDevice(ctx context.Context, deviceID string, limit int) ([]model.SensorReading, error)
	GetReadingsByTimeRange(ctx context.Context, startTime time.Time, endTime time.Time, limit int) ([]model.SensorReading, error)
	HandleBatch(ctx context.Context, msgs []kafka.ConsumerMessage) error
	SaveBenchmarkMetrics(ctx context.Context, metrics model.BenchmarkMetrics) error
	SaveMetrics(ctx context.Context) error
}

// impl handles IoT data operations
type impl struct {
	repo         repository.Registry
	batchSize    int
	metrics      *BenchmarkMetrics
	metricsMutex sync.RWMutex
	stopChan     chan struct{}
	wg           sync.WaitGroup
}

// New creates a new IoT controller
func New(repo repository.Registry) *impl {
	batchSize, _ := strconv.Atoi(env.GetwithDefault("BATCH_SIZE", "100"))
	return &impl{
		repo:      repo,
		batchSize: batchSize,
		metrics: &BenchmarkMetrics{
			StartTime: time.Now(),
		},
		stopChan: make(chan struct{}),
	}
}

// GetDevices retrieves all IoT devices
func (c *impl) GetDevices(ctx context.Context) ([]model.IoTDevice, error) {
	devices, err := c.repo.IoT().GetDevices(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get devices: %w", err)
	}
	return devices, nil
}

// GetReadings retrieves sensor readings with filters
func (c *impl) GetReadings(ctx context.Context, input model.GetReadingsInput) ([]model.SensorReading, error) {
	readings, err := c.repo.IoT().GetReadings(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get readings: %w", err)
	}
	return readings, nil
}

// GetLatestReadings retrieves the latest reading for each device
func (c *impl) GetLatestReadings(ctx context.Context) ([]model.SensorReading, error) {
	readings, err := c.repo.IoT().GetLatestReadings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get latest readings: %w", err)
	}
	return readings, nil
}

// GetReadingsByDevice retrieves readings for a specific device
func (c *impl) GetReadingsByDevice(ctx context.Context, deviceID string, limit int) ([]model.SensorReading, error) {
	input := model.GetReadingsInput{
		DeviceID: deviceID,
		Limit:    limit,
	}
	readings, err := c.repo.IoT().GetReadings(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get readings for device %s: %w", deviceID, err)
	}
	return readings, nil
}

// GetReadingsByTimeRange retrieves readings within a time range
func (c *impl) GetReadingsByTimeRange(ctx context.Context, startTime, endTime time.Time, limit int) ([]model.SensorReading, error) {
	input := model.GetReadingsInput{
		StartTime: startTime,
		EndTime:   endTime,
		Limit:     limit,
	}
	readings, err := c.repo.IoT().GetReadings(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to get readings by time range: %w", err)
	}
	return readings, nil
}

// GetBenchmarkMetrics retrieves benchmark performance metrics
func (c *impl) GetBenchmarkMetrics(ctx context.Context, limit int) ([]model.BenchmarkMetrics, error) {
	metrics, err := c.repo.IoT().GetBenchmarkMetrics(ctx, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get benchmark metrics: %w", err)
	}
	return metrics, nil
}

// SaveBenchmarkMetrics saves benchmark performance metrics
func (c *impl) SaveBenchmarkMetrics(ctx context.Context, metrics model.BenchmarkMetrics) error {
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
	err := c.repo.IoT().SaveBenchmarkMetrics(ctx, metrics)
	if err != nil {
		return fmt.Errorf("failed to save benchmark metrics: %w", err)
	}
	log.Printf("Saved benchmark metrics: %+v", metrics)
	return nil
}
