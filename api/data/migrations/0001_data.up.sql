-- Create users table
CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create IoT devices table
CREATE TABLE iot_devices
(
    id BIGSERIAL PRIMARY KEY,
    device_id VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    location VARCHAR(255) NOT NULL,
    floor INTEGER NOT NULL,
    zone VARCHAR(100) NOT NULL,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create sensor readings table with TimescaleDB hypertable
CREATE TABLE sensor_readings
(
    id BIGSERIAL,
    device_id VARCHAR(255) NOT NULL,
    device_name VARCHAR(255) NOT NULL,
    device_type VARCHAR(50) NOT NULL,
    location VARCHAR(255) NOT NULL,
    floor INTEGER NOT NULL,
    zone VARCHAR(100) NOT NULL,
    temperature DOUBLE PRECISION,
    humidity DOUBLE PRECISION,
    co2 DOUBLE PRECISION,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, timestamp)
);

-- Create index on device_id and timestamp for efficient queries
CREATE INDEX idx_sensor_readings_device_timestamp ON sensor_readings (device_id, timestamp DESC);
CREATE INDEX idx_sensor_readings_location ON sensor_readings (location);
CREATE INDEX idx_sensor_readings_zone ON sensor_readings (zone);
CREATE INDEX idx_sensor_readings_floor ON sensor_readings (floor);

-- Create benchmark metrics table
CREATE TABLE benchmark_metrics
(
    id BIGSERIAL PRIMARY KEY,
    total_records BIGINT NOT NULL,
    processed_records BIGINT NOT NULL,
    failed_records BIGINT NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    average_latency DOUBLE PRECISION NOT NULL,
    throughput DOUBLE PRECISION NOT NULL,
    batch_size INTEGER NOT NULL,
    database_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);