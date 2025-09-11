SELECT * FROM public.benchmark_metrics
ORDER BY id DESC limit 10;

SELECT * from public.sensor_readings limit 1;

SELECT count(*) from public.sensor_readings;

-- Basic time window and throughput calculation
WITH metrics AS (
    SELECT 
        count(*) as total_records,
        min(timestamp) as window_start,
        max(timestamp) as window_end,
        extract(epoch from (max(timestamp) - min(timestamp))) as window_seconds
    FROM sensor_readings
)
SELECT 
    total_records,
    window_start,
    window_end,
    window_seconds,
    CASE 
        WHEN window_seconds > 0 
        THEN round((total_records::numeric / window_seconds)::numeric, 2)
        ELSE 0 
    END as events_per_second
FROM metrics;