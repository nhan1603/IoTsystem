-- Insert sample floors
INSERT INTO floors 
    (floor_number, description, total_area)
VALUES
    (1, 'Ground Floor', 1000.0),
    (2, 'First Floor', 1000.0);

-- Insert zones for each floor
INSERT INTO zones
    (floor_id, zone_name, zone_type, description, area)
VALUES
    -- Ground Floor (floor_id = 1) zones
    (1, 'Zone A', 'office', 'Main Office Space', 400.0),
    (1, 'Zone B', 'meeting', 'Meeting Rooms', 300.0),
    -- First Floor (floor_id = 2) zones
    (2, 'Zone A', 'laboratory', 'Research Lab', 500.0),
    (2, 'Zone B', 'storage', 'Storage Area', 200.0);

-- Insert IoT devices with proper foreign keys
INSERT INTO iot_devices
    (device_id, name, type, location, floor_id, zone_id)
VALUES
    -- Ground Floor, Zone A devices
    ('TEMP_001', 'Temperature Sensor 1', 'temperature', 'Main Building', 1, 1),
    ('HUM_001', 'Humidity Sensor 1', 'humidity', 'Main Building', 1, 1),
    ('CO2_001', 'CO2 Sensor 1', 'co2', 'Main Building', 1, 1),
    ('MULTI_001', 'Multi Sensor 1', 'multi', 'Main Building', 1, 1),
    
    -- Ground Floor, Zone B devices
    ('TEMP_002', 'Temperature Sensor 2', 'temperature', 'Main Building', 1, 2),
    ('HUM_101', 'Humidity Sensor 1.5', 'humidity', 'Main Building', 1, 2),
    
    -- First Floor, Zone A devices
    ('HUM_002', 'Humidity Sensor 2', 'humidity', 'Main Building', 2, 3),
    ('CO2_002', 'CO2 Sensor 2', 'co2', 'Main Building', 2, 3),
    
    -- First Floor, Zone B devices
    ('MULTI_002', 'Multi Sensor 2', 'multi', 'Main Building', 2, 4),
    ('CO2_003', 'CO2 Sensor 3', 'co2', 'Main Building', 2, 4);