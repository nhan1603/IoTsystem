package simulator

import (
	"context"
	"log"
	"strconv"
	"sync"
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
	log.Printf("[Simulator] Loaded %d devices for simulation\n", len(listDevices))

	// Create context with timeout
	ctxCancel, cancel := context.WithTimeout(ctx, 3*time.Minute)
	defer cancel()

	// Create done channel for cleanup
	done := make(chan struct{})

	go func() {
		executeSensorSimulation(ctxCancel, listDevices, sensorInterval*time.Second, i.topic, i.producer)
		close(done)
	}()

	// Wait for either context timeout or done signal
	select {
	case <-ctx.Done():
		log.Println("[Simulator] Stopping after 5 minutes runtime")
	case <-done:
		log.Println("[Simulator] Simulation completed")
	}
}

// executeSensorSimulation sends random sensor data at intervals for each device
func executeSensorSimulation2(ctx context.Context, listDevices []model.IoTDevice, interval time.Duration, topic string, producer *kafka.SyncProducer) {
	var sentCount int64
	batchSize, _ := strconv.Atoi(env.GetwithDefault("PRODUCER_RATE", "100"))
	generate := func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			for i := 0; i < batchSize; i++ {
				// Pick a random device for each message
				device := listDevices[i%len(listDevices)]
				reading := generateSensorReading(device)
				// b, _ := json.Marshal(reading)
				// fmt.Println(string(b))
				_ = sendMessage(ctx, reading, topic, producer)
				newCount := atomic.AddInt64(&sentCount, 1)
				if newCount%1000 == 0 {
					log.Printf("[Simulator] Total messages sent: %d", newCount)
				}
			}
		}
	}
	for i := 0; i < len(listDevices); i++ {
		go generate()
	}
}

func executeSensorSimulation(ctx context.Context, listDevices []model.IoTDevice, interval time.Duration, topic string, producer *kafka.SyncProducer) {
	var sentCount int64
	batchSize, _ := strconv.Atoi(env.GetwithDefault("PRODUCER_RATE", "100"))

	var wg sync.WaitGroup

	generate := func() {
		defer wg.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for i := 0; i < batchSize; i++ {
					device := listDevices[i%len(listDevices)]
					reading := generateSensorReading(device)
					if err := sendMessage(ctx, reading, topic, producer); err != nil {
						log.Printf("[Simulator] Error sending message: %v", err)
						continue
					}
					newCount := atomic.AddInt64(&sentCount, 1)
					if newCount%1000 == 0 {
						log.Printf("[Simulator] Total messages sent: %d", newCount)
					}
				}
			}
		}
	}

	for i := 0; i < len(listDevices); i++ {
		wg.Add(1)
		go generate()
	}

	wg.Wait() // Wait for all goroutines to finish

	log.Printf("[Simulator] Total messages sent: %d", sentCount)
}
