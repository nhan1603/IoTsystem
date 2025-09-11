package cassiot

import (
	"context"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/nhan1603/IoTsystem/api/internal/model"
	"github.com/nhan1603/IoTsystem/api/internal/repository/iotsystem"
	"github.com/scylladb/gocqlx/v2"
	"github.com/scylladb/gocqlx/v2/qb"
)

type cassandraImpl struct {
	session   *gocqlx.Session
	tsDeduper *TsDeduper
}

func NewCassandra(session gocqlx.Session) iotsystem.Repository {
	return &cassandraImpl{
		session:   &session,
		tsDeduper: NewTsDeduper(),
	}
}

func (c *cassandraImpl) GetDevices(ctx context.Context) ([]model.IoTDevice, error) {
	// Add context timeout
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Only select active devices and add pagination
	stmt, names := qb.Select("iot_devices").
		Columns("id", "device_id", "name", "type", "location", "floor_id", "zone_id",
						"is_active", "created_at", "updated_at").
		Where(qb.Eq("is_active")). // Only get active devices
		AllowFiltering().          // Required for non-partition key filters
		ToCql()

	q := c.session.Query(stmt, names).
		BindMap(qb.M{"is_active": true}).
		WithContext(ctx)
	defer q.Release()

	// Use scanner for better memory efficiency
	scanner := q.Iter().Scanner()
	var devices []model.IoTDevice

	for scanner.Next() {
		var device model.IoTDevice
		var id gocql.UUID // Use UUID instead of int64

		err := scanner.Scan(&id,
			&device.DeviceID,
			&device.Name,
			&device.Type,
			&device.Location,
			&device.Floor,
			&device.Zone,
			&device.IsActive,
			&device.CreatedAt,
			&device.UpdatedAt)

		if err != nil {
			return nil, fmt.Errorf("failed to scan device: %w", err)
		}
		devices = append(devices, device)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error during scan: %w", err)
	}

	return devices, nil
}

func (c *cassandraImpl) BatchInsertReadings(ctx context.Context, readings []model.SensorReading) error {
	batch := c.session.Session.NewBatch(gocql.UnloggedBatch)
	batch.SetConsistency(gocql.One)
	stmt, _ := qb.Insert("sensor_readings").
		Columns("id", "device_id", "device_name", "device_type", "location",
			"floor_id", "zone_id", "temperature", "humidity", "co2", "timestamp",
			"created_at", "heat_index", "air_quality_index", "durable_write_ts").
		ToCql()

	for _, reading := range readings {

		// Calculate heat index
		heatIndex := 0.5 * (reading.Temperature + 61.0 + ((reading.Temperature - 68.0) * 1.2) + (reading.Humidity * 0.094))

		// Calculate air quality index
		airQualityIndex := 1
		switch {
		case reading.CO2 >= 5000:
			airQualityIndex = 4
		case reading.CO2 >= 2000:
			airQualityIndex = 3
		case reading.CO2 >= 1000:
			airQualityIndex = 2
		}

		ts := c.tsDeduper.Next(reading.DeviceID, reading.Timestamp)

		batch.Query(stmt,
			gocql.TimeUUID(),
			reading.DeviceID,
			reading.DeviceName,
			reading.DeviceType,
			reading.Location,
			reading.Floor,
			reading.Zone,
			reading.Temperature,
			reading.Humidity,
			reading.CO2,
			ts,
			time.Now(),
			heatIndex,
			airQualityIndex,
			`toTimestamp(now())`,
		)
	}

	if err := c.session.ExecuteBatch(batch); err != nil {
		return fmt.Errorf("failed to batch insert readings: %w", err)
	}

	return nil

	// stmt, _ := qb.Insert("sensor_readings").
	// 	Columns("id", "device_id", "device_name", "device_type", "location",
	// 		"floor_id", "zone_id", "temperature", "humidity", "co2",
	// 		"timestamp", "created_at", "heat_index", "air_quality_index").
	// 	ToCql()

	// // Limit in-flight writes
	// const maxInFlight = 128
	// sem := make(chan struct{}, maxInFlight)
	// errCh := make(chan error, len(readings))
	// var wg sync.WaitGroup

	// for _, r := range readings {
	// 	r := r
	// 	wg.Add(1)
	// 	sem <- struct{}{}
	// 	go func() {
	// 		defer wg.Done()
	// 		defer func() { <-sem }()
	// 		heatIndex := 0.5 * (r.Temperature + 61.0 + ((r.Temperature - 68.0) * 1.2) + (r.Humidity * 0.094))
	// 		aqi := 1
	// 		switch {
	// 		case r.CO2 >= 5000:
	// 			aqi = 4
	// 		case r.CO2 >= 2000:
	// 			aqi = 3
	// 		case r.CO2 >= 1000:
	// 			aqi = 2
	// 		}
	// 		err := c.session.Session.Query(stmt,
	// 			gocql.TimeUUID(), r.DeviceID, r.DeviceName, r.DeviceType, r.Location,
	// 			r.Floor, r.Zone, r.Temperature, r.Humidity, r.CO2,
	// 			r.Timestamp, time.Now(), heatIndex, aqi,
	// 		).Consistency(gocql.One).WithContext(ctx).Exec()
	// 		if err != nil {
	// 			errCh <- err
	// 		}
	// 	}()
	// }
	// wg.Wait()
	// close(errCh)
	// for err := range errCh {
	// 	if err != nil {
	// 		return fmt.Errorf("cass write error: %w", err)
	// 	}
	// }
}

func (c *cassandraImpl) GetReadings(ctx context.Context, input model.GetReadingsInput) ([]model.SensorReading, error) {
	builder := qb.Select("sensor_readings").
		Columns("id", "device_id", "device_name", "device_type", "location",
			"floor_id", "zone_id", "temperature", "humidity", "co2",
			"timestamp", "created_at")

	if input.DeviceID != "" {
		builder = builder.Where(qb.Eq("device_id"))
	}
	if !input.StartTime.IsZero() {
		builder = builder.Where(qb.GtOrEq("timestamp"))
	}
	if !input.EndTime.IsZero() {
		builder = builder.Where(qb.LtOrEq("timestamp"))
	}

	stmt, names := builder.ToCql()
	q := c.session.Query(stmt, names)
	defer q.Release()

	if input.DeviceID != "" {
		q.BindMap(qb.M{"device_id": input.DeviceID})
	}
	if !input.StartTime.IsZero() {
		q.BindMap(qb.M{"timestamp": input.StartTime})
	}
	if !input.EndTime.IsZero() {
		q.BindMap(qb.M{"timestamp": input.EndTime})
	}

	var results []map[string]interface{}
	if err := q.SelectRelease(&results); err != nil {
		return nil, fmt.Errorf("failed to query readings: %w", err)
	}

	var readings []model.SensorReading
	for _, row := range results {
		readings = append(readings, model.SensorReading{
			ID:          row["id"].(int64),
			DeviceID:    row["device_id"].(string),
			DeviceName:  row["device_name"].(string),
			DeviceType:  row["device_type"].(string),
			Location:    row["location"].(string),
			Floor:       row["floor_id"].(int),
			Zone:        row["zone_id"].(int),
			Temperature: row["temperature"].(float64),
			Humidity:    row["humidity"].(float64),
			CO2:         row["co2"].(float64),
			Timestamp:   row["timestamp"].(time.Time),
			CreatedAt:   row["created_at"].(time.Time),
		})
	}

	return readings, nil
}

func (c *cassandraImpl) GetBenchmarkMetrics(ctx context.Context, limit int) ([]model.BenchmarkMetrics, error) {
	stmt, names := qb.Select("benchmark_metrics").
		Columns("id", "total_records", "processed_records", "failed_records",
			"start_time", "end_time", "average_latency", "throughput",
			"batch_size", "database_type", "created_at").
		Limit(uint(limit)).
		OrderBy("created_at", qb.DESC).
		ToCql()

	q := c.session.Query(stmt, names)
	defer q.Release()

	var results []map[string]interface{}
	if err := q.SelectRelease(&results); err != nil {
		return nil, fmt.Errorf("failed to query benchmark metrics: %w", err)
	}

	var metrics []model.BenchmarkMetrics
	for _, row := range results {
		metrics = append(metrics, model.BenchmarkMetrics{
			TotalRecords:     row["total_records"].(int64),
			ProcessedRecords: row["processed_records"].(int64),
			FailedRecords:    row["failed_records"].(int64),
			StartTime:        row["start_time"].(time.Time),
			EndTime:          row["end_time"].(time.Time),
			AverageLatency:   row["average_latency"].(float64),
			Throughput:       row["throughput"].(float64),
			BatchSize:        row["batch_size"].(int),
			DatabaseType:     row["database_type"].(string),
		})
	}

	return metrics, nil
}

func (c *cassandraImpl) SaveBenchmarkMetrics(ctx context.Context, metrics model.BenchmarkMetrics) error {
	stmt, names := qb.Insert("benchmark_metrics").
		Columns("id", "total_records", "processed_records", "failed_records",
			"start_time", "end_time", "average_latency", "end_to_end_latency", "throughput",
			"batch_size", "database_type", "created_at").
		ToCql()

	q := c.session.Query(stmt, names).BindMap(qb.M{
		"id":                 gocql.TimeUUID(),
		"total_records":      metrics.TotalRecords,
		"processed_records":  metrics.ProcessedRecords,
		"failed_records":     metrics.FailedRecords,
		"start_time":         metrics.StartTime,
		"end_time":           metrics.EndTime,
		"average_latency":    metrics.AverageLatency,
		"end_to_end_latency": metrics.EndToEndLatency,
		"throughput":         metrics.Throughput,
		"batch_size":         metrics.BatchSize,
		"database_type":      metrics.DatabaseType,
		"created_at":         time.Now(),
	})

	if err := q.Exec(); err != nil {
		return fmt.Errorf("failed to save benchmark metrics: %w", err)
	}

	q.Release()

	return nil
}

func (c *cassandraImpl) GetLatestReadings(ctx context.Context) ([]model.SensorReading, error) {
	// Get all unique device IDs first
	deviceStmt, deviceNames := qb.Select("iot_devices").
		Columns("device_id").
		ToCql()

	deviceQuery := c.session.Query(deviceStmt, deviceNames)
	defer deviceQuery.Release()

	var deviceResults []map[string]interface{}
	if err := deviceQuery.SelectRelease(&deviceResults); err != nil {
		return nil, fmt.Errorf("failed to query devices: %w", err)
	}

	// For each device, get the latest reading
	var readings []model.SensorReading
	for _, device := range deviceResults {
		deviceID := device["device_id"].(string)

		stmt, names := qb.Select("sensor_readings").
			Columns("id", "device_id", "device_name", "device_type", "location",
				"floor_id", "zone_id", "temperature", "humidity", "co2",
				"timestamp", "created_at", "heat_index", "air_quality_index").
			Where(qb.Eq("device_id")).
			OrderBy("timestamp", qb.DESC).
			Limit(1).
			ToCql()

		q := c.session.Query(stmt, names).BindMap(qb.M{
			"device_id": deviceID,
		})
		defer q.Release()

		var results []map[string]interface{}
		if err := q.SelectRelease(&results); err != nil {
			return nil, fmt.Errorf("failed to query latest reading for device %s: %w", deviceID, err)
		}

		// Skip if no readings for this device
		if len(results) == 0 {
			continue
		}

		row := results[0] // Get the latest reading
		readings = append(readings, model.SensorReading{
			DeviceID:    row["device_id"].(string),
			DeviceName:  row["device_name"].(string),
			DeviceType:  row["device_type"].(string),
			Location:    row["location"].(string),
			Floor:       row["floor_id"].(int),
			Zone:        row["zone_id"].(int),
			Temperature: row["temperature"].(float64),
			Humidity:    row["humidity"].(float64),
			CO2:         row["co2"].(float64),
			Timestamp:   row["timestamp"].(time.Time),
			CreatedAt:   row["created_at"].(time.Time),
		})
	}

	return readings, nil
}
