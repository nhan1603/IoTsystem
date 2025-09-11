# ----------------------------
# Env Variables
# ----------------------------
DOCKER_COMPOSE_FILE ?= build/docker-compose.local.yaml
DOCKER_COMPOSE_OBSERVABILITY ?= build/docker-compose.observability.yaml
DATABASE_CONTAINER ?= database
CASSANDRA_CONTAINER ?= cassandra
CASS_CONTAINTER_NAME ?= cass1
API_CONTAINER ?= server
PROJECT_NAME ?= iotsystem

build-local-go-image:
	docker build -f build/local.go.Dockerfile -t ${PROJECT_NAME}-go-local:latest .
	-docker images -q -f "dangling=true" | xargs docker rmi -f

## run: starts containers to run api server
run: api-create

## setup: executes pre-defined steps to setup api server
setup:
	docker image inspect ${PROJECT_NAME}-go-local:latest >/dev/null 2>&1 || make build-local-go-image
setup: pg-create pg-migrate

## api-create: starts api server
api-create:
	@echo Starting Api container
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} up ${API_CONTAINER}
	@echo Api container started!

## api-gen-models: executes CLI command to generate new database models
api-gen-models:
	@echo Starting generate db model...
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} run -T --rm --service-ports -w /app server sh -c 'sqlboiler --wipe psql && GOFLAGS="-mod=vendor" goimports -w internal/repository/dbmodel/*.go'
	@echo Done!

## pg-create: starts postgres container
pg-create:
	@echo Starting Postgres database container
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} up -d ${DATABASE_CONTAINER}
	@echo Database container started!

## pg-migrate: executes latest migration files
pg-migrate:
	@echo Running migration
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} --profile tools run --rm migrate up
	@echo Migration done!

## api-gen-mocks: generates mock files for testing purpose
api-gen-mocks:
	@echo Starting generate Mock files...
	docker compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} run --name mockery --rm -w /api --entrypoint '' mockery /bin/sh -c "\
			mockery --dir internal/controller --all --recursive --inpackage && \
			mockery --dir internal/repository --all --recursive --inpackage"
	@echo Done!

## test: executes all test cases
test:
	cd api; \
	env $$(grep '^PG_URL=' ./local.env) \
	sh -c 'go test -mod=vendor -p 1 -coverprofile=c.out -failfast -timeout 5m ./... | grep -v pkg'

## pg-drop: reset db to blank
pg-drop:
	@echo Dropping database...
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} --profile tools run --rm migrate down
	@echo Done!

down:
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} down -v
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} rm --force --stop -v

# ----------------------------
# Base infrastructure
# ----------------------------
zookeeper:
	docker compose -f ${DOCKER_COMPOSE_FILE} up -d zookeeper

kafka:
	docker compose -f ${DOCKER_COMPOSE_FILE} up -d kafka

kafka-topic:
	docker compose -f ${DOCKER_COMPOSE_FILE} up -d kafka-topic

infra: zookeeper kafka kafka-topic

# ----------------------------
# Alternate DB
# ----------------------------
## cass-create: starts Cassandra container
cass-wait:
	@echo "Waiting for Cassandra to be ready..."
	@powershell -NoProfile -Command "$$attempts = 60; while ($$attempts -gt 0) { try { $$null = docker-compose -f '${DOCKER_COMPOSE_FILE}' -p '${PROJECT_NAME}' exec -T ${CASSANDRA_CONTAINER} cqlsh -e 'SELECT NOW() FROM system.local;'; if ($$?) { Write-Host 'Cassandra is ready!'; exit 0; }} catch { Write-Host \"Waiting for Cassandra... $$attempts attempts remaining...\"; Start-Sleep -Seconds 1; $$attempts--; }}; if ($$attempts -eq 0) { Write-Host 'Cassandra failed to start!'; exit 1; }"

cass-create:
	@echo Starting Cassandra database container
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} up -d ${CASSANDRA_CONTAINER}
	@echo Cassandra container started!

cass-copy-scripts:
	@echo "Copying CQL scripts to container..."
	docker cp api/data/cassandra/0001_data.up.cql ${CASS_CONTAINTER_NAME}:/schema.cql
	docker cp api/data/cassandra/0002_seed_data.up.cql ${CASS_CONTAINTER_NAME}:/seed.cql
	@echo "Scripts copied successfully!"

## cass-migrate: executes Cassandra schema migrations
cass-migrate: cass-wait cass-copy-scripts
	@echo "Applying Cassandra schema..."
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} exec -T ${CASSANDRA_CONTAINER} cqlsh -f /schema.cql
	@echo "Schema applied successfully!"

cass-cleanup-scripts:
	@echo "Cleaning up CQL scripts from container..."
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} exec -T ${CASSANDRA_CONTAINER} rm -f /schema.cql /seed.cql
	@echo "Scripts cleaned up successfully!"

## cass-seed: seeds initial data into Cassandra
cass-seed: cass-wait cass-copy-scripts
	@echo "Seeding Cassandra data..."
	docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} exec -T ${CASSANDRA_CONTAINER} cqlsh -f /seed.cql
	@make cass-cleanup-scripts
	@echo "Data seeded successfully!"

cass-verify:
	@echo "Verifying Cassandra migration..."
	@docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} exec -T ${CASSANDRA_CONTAINER} cqlsh -e "\
		SELECT COUNT(*) FROM iotsystem.iot_devices; \
		SELECT device_id, floor_id, zone_id FROM iotsystem.iot_devices;"
	@echo "Verification complete!"

## cass-setup: complete Cassandra setup (create, migrate, seed)
cass-setup: cass-create cass-wait cass-drop cass-migrate cass-seed
	@echo "Cassandra setup completed!"

cass-reset: cass-drop cass-migrate cass-seed cass-verify
	@echo "Cassandra reset completed!"

cass-drop:
	@echo "Dropping Cassandra keyspace..."
	@docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} exec -T ${CASSANDRA_CONTAINER} cqlsh -e "\
		DROP KEYSPACE IF EXISTS iotsystem;"
	@echo "Keyspace dropped successfully!"

cass-testing:
	@echo "Counting Cassandra migration..."
	@docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} exec -T ${CASSANDRA_CONTAINER} cqlsh -e "\
		SELECT COUNT(*) FROM iotsystem.sensor_readings; \
		SELECT processed_records FROM iotsystem.benchmark_metrics;"
	@echo "Counting complete!"


# Get stats for a specific device (usage: make cass-device-stats DEVICE=TEMP_001)
cass-device-stats:
	@echo "Getting stats for device ${DEVICE}..."
	@docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} exec -T ${CASSANDRA_CONTAINER} cqlsh -e "\
		SELECT \
			device_id, \
			COUNT(1) as record_count, \
			MIN(timestamp) as first_reading, \
			MAX(timestamp) as last_reading \
		FROM iotsystem.sensor_readings \
		WHERE device_id='${DEVICE}';"

# Loop through all devices (PowerShell script)
cass-all-device-stats:
	@echo "Analyzing all devices..."
	@powershell -NoProfile -Command "\
		$$devices = @('TEMP_001', 'HUM_001', 'CO2_001', 'MULTI_001', 'TEMP_002', 'HUM_101', 'HUM_002', 'CO2_002', 'MULTI_002', 'CO2_003'); \
		foreach ($$device in $$devices) { \
			Write-Host \"\nAnalyzing $$device...\"; \
			make cass-device-stats DEVICE=$$device \
		}"
	@echo "Analysis complete!"

cass-device-stats-mini:
	@echo "Getting stats for device ${DEVICE}..."
	@docker-compose -f ${DOCKER_COMPOSE_FILE} -p=${PROJECT_NAME} exec -T ${CASSANDRA_CONTAINER} cqlsh -e "\
		SELECT \
			COUNT(1) as record_count, \
			MIN(timestamp) as first_reading, \
			MAX(timestamp) as last_reading \
		FROM iotsystem.sensor_readings;"

# ----------------------------
# simulator
# ----------------------------
api-simulator-run:
	docker compose -f ${DOCKER_COMPOSE_FILE} up -d simulator

# ----------------------------
# Observability
# ----------------------------
observability-up:
	@echo Starting observability stack...
	docker-compose -f ${DOCKER_COMPOSE_OBSERVABILITY} -p=${PROJECT_NAME} up -d
	@echo Observability stack started!

observability-down:
	@echo Stopping observability stack...
	docker-compose -f ${DOCKER_COMPOSE_OBSERVABILITY} -p=${PROJECT_NAME} down
	@echo Observability stack stopped!

observability-logs:
	docker-compose -f ${DOCKER_COMPOSE_OBSERVABILITY} -p=${PROJECT_NAME} logs -f

# Start observability with database exporter (requires main system database)
observability-with-db:
	@echo Starting observability stack with database exporter...
	docker-compose -f ${DOCKER_COMPOSE_OBSERVABILITY} -p=${PROJECT_NAME} --profile database up -d
	@echo Observability stack with database exporter started!

# Start observability with Kafka exporter (requires main system Kafka)
observability-with-kafka:
	@echo Starting observability stack with Kafka exporter...
	docker-compose -f ${DOCKER_COMPOSE_OBSERVABILITY} -p=${PROJECT_NAME} --profile kafka up -d
	@echo Observability stack with Kafka exporter started!

# Start observability with both database and Kafka exporters
observability-with-all:
	@echo Starting observability stack with all exporters...
	docker-compose -f ${DOCKER_COMPOSE_OBSERVABILITY} -p=${PROJECT_NAME} --profile database --profile kafka up -d
	@echo Observability stack with all exporters started!

# Start full system with observability
run-with-monitoring: setup observability-up api-create

# Start full system with observability and database exporter
run-with-monitoring-db: setup observability-with-db api-create

# Start full system with observability and all exporters
run-with-monitoring-all: setup observability-with-all api-create

# Start only monitoring services
monitoring-only: observability-up

# Stop monitoring services
monitoring-stop: observability-down

# View monitoring URLs
monitoring-urls:
	@echo "=== IoT System Monitoring URLs ==="
	@echo "Grafana Dashboard: http://localhost:3000 (admin/admin)"
	@echo "Prometheus: http://localhost:9090"
	@echo "cAdvisor (Container Metrics): http://localhost:8080"
	@echo "Node Exporter (Host Metrics): http://localhost:9100/metrics"
	@echo "IoT API Metrics: http://localhost:3001/metrics"
	@echo ""
	@echo "Optional Exporters (when enabled):"
	@echo "PostgreSQL Exporter: http://localhost:9187/metrics"
	@echo "Kafka Exporter: http://localhost:9308/metrics"
	@echo "=================================="