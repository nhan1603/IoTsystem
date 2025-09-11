package model

import (
	"time"
)

// IoTDevice represents an IoT sensor device
type IoTDevice struct {
	ID        int64
	DeviceID  string
	Name      string
	Type      string
	Location  string
	Floor     int
	Zone      int
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SensorReading represents a single sensor data reading
type SensorReading struct {
	ID          int64
	DeviceID    string
	DeviceName  string
	DeviceType  string
	Location    string
	Floor       int
	Zone        int
	Temperature float64
	Humidity    float64
	CO2         float64
	Timestamp   time.Time
	CreatedAt   time.Time
}

// IoTDataMessage represents the Kafka message format for IoT data
type IoTDataMessage struct {
	DeviceID    string    `json:"device_id"`
	DeviceName  string    `json:"device_name"`
	DeviceType  string    `json:"device_type"`
	Location    string    `json:"location"`
	Floor       int       `json:"floor"`
	Zone        int       `json:"zone"`
	Temperature float64   `json:"temperature"`
	Humidity    float64   `json:"humidity"`
	CO2         float64   `json:"co2"`
	Timestamp   time.Time `json:"timestamp"`
}

// GetReadingsInput represents input for querying sensor readings
type GetReadingsInput struct {
	DeviceID   string
	DeviceType string
	Location   string
	Floor      int
	Zone       int
	StartTime  time.Time
	EndTime    time.Time
	Limit      int
	Page       int
}

// DeviceType is enum for device types
type DeviceType string

const (
	// DeviceTypeTemperature is a temperature sensor
	DeviceTypeTemperature DeviceType = "temperature"
	// DeviceTypeHumidity is a humidity sensor
	DeviceTypeHumidity DeviceType = "humidity"
	// DeviceTypeCO2 is a CO2 sensor
	DeviceTypeCO2 DeviceType = "co2"
	// DeviceTypeMulti is a multi-sensor device
	DeviceTypeMulti DeviceType = "multi"
)

// ToString convert enum to string
func (dt DeviceType) ToString() string {
	return string(dt)
}

// BenchmarkMetrics represents performance metrics for benchmarking
type BenchmarkMetrics struct {
	TotalRecords     int64
	ProcessedRecords int64
	FailedRecords    int64
	StartTime        time.Time
	EndTime          time.Time
	AverageLatency   float64
	EndToEndLatency  float64
	Throughput       float64 // records per second
	BatchSize        int
	DatabaseType     string
}
