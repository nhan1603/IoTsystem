CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Create users table
CREATE TABLE users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE floors (
    id SERIAL PRIMARY KEY,
    floor_number INTEGER NOT NULL UNIQUE,
    description VARCHAR(255),
    total_area DOUBLE PRECISION,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Add building zones table
CREATE TABLE zones (
    id SERIAL PRIMARY KEY,
    floor_id INTEGER REFERENCES floors(id),
    zone_name VARCHAR(100) NOT NULL,
    zone_type VARCHAR(50) NOT NULL,
    description TEXT,
    area DOUBLE PRECISION,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (floor_id, zone_name)
);

-- Create IoT devices table
CREATE TABLE iot_devices
(
    id BIGSERIAL PRIMARY KEY,
    device_id VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    location VARCHAR(255) NOT NULL,
    floor_id INTEGER NOT NULL REFERENCES floors(id),
    zone_id INTEGER NOT NULL REFERENCES zones(id),
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
    floor_id INTEGER NOT NULL REFERENCES floors(id),
    zone_id INTEGER NOT NULL REFERENCES zones(id),
    temperature DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    humidity DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    co2 DOUBLE PRECISION NOT NULL DEFAULT 0.0,
    timestamp TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id, timestamp),
    -- Idempotency key for exactly-once-ish writes:
    CONSTRAINT uq_device_ts UNIQUE (device_id, "timestamp")
);

SELECT create_hypertable('sensor_readings', 'timestamp', if_not_exists => TRUE);

-- Create index on device_id and timestamp for efficient queries
CREATE INDEX idx_sensor_readings_device_timestamp ON sensor_readings (device_id, timestamp DESC);
CREATE INDEX idx_sensor_readings_location ON sensor_readings (location);
CREATE INDEX idx_sensor_readings_floor_id ON sensor_readings(floor_id);
CREATE INDEX idx_sensor_readings_zone_id ON sensor_readings(zone_id);
CREATE INDEX idx_zones_floor_id ON zones(floor_id);

ALTER TABLE sensor_readings
    ADD COLUMN heat_index DOUBLE PRECISION GENERATED ALWAYS AS (
        0.5 * (temperature + 61.0 + ((temperature-68.0)*1.2) + (humidity*0.094))
    ) STORED,
    ADD COLUMN air_quality_index INTEGER GENERATED ALWAYS AS (
        CASE 
            WHEN co2 < 1000 THEN 1
            WHEN co2 < 2000 THEN 2
            WHEN co2 < 5000 THEN 3
            ELSE 4
        END
    ) STORED;

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