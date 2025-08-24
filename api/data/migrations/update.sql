-- Enable TimescaleDB and hypertable on the time column
CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE sensor_readings (
    id            BIGSERIAL PRIMARY KEY,
    device_id     VARCHAR(255) NOT NULL,
    device_name   VARCHAR(255) NOT NULL,
    device_type   VARCHAR(50)  NOT NULL,
    location      VARCHAR(255) NOT NULL,
    floor         INTEGER      NOT NULL,
    zone          VARCHAR(100) NOT NULL,
    temperature   DOUBLE PRECISION,
    humidity      DOUBLE PRECISION,
    co2           DOUBLE PRECISION,
    "timestamp"   TIMESTAMP    NOT NULL,
    created_at    TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    -- Idempotency key for exactly-once-ish writes:
    CONSTRAINT uq_device_ts UNIQUE (device_id, "timestamp")
);

SELECT create_hypertable('sensor_readings', 'timestamp', if_not_exists => TRUE);

-- Helpful indexes for queries
CREATE INDEX IF NOT EXISTS idx_readings_device_ts_desc
  ON sensor_readings (device_id, "timestamp" DESC);
CREATE INDEX IF NOT EXISTS idx_readings_location ON sensor_readings (location);
CREATE INDEX IF NOT EXISTS idx_readings_zone     ON sensor_readings (zone);
CREATE INDEX IF NOT EXISTS idx_readings_floor    ON sensor_readings (floor);