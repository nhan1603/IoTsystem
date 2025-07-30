package simulator

import (
	"context"
	"strconv"
	"time"

	"github.com/nhan1603/IoTsystem/api/internal/model"
	"github.com/nhan1603/IoTsystem/api/internal/pkg/env"
	"github.com/nhan1603/IoTsystem/api/internal/pkg/kafka"
)

const sensorInterval = 1 // seconds

// Simulate simulates IoT sensor data
func (i impl) Simulate(ctx context.Context) {
	listDevices, err := i.repo.IoT().GetDevices(ctx)
	if err != nil {
		return
	}
	executeSensorSimulation(ctx, listDevices, sensorInterval*time.Second, i.topic, i.producer)

	select {}
}

// executeSensorSimulation sends random sensor data at intervals for each device
func executeSensorSimulation(ctx context.Context, listDevices []model.IoTDevice, interval time.Duration, topic string, producer *kafka.SyncProducer) {
	ticker := time.NewTicker(interval)
	generate := func() {
		for range ticker.C {
			for _, device := range listDevices {
				reading := generateSensorReading(device)
				_ = sendMessage(ctx, reading, topic, producer)
			}
		}
	}
	batchSize, _ := strconv.Atoi(env.GetwithDefault("BATCH_SIZE", "100"))
	for range batchSize {
		go generate()
	}
}
