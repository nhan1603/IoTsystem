package main

import (
	"context"
	"log"
	"os"

	"github.com/nhan1603/IoTsystem/api/internal/controller/simulator"
	"github.com/nhan1603/IoTsystem/api/internal/pkg/kafka"
	"github.com/nhan1603/IoTsystem/api/internal/repository"
	"github.com/nhan1603/IoTsystem/api/internal/repository/generator"
)

func main() {
	ctx := context.Background()

	// Initial DB connection
	cfg := repository.FromEnv()
	repo, cleanup, err := repository.NewFromConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cleanup()

	// Initial producer kafka
	producer, err := kafka.NewSyncProducer(ctx, os.Getenv("KAFKA_BROKER"))
	if err != nil {
		log.Fatal(err)
	}

	// Initial snowflake generator
	if err := generator.InitSnowflakeGenerators(); err != nil {
		log.Printf("Error when init snowflake, %v", err)
		return
	}

	// Initial Simulate
	ctrl := simulator.New(
		repo,
		producer,
		os.Getenv("IOT_TOPIC"),
	)
	ctrl.Simulate(ctx)
	select {}
}
