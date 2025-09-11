USE iotsystem;

-- Insert sample devices
INSERT INTO iot_devices (id, device_id, name, type, location, floor_id, zone_id, is_active, created_at, updated_at)
VALUES (now(), 'TEMP_001', 'Temperature Sensor 1', 'temperature', 'Main Building', 1, 1, true, toTimestamp(now()), toTimestamp(now()));

INSERT INTO iot_devices (id, device_id, name, type, location, floor_id, zone_id, is_active, created_at, updated_at)
VALUES (now(), 'HUM_001', 'Humidity Sensor 1', 'humidity', 'Main Building', 1, 1, true, toTimestamp(now()), toTimestamp(now()));

INSERT INTO iot_devices (id, device_id, name, type, location, floor_id, zone_id, is_active, created_at, updated_at)
VALUES (now(), 'CO2_001', 'CO2 Sensor 1', 'co2', 'Main Building', 1, 1, true, toTimestamp(now()), toTimestamp(now()));

INSERT INTO iot_devices (id, device_id, name, type, location, floor_id, zone_id, is_active, created_at, updated_at)
VALUES (now(), 'MULTI_001', 'Multi Sensor 1', 'multi', 'Main Building', 1, 1, true, toTimestamp(now()), toTimestamp(now()));

-- Ground Floor, Zone B devices
INSERT INTO iot_devices (id, device_id, name, type, location, floor_id, zone_id, is_active, created_at, updated_at)
VALUES (now(), 'TEMP_002', 'Temperature Sensor 2', 'temperature', 'Main Building', 1, 2, true, toTimestamp(now()), toTimestamp(now()));

INSERT INTO iot_devices (id, device_id, name, type, location, floor_id, zone_id, is_active, created_at, updated_at)
VALUES (now(), 'HUM_101', 'Humidity Sensor 1.5', 'humidity', 'Main Building', 1, 2, true, toTimestamp(now()), toTimestamp(now()));


-- First Floor, Zone A devices
INSERT INTO iot_devices (id, device_id, name, type, location, floor_id, zone_id, is_active, created_at, updated_at)
VALUES (now(), 'HUM_002', 'Humidity Sensor 2', 'humidity', 'Main Building', 2, 3, true, toTimestamp(now()), toTimestamp(now()));

INSERT INTO iot_devices (id, device_id, name, type, location, floor_id, zone_id, is_active, created_at, updated_at)
VALUES (now(), 'CO2_002', 'CO2 Sensor 2', 'co2', 'Main Building', 2, 3, true, toTimestamp(now()), toTimestamp(now()));

INSERT INTO iot_devices (id, device_id, name, type, location, floor_id, zone_id, is_active, created_at, updated_at)
VALUES (now(), 'MULTI_002', 'Multi Sensor 2', 'multi', 'Main Building', 2, 4, true, toTimestamp(now()), toTimestamp(now()));

INSERT INTO iot_devices (id, device_id, name, type, location, floor_id, zone_id, is_active, created_at, updated_at)
VALUES (now(), 'CO2_003', 'CO2 Sensor 3', 'co2', 'Main Building', 2, 4, true, toTimestamp(now()), toTimestamp(now()));