DELETE FROM benchmark_metrics;
DELETE FROM sensor_readings;
DELETE FROM iot_devices;
DELETE FROM zones;
DELETE FROM floors;
-- Note: Users table is not dropped to preserve user data