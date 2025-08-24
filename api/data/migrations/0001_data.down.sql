DROP INDEX IF EXISTS idx_sensor_readings_floor;
DROP INDEX IF EXISTS idx_sensor_readings_zone;
-- Drop tables in reverse order
DROP TABLE IF EXISTS benchmark_metrics;
DROP TABLE IF EXISTS sensor_readings;
DROP TABLE IF EXISTS iot_devices; 
DROP TABLE IF EXISTS users;