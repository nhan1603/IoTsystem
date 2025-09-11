Write-Host "Starting performance analysis..." -ForegroundColor Green

# Export data from Cassandra
Write-Host "Exporting data from Cassandra..."
try {
    docker exec cass1 cqlsh -e "COPY iotsystem.sensor_readings (device_id, timestamp, created_at) TO 'readings.csv' WITH HEADER = true;"
    docker cp cass1:readings.csv .
} catch {
    Write-Error "Failed to export data: $_"
    exit 1
}

# Process the data
try {
    # Read and parse CSV
    $readings = Import-Csv readings.csv
    
    # Calculate latencies in milliseconds
    $latencies = $readings | ForEach-Object {
        $created = [datetime]$_.created_at
        $timestamp = [datetime]$_.timestamp
        ($created - $timestamp).TotalMilliseconds
    } | Sort-Object

    # Group by minute for throughput calculation
    $byMinute = $readings | Group-Object { 
        [datetime]$_.created_at | Get-Date -Format "yyyy-MM-dd HH:mm"
    } | Select-Object @{
        Name='Count'; 
        Expression={$_.Count}
    } | Sort-Object Count

    # Calculate metrics
    $count = $latencies.Count
    $p50 = $latencies[[math]::Floor($count * 0.50)]
    $p95 = $latencies[[math]::Floor($count * 0.95)]
    $p99 = $latencies[[math]::Floor($count * 0.99)]
    $mean = ($latencies | Measure-Object -Average).Average

    # Calculate sustained events/second (median throughput per minute)
    $medianPerMinute = $byMinute[[math]::Floor($byMinute.Count * 0.50)].Count
    $sustainedEventsSec = $medianPerMinute / 60.0

    # Calculate overall throughput
    $firstReading = ($readings | Select-Object @{Name='ts';Expression={[datetime]$_.timestamp}} | Sort-Object ts | Select-Object -First 1).ts
    $lastReading = ($readings | Select-Object @{Name='ts';Expression={[datetime]$_.timestamp}} | Sort-Object ts -Descending | Select-Object -First 1).ts
    $totalSeconds = ($lastReading - $firstReading).TotalSeconds
    $overallThroughput = $count / $totalSeconds

    # Output results
    Write-Host "`nPerformance Metrics:" -ForegroundColor Cyan
    Write-Host "  Overall Throughput: $($overallThroughput.ToString('F2')) events/second"
    Write-Host "Latency Percentiles:"
    Write-Host "  P50: $($p50.ToString('F2')) ms"
    Write-Host "  P95: $($p95.ToString('F2')) ms"
    Write-Host "  P99: $($p99.ToString('F2')) ms"
    Write-Host "Mean Latency: $($mean.ToString('F2')) ms"
    Write-Host "Sustained Events/Second (median): $($sustainedEventsSec.ToString('F2'))"

    # Cleanup
    Remove-Item readings.csv -ErrorAction SilentlyContinue
} catch {
    Write-Error "Failed to process data: $_"
    exit 1
}