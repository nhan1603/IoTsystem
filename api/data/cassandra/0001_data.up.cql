CREATE KEYSPACE IF NOT EXISTS iotsystem
    WITH replication = {
        'class': 'SimpleStrategy',
        'replication_factor': 1
    }
    AND durable_writes = true;

USE iotsystem;

CREATE TABLE IF NOT EXISTS iot_devices (
    id uuid,
    device_id text,
    name text,
    type text,
    location text,
    floor_id int,
    zone_id int,
    is_active boolean,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY (device_id)
);

CREATE TABLE IF NOT EXISTS sensor_readings (
    id timeuuid,
    device_id text,
    device_name text,
    device_type text,
    location text,
    floor_id int,
    zone_id int,
    temperature double,
    humidity double,
    co2 double,
    timestamp timestamp,
    created_at timestamp,
    // Add computed columns
    heat_index double,      // Storing pre-computed heat index
    air_quality_index int,  // Storing pre-computed air quality index
    durable_write_ts timestamp,
    PRIMARY KEY ((device_id), timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC);

CREATE TABLE IF NOT EXISTS benchmark_metrics (
    id timeuuid,
    total_records bigint,
    processed_records bigint,
    failed_records bigint,
    start_time timestamp,
    end_time timestamp,
    average_latency double,
    end_to_end_latency double,
    throughput double,
    batch_size int,
    database_type text,
    created_at timestamp,
    PRIMARY KEY (id)
);