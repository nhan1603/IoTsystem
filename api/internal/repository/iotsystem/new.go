package iotsystem

import (
	"context"
	"fmt"
	"strings"

	"github.com/aarondl/null/v8"
	"github.com/aarondl/sqlboiler/v4/boil"
	"github.com/aarondl/sqlboiler/v4/queries/qm"
	"github.com/nhan1603/IoTsystem/api/internal/model"
	"github.com/nhan1603/IoTsystem/api/internal/repository/dbmodel"
)

// Repository defines IoT data repository interface
type Repository interface {
	GetDevices(ctx context.Context) ([]model.IoTDevice, error)
	GetReadings(ctx context.Context, input model.GetReadingsInput) ([]model.SensorReading, error)
	GetLatestReadings(ctx context.Context) ([]model.SensorReading, error)
	BatchInsertReadings(ctx context.Context, readings []model.SensorReading) error
	GetBenchmarkMetrics(ctx context.Context, limit int) ([]model.BenchmarkMetrics, error)
	SaveBenchmarkMetrics(ctx context.Context, metrics model.BenchmarkMetrics) error
}

// New returns an implementation instance satisfying Repository
func New(dbConn boil.ContextExecutor) Repository {
	return impl{
		dbConn: dbConn,
	}

}

type impl struct {
	dbConn boil.ContextExecutor
}

// GetDevices retrieves all IoT devices
func (r impl) GetDevices(ctx context.Context) ([]model.IoTDevice, error) {
	dbDevice, err := dbmodel.IotDevices(
		dbmodel.IotDeviceWhere.IsActive.EQ(null.BoolFrom(true)),
		qm.OrderBy(dbmodel.IotDeviceColumns.DeviceID)).All(ctx, r.dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to query devices: %w", err)
	}
	var devices []model.IoTDevice
	for _, row := range dbDevice {
		devices = append(devices, model.IoTDevice{
			ID:        row.ID,
			DeviceID:  row.DeviceID,
			Name:      row.Name,
			Type:      row.Type,
			Location:  row.Location,
			Floor:     row.FloorID,
			Zone:      row.ZoneID,
			IsActive:  row.IsActive.Bool,
			CreatedAt: row.CreatedAt.Time,
			UpdatedAt: row.UpdatedAt.Time,
		})
	}

	return devices, nil
}

// GetReadings retrieves sensor readings with filters
func (r impl) GetReadings(ctx context.Context, input model.GetReadingsInput) ([]model.SensorReading, error) {
	mods := []qm.QueryMod{}

	if input.DeviceID != "" {
		mods = append(mods, dbmodel.SensorReadingWhere.DeviceID.EQ(input.DeviceID))
	}
	if input.DeviceType != "" {
		mods = append(mods, dbmodel.SensorReadingWhere.DeviceType.EQ(input.DeviceType))
	}
	if input.Location != "" {
		mods = append(mods, dbmodel.SensorReadingWhere.Location.EQ(input.Location))
	}
	if input.Floor > 0 {
		mods = append(mods, dbmodel.SensorReadingWhere.FloorID.EQ(input.Floor))
	}
	if input.Zone > 0 {
		mods = append(mods, dbmodel.SensorReadingWhere.ZoneID.EQ(input.Zone))
	}
	if !input.StartTime.IsZero() {
		mods = append(mods, dbmodel.SensorReadingWhere.Timestamp.GTE(input.StartTime))
	}
	if !input.EndTime.IsZero() {
		mods = append(mods, dbmodel.SensorReadingWhere.Timestamp.LTE(input.EndTime))
	}
	mods = append(mods, qm.OrderBy("timestamp DESC"))
	if input.Limit > 0 {
		mods = append(mods, qm.Limit(input.Limit))
	}

	dbReadings, err := dbmodel.SensorReadings(mods...).All(ctx, r.dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to query readings: %w", err)
	}

	var readings []model.SensorReading
	for _, row := range dbReadings {
		readings = append(readings, model.SensorReading{
			ID:          row.ID,
			DeviceID:    row.DeviceID,
			DeviceName:  row.DeviceName,
			DeviceType:  row.DeviceType,
			Location:    row.Location,
			Floor:       row.FloorID,
			Zone:        row.ZoneID,
			Temperature: row.Temperature,
			Humidity:    row.Humidity,
			CO2:         row.Co2,
			Timestamp:   row.Timestamp,
			CreatedAt:   row.CreatedAt.Time,
		})
	}
	return readings, nil
}

// GetLatestReadings retrieves the latest reading for each device
func (r impl) GetLatestReadings(ctx context.Context) ([]model.SensorReading, error) {
	query := `
		SELECT DISTINCT ON (device_id) 
		       id, device_id, device_name, device_type, location, floor, zone,
		       temperature, humidity, co2, timestamp, created_at
		FROM sensor_readings
		ORDER BY device_id, timestamp DESC
	`

	rows, err := r.dbConn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest readings: %w", err)
	}
	defer rows.Close()

	var readings []model.SensorReading
	for rows.Next() {
		var reading model.SensorReading
		err := rows.Scan(
			&reading.ID,
			&reading.DeviceID,
			&reading.DeviceName,
			&reading.DeviceType,
			&reading.Location,
			&reading.Floor,
			&reading.Zone,
			&reading.Temperature,
			&reading.Humidity,
			&reading.CO2,
			&reading.Timestamp,
			&reading.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan reading: %w", err)
		}
		readings = append(readings, reading)
	}

	return readings, nil
}

// Do this inTX
// BatchInsertReadings inserts multiple sensor readings in a batch
func (r impl) BatchInsertReadings(ctx context.Context, readings []model.SensorReading) error {
	if len(readings) == 0 {
		return nil
	}

	// Prepare batch insert statement
	query := `
		INSERT INTO sensor_readings 
		(device_id, device_name, device_type, location, floor_id, zone_id, 
		 temperature, humidity, co2, timestamp, created_at)
		VALUES 
	`

	var valueStrings []string
	var valueArgs []interface{}
	argCount := 0

	for _, reading := range readings {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			argCount+1, argCount+2, argCount+3, argCount+4, argCount+5, argCount+6,
			argCount+7, argCount+8, argCount+9, argCount+10, argCount+11))

		valueArgs = append(valueArgs,
			reading.DeviceID,
			reading.DeviceName,
			reading.DeviceType,
			reading.Location,
			reading.Floor,
			reading.Zone,
			reading.Temperature,
			reading.Humidity,
			reading.CO2,
			reading.Timestamp,
			reading.CreatedAt,
		)
		argCount += 11
	}

	query += strings.Join(valueStrings, ",")
	query += `
        ON CONFLICT ON CONSTRAINT uq_device_ts 
        DO UPDATE SET 
            temperature = sensor_readings.temperature,
            humidity = sensor_readings.humidity,
            co2 = sensor_readings.co2
    `

	// Execute batch insert
	_, err := r.dbConn.ExecContext(ctx, query, valueArgs...)
	if err != nil {
		return fmt.Errorf("failed to batch insert readings: %w", err)
	}

	return nil
}

// GetBenchmarkMetrics retrieves benchmark performance metrics using SQLBoiler query mods
func (r impl) GetBenchmarkMetrics(ctx context.Context, limit int) ([]model.BenchmarkMetrics, error) {
	mods := []qm.QueryMod{
		qm.OrderBy("created_at DESC"),
	}
	if limit > 0 {
		mods = append(mods, qm.Limit(limit))
	}

	dbMetrics, err := dbmodel.BenchmarkMetrics(mods...).All(ctx, r.dbConn)
	if err != nil {
		return nil, fmt.Errorf("failed to query benchmark metrics: %w", err)
	}

	var metrics []model.BenchmarkMetrics
	for _, row := range dbMetrics {
		metrics = append(metrics, model.BenchmarkMetrics{
			TotalRecords:     row.TotalRecords,
			ProcessedRecords: row.ProcessedRecords,
			FailedRecords:    row.FailedRecords,
			StartTime:        row.StartTime,
			EndTime:          row.EndTime,
			AverageLatency:   row.AverageLatency,
			Throughput:       row.Throughput,
			BatchSize:        row.BatchSize,
			DatabaseType:     row.DatabaseType,
		})
	}

	return metrics, nil
}

// SaveBenchmarkMetrics saves benchmark performance metrics
func (r impl) SaveBenchmarkMetrics(ctx context.Context, metrics model.BenchmarkMetrics) error {

	// Prepare insert query
	metricsDb := dbmodel.BenchmarkMetric{
		TotalRecords:     metrics.TotalRecords,
		ProcessedRecords: metrics.ProcessedRecords,
		FailedRecords:    metrics.FailedRecords,
		StartTime:        metrics.StartTime,
		EndTime:          metrics.EndTime,
		AverageLatency:   metrics.AverageLatency,
		Throughput:       metrics.Throughput,
		BatchSize:        metrics.BatchSize,
		DatabaseType:     metrics.DatabaseType,
	}

	err := metricsDb.Insert(ctx, r.dbConn, boil.Infer())
	if err != nil {
		return fmt.Errorf("failed to save benchmark metrics: %w", err)
	}

	return nil
}
