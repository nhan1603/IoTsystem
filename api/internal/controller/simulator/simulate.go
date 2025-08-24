package simulator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync/atomic"
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
	var sentCount int64
	ticker := time.NewTicker(interval)
	batchSize, _ := strconv.Atoi(env.GetwithDefault("BATCH_SIZE", "100"))
	generate := func() {
		for range ticker.C {
			for i := 0; i < batchSize; i++ {
				// Pick a random device for each message
				device := listDevices[i%len(listDevices)]
				reading := generateSensorReading(device)
				b, _ := json.Marshal(reading)
				fmt.Println(string(b))
				_ = sendMessage(ctx, reading, topic, producer)
				newCount := atomic.AddInt64(&sentCount, 1)
				if newCount%1000 == 0 {
					log.Printf("[Simulator] Total messages sent: %d", newCount)
				}
			}
		}
	}
	for i := 0; i < 2; i++ {
		go generate()
	}
}
