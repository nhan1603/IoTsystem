-- Insert sample IoT devices
INSERT INTO iot_devices
    (device_id, name, type, location, floor, zone)
VALUES
    ('TEMP_001', 'Temperature Sensor 1', 'temperature', 'Main Building', 1, 'Zone A'),
    ('TEMP_002', 'Temperature Sensor 2', 'temperature', 'Main Building', 2, 'Zone B'),
    ('HUM_001', 'Humidity Sensor 1', 'humidity', 'Main Building', 1, 'Zone A'),
    ('HUM_002', 'Humidity Sensor 2', 'humidity', 'Main Building', 2, 'Zone B'),
    ('CO2_001', 'CO2 Sensor 1', 'co2', 'Main Building', 1, 'Zone A'),
    ('CO2_002', 'CO2 Sensor 2', 'co2', 'Main Building', 2, 'Zone B'),
    ('MULTI_001', 'Multi Sensor 1', 'multi', 'Main Building', 1, 'Zone A'),
    ('MULTI_002', 'Multi Sensor 2', 'multi', 'Main Building', 2, 'Zone B'); 