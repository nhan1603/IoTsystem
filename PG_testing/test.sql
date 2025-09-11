-- Overall system throughput
SELECT database_type,
    AVG(throughput) as avg_throughput,
    MAX(throughput) as max_throughput,
    MIN(throughput) as min_throughput,
    AVG(average_latency) as avg_latency
FROM benchmark_metrics
WHERE created_at >= NOW() - INTERVAL '24 hours'
GROUP BY database_type;

-- Batch size impact on performance
SELECT batch_size, AVG(throughput) as avg_throughput,
    AVG(average_latency) as avg_latency, COUNT(*) as num_batches
FROM benchmark_metrics
WHERE created_at >= NOW() - INTERVAL '24 hours'
GROUP BY batch_size
ORDER BY batch_size;

-- Success rate calculation
SELECT database_type,
    SUM(processed_records) as total_processed,
    SUM(failed_records) as total_failed,
    ROUND((SUM(processed_records)::float / NULLIF(SUM(total_records), 0) * 100), 2) as success_rate
FROM benchmark_metrics
WHERE created_at >= NOW() - INTERVAL '24 hours'
GROUP BY database_type;


-- Count messages per device with timing info for last hour
SELECT device_id, COUNT(*) as message_count,
    MIN(timestamp) as first_message, MAX(timestamp) as last_message,
    (MAX(timestamp) - MIN(timestamp)) as time_span
FROM sensor_readings
WHERE timestamp >= NOW() - INTERVAL '1 hour'
GROUP BY device_id;