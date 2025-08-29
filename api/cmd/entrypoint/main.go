package main

import (
	"context"
	"database/sql"
	"log"
	"os"

	"github.com/nhan1603/IoTsystem/api/internal/appconfig/db/pg"
	"github.com/nhan1603/IoTsystem/api/internal/appconfig/httpserver"
	"github.com/nhan1603/IoTsystem/api/internal/appconfig/iam"
	"github.com/nhan1603/IoTsystem/api/internal/controller/auth"
	"github.com/nhan1603/IoTsystem/api/internal/controller/iot"
	"github.com/nhan1603/IoTsystem/api/internal/pkg/obsmetrics"
	"github.com/nhan1603/IoTsystem/api/internal/repository"
	"github.com/nhan1603/IoTsystem/api/internal/repository/generator"
)

func main() {
	log.Println("IOT Ingestion system")
	ctx := context.Background()

	iamConfig := iam.NewConfig()
	ctx = iam.SetConfigToContext(ctx, iamConfig)

	// Initialize Prometheus metrics
	promMetrics := obsmetrics.NewMetrics("iotsystem")

	// Start metrics server
	if err := promMetrics.StartMetricsServer(ctx, ":9091"); err != nil {
		log.Printf("Failed to start metrics server: %v", err)
	} else {
		log.Println("Prometheus metrics server started on :9091")
	}

	// Initial DB connection
	conn, err := pg.Connect(os.Getenv("PG_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Initial snowflake generator
	if err := generator.InitSnowflakeGenerators(); err != nil {
		log.Fatal(err)
	}

	// Initial router
	rtr, err := initRouter(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}

	// start the kafka consumer
	rtr.initKafkaConsumer()

	// Start server
	httpserver.Start(httpserver.Handler(ctx, rtr.routes))
}

func initRouter(
	ctx context.Context,
	db *sql.DB,
) (router, error) {
	repo := repository.New(db)

	return router{
		ctx:      ctx,
		authCtrl: auth.New(repo, iam.ConfigFromContext(ctx)),
		iotCtrol: iot.New(repo),
	}, nil
}
