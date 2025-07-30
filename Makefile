# ----------------------------
# Env Variables
# ----------------------------
DOCKER_COMPOSE_FILE ?= build/docker-compose.local.yaml
DATABASE_CONTAINER ?= database
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

# ----------------------------
# simulator
# ----------------------------
api-simulator-run:
	docker compose -f ${DOCKER_COMPOSE_FILE} up -d simulator