# IoT System Observability

This document describes the comprehensive observability setup for the IoT system, including Prometheus exporters, Grafana dashboards, and CPU tracking in Docker.

## Overview

The observability stack provides:

- **Metrics Collection**: Prometheus with custom IoT metrics
- **Visualization**: Grafana dashboards
- **System Monitoring**: Node Exporter, cAdvisor for container metrics
- **Database Monitoring**: PostgreSQL Exporter
- **Message Queue Monitoring**: Kafka Exporter

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   IoT Devices   │───▶│   Kafka Queue   │───▶│  IoT Consumer   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Grafana       │◀───│   Prometheus    │◀───│  Metrics Server │
│   Dashboards    │    │   (Metrics DB)  │    │  (:9091)        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  Node Exporter  │    │    cAdvisor     │    │  PostgreSQL     │
│  (Host Metrics) │    │ (Container Metrics) │  Exporter       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Quick Start

### 1. Start the Full System with Monitoring

```bash
# Start everything including observability stack
make run-with-monitoring
```

### 2. Start Only Monitoring Services

```bash
# Start only the observability stack
make monitoring-only
```

### 3. View Monitoring URLs

```bash
make monitoring-urls
```

## Monitoring URLs

| Service                      | URL                           | Credentials |
| ---------------------------- | ----------------------------- | ----------- |
| Grafana Dashboard            | http://localhost:3000         | admin/admin |
| Prometheus                   | http://localhost:9090         | -           |
| cAdvisor (Container Metrics) | http://localhost:8080         | -           |
| Node Exporter (Host Metrics) | http://localhost:9100/metrics | -           |
| IoT API Metrics              | http://localhost:3001/metrics | -           |

## Metrics Collected

### System Metrics

- CPU usage percentage
- Memory allocation and heap usage
- Disk usage
- Number of goroutines
- Container resource usage (CPU, memory, network)

### IoT Application Metrics

- IoT messages received/processed/failed
- Batch processing time and size
- Sensor readings count
- Device count and online status
- Anomaly detection rate

### Kafka Metrics

- Messages consumed/produced
- Consumer lag
- Producer errors
- Topic metrics

### Database Metrics

- Active/idle connections
- Query duration
- Transaction count
- Error count

### HTTP Metrics

- Request rate by endpoint
- Response time
- Error rates (4xx, 5xx)
- Requests in flight

## Grafana Dashboards

### IoT System Overview Dashboard

- **System Overview**: CPU, memory, disk, goroutines
- **IoT Metrics**: Message rates, processing times, device counts
- **Kafka Metrics**: Consumer lag, message rates
- **Database Metrics**: Connections, query performance
- **HTTP Metrics**: Request rates, error rates

### Key Panels

1. **System Health**: Real-time system resource usage
2. **IoT Processing**: Message throughput and latency
3. **Kafka Performance**: Consumer lag and message rates
4. **Database Performance**: Connection pool and query metrics
5. **API Performance**: HTTP request metrics and error rates

## Monitoring and Visualization

The system provides comprehensive monitoring and visualization capabilities without automated alerting. Key features include:

### Real-time Monitoring

- Live system resource usage tracking
- IoT message processing metrics
- Infrastructure performance monitoring
- Business metrics visualization

### Manual Analysis

- Historical data analysis through Prometheus queries
- Custom dashboard creation in Grafana
- Performance trend analysis
- Capacity planning insights

### Health Checks

- Service availability monitoring
- Performance threshold tracking
- Resource utilization analysis
- Error rate monitoring

Note: While automated alerting is not implemented, the monitoring system provides comprehensive visibility for manual analysis and proactive monitoring.

## Configuration Files

### Prometheus Configuration

- **File**: `monitoring/prometheus/prometheus.yml`
- **Scrape Interval**: 15s (10s for IoT services)
- **Retention**: 200 hours
- **Targets**: All IoT services, system metrics, database, Kafka

### Grafana Configuration

- **Datasource**: Prometheus (auto-configured)
- **Dashboards**: Auto-provisioned from JSON files
- **Refresh**: 5s

## Custom Metrics

### Adding New Metrics

1. **Define in Prometheus Metrics Package**:

```go
// In api/internal/pkg/metrics/prometheus.go
newMetric := promauto.NewCounter(prometheus.CounterOpts{
    Namespace: namespace,
    Name:      "new_metric_total",
    Help:      "Description of the new metric",
})
```

2. **Record in Application Code**:

```go
// In your application code
if c.promMetrics != nil {
    c.promMetrics.RecordNewMetric()
}
```

3. **Add to Grafana Dashboard**:

```json
{
  "targets": [
    {
      "expr": "rate(new_metric_total[5m])",
      "legendFormat": "New Metric Rate"
    }
  ]
}
```

## Troubleshooting

### Common Issues

1. **Prometheus Can't Scrape Metrics**

   - Check if metrics server is running on :9091
   - Verify network connectivity between containers
   - Check Prometheus configuration

2. **Grafana Can't Connect to Prometheus**

   - Verify Prometheus is running on :9090
   - Check datasource configuration
   - Ensure containers are on the same network

3. **No Metrics Appearing**
   - Check if IoT application is generating metrics
   - Verify metrics server is started
   - Check Prometheus targets page

### Debug Commands

```bash
# Check container status
docker ps

# View logs
make observability-logs

# Check Prometheus targets
curl http://localhost:9090/api/v1/targets

# Check metrics endpoint
curl http://localhost:3001/metrics

# Test Prometheus rules
curl http://localhost:9090/api/v1/rules
```

## Performance Considerations

### Resource Usage

- **Prometheus**: ~100MB RAM, 1-2 CPU cores
- **Grafana**: ~50MB RAM, 0.5 CPU cores
- **Node Exporter**: ~10MB RAM, minimal CPU
- **cAdvisor**: ~50MB RAM, 0.5 CPU cores

### Scaling

- For high-volume IoT deployments, consider:
  - Prometheus federation for multiple instances
  - Grafana clustering
  - Dedicated monitoring infrastructure
  - Metrics aggregation and sampling

## Security

### Access Control

- Grafana: admin/admin (change in production)
- Prometheus: No authentication (add reverse proxy)

### Network Security

- All monitoring services run on internal Docker network
- External access only through port mappings
- Consider VPN or reverse proxy for production

## Production Deployment

### Recommendations

1. **Use dedicated monitoring infrastructure**
2. **Implement proper authentication**
3. **Set up backup and retention policies**
4. **Configure manual monitoring procedures**
5. **Monitor the monitoring system itself**
6. **Use persistent volumes for data storage**

### Environment Variables

```bash
# Grafana
GF_SECURITY_ADMIN_PASSWORD=secure_password
GF_USERS_ALLOW_SIGN_UP=false

# Prometheus
PROMETHEUS_RETENTION_TIME=30d
```

## Manual Monitoring Procedures

Since automated alerting is not implemented, consider these manual monitoring practices:

### Daily Checks

- Review Grafana dashboards for system health
- Check Prometheus targets for any down services
- Monitor resource usage trends
- Review error rates and performance metrics

### Weekly Analysis

- Analyze performance trends
- Review capacity utilization
- Check for any anomalies in metrics
- Plan capacity adjustments if needed

### Monthly Reviews

- Comprehensive system health review
- Performance optimization opportunities
- Capacity planning for growth
- Monitoring system improvements

## Support

For issues with the observability stack:

1. Check the troubleshooting section
2. Review container logs
3. Verify configuration files
4. Test individual components
5. Check network connectivity between services
