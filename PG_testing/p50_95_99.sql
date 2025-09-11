WITH win AS (
  SELECT EXTRACT(EPOCH FROM (durable_write_ts - "timestamp")) * 1000.0 AS lat_ms
  FROM sensor_readings
),
per_min AS (
  SELECT date_trunc('minute', durable_write_ts) AS m, COUNT(*) AS rows
  FROM sensor_readings
  GROUP BY 1
)
SELECT
  -- latency percentiles over the window
  percentile_disc(0.50) WITHIN GROUP (ORDER BY lat_ms) AS p50_ms,
  percentile_disc(0.95) WITHIN GROUP (ORDER BY lat_ms) AS p95_ms,
  percentile_disc(0.99) WITHIN GROUP (ORDER BY lat_ms) AS p99_ms,
  AVG(lat_ms)                                          AS mean_ms,
  (SELECT percentile_disc(0.50) WITHIN GROUP (ORDER BY rows)
     FROM per_min) / 60.0                               AS sustained_ev_s_median
FROM win;
