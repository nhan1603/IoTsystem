# IoTSystem

Simulation of IoT data ingestion system with comprehensive observability and monitoring capabilities.

## Features

- **IoT Data Ingestion**: Kafka-based message processing with PostgreSQL/TimescaleDB
- **Real-time Monitoring**: Prometheus metrics collection and Grafana dashboards
- **System Observability**: CPU, memory, disk, and container metrics tracking
- **Performance Monitoring**: HTTP, database, and Kafka metrics
- **Business Intelligence**: IoT device tracking and anomaly detection

## Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   IoT Simulator │──▶│   Kafka Queue   │───▶│  IoT Consumer   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                                       │
                                                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Grafana       │◀───│   Prometheus    │◀──│  Metrics Server │
│   Dashboards    │    │   (Metrics DB)  │    │  (:9091)        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                                 │                       │
                                 ▼                       ▼
                       ┌─────────────────┐    ┌─────────────────┐
                       │  Node Exporter  │    │    cAdvisor     │
                       │  (Host Metrics) │    │ (Container Metrics)
                       └─────────────────┘    └─────────────────┘
```

## Quick Start

### Backend Setup

- Install go 1.23
- Install `make` command
- cd `/api`
- `go mod vendor`
- Setup Docker and Docker Compose, and make sure it is running
- `cd ..`
- Setup project by running `make setup`
- Seed the necessary data into the database: `make pg-drop` then `make pg-migrate`

### Start the System

#### Option 1: Full System with Basic Monitoring (Recommended)

```bash
# Start everything including basic observability stack
make run-with-monitoring
```

#### Option 2: Full System with All Monitoring (Database + Kafka)

```bash
# Start everything including all monitoring exporters
make run-with-monitoring-all
```

#### Option 3: Basic System Only

```bash
# Start only the core system
make run
```

#### Option 4: Monitoring Only

```bash
# Start only the observability stack
make monitoring-only
```

### Frontend Setup

- Install nodejs
- cd `/web`
- `npm install`
- `npm start`
- The Web URL: `http://localhost:3000`

## Monitoring URLs

| Service                      | URL                           | Credentials | Status                                   |
| ---------------------------- | ----------------------------- | ----------- | ---------------------------------------- |
| Grafana Dashboard            | http://localhost:3000         | admin/admin | Always Available                         |
| Prometheus                   | http://localhost:9090         | -           | Always Available                         |
| cAdvisor (Container Metrics) | http://localhost:8080         | -           | Always Available                         |
| Node Exporter (Host Metrics) | http://localhost:9100/metrics | -           | Always Available                         |
| IoT API Metrics              | http://localhost:3001/metrics | -           | When API is running                      |
| PostgreSQL Exporter          | http://localhost:9187/metrics | -           | Optional (requires `--profile database`) |
| Kafka Exporter               | http://localhost:9308/metrics | -           | Optional (requires `--profile kafka`)    |

## Available Make Commands

### Core System

- `make setup` - Setup the project
- `make run` - Start the API server
- `make down` - Stop all containers
- `make test` - Run tests

### Observability

- `make run-with-monitoring` - Start full system with basic monitoring
- `make run-with-monitoring-all` - Start full system with all monitoring exporters
- `make monitoring-only` - Start only monitoring services
- `make monitoring-stop` - Stop monitoring services
- `make monitoring-urls` - Display monitoring URLs
- `make observability-logs` - View monitoring logs

#### Optional Monitoring Services

- `make observability-with-db` - Start monitoring with database exporter
- `make observability-with-kafka` - Start monitoring with Kafka exporter
- `make observability-with-all` - Start monitoring with all exporters

### Infrastructure

- `make zookeeper` - Start Zookeeper
- `make kafka` - Start Kafka
- `make kafka-topic` - Create Kafka topics
- `make api-simulator-run` - Start IoT simulator

## Metrics Collected

### System Metrics

- CPU usage percentage
- Memory allocation and heap usage
- Disk usage
- Number of goroutines
- Container resource usage

### IoT Application Metrics

- IoT messages received/processed/failed
- Batch processing time and size
- Sensor readings count
- Device count and online status
- Anomaly detection rate

### Infrastructure Metrics

- Kafka consumer lag and message rates
- Database connections and query performance
- HTTP request rates and error rates

## Monitoring and Visualization

The system includes comprehensive monitoring and visualization capabilities:

- **Real-time Dashboards**: Grafana dashboards for system and application metrics
- **Metrics Collection**: Prometheus for time-series data storage
- **Performance Monitoring**: HTTP, database, and Kafka metrics
- **System Health**: CPU, memory, disk, and container resource tracking

Note: Alerting functionality is not implemented in this setup. Metrics can be viewed through Grafana dashboards and Prometheus queries for manual monitoring and analysis.

## Documentation

- [Observability Guide](OBSERVABILITY.md) - Comprehensive monitoring documentation
- [Dataflow](Dataflow.txt) - System data flow diagram
- [Architecture](High_level_achitecture.drawio.xml) - High-level system architecture

## API Endpoints

- Health Check: `http://localhost:3001/api/public/ping`
- Metrics: `http://localhost:3001/metrics`
- IoT Devices: `http://localhost:3001/api/authenticated/v1/devices`
- Sensor Readings: `http://localhost:3001/api/authenticated/v1/readings`

## Development

### Adding New Metrics

1. Define metrics in `api/internal/pkg/obsmetrics/prometheus.go`
2. Record metrics in your application code
3. Add visualization panels to Grafana dashboards

## Troubleshooting

### Common Issues

1. **Prometheus can't scrape metrics**: Check if metrics server is running on :9091
2. **Grafana can't connect**: Verify Prometheus is running and accessible
3. **No metrics appearing**: Check if IoT application is generating metrics

### Debug Commands

```bash
# Check container status
docker ps

# View monitoring logs
make observability-logs

# Check Prometheus targets
curl http://localhost:9090/api/v1/targets

# Check metrics endpoint
curl http://localhost:3001/metrics
```

## Production Considerations

- Change default passwords for Grafana
- Configure proper authentication for monitoring services
- Set up persistent volumes for data storage
- Configure proper alerting channels
- Use dedicated monitoring infrastructure for high-volume deployments

## Support

For issues with the observability stack, see the [Observability Guide](OBSERVABILITY.md) for detailed troubleshooting steps.
