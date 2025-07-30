FROM --platform=linux/amd64 golang:1.23.8-alpine3.21 AS base
RUN apk --no-cache add \
    bash \
    build-base \
    git 

#################

FROM base AS builder

WORKDIR /IoTsystem/api

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    GOOS=linux GOARCH=amd64 go build -o /IoTsystem/api/cmd/entrypoint ./cmd/entrypoint

###################

FROM --platform=linux/amd64 alpine:3.19
RUN apk --no-cache add \
    ca-certificates \
    tzdata
COPY --from=builder /IoTsystem/api/cmd/entrypoint /
COPY ./templates ./templates

RUN adduser -D -H -u 1000 IoTsystem
USER IoTsystem

EXPOSE 3001
CMD /entrypoint
