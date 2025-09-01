package main

import (
	"context"
	"log"
	"os"

	"github.com/nhan1603/IoTsystem/api/internal/controller/simulator"
	"github.com/nhan1603/IoTsystem/api/internal/pkg/kafka"
	"github.com/nhan1603/IoTsystem/api/internal/repository"
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

	log.Println("Connected to database successfully")

	// Initial producer kafka
	producer, err := kafka.NewSyncProducer(ctx, os.Getenv("KAFKA_BROKER"))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Kafka successfully")

	// Initial Simulate
	ctrl := simulator.New(
		repo,
		producer,
		os.Getenv("IOT_TOPIC"),
	)
	ctrl.Simulate(ctx)
	select {}
}
