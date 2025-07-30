package simulator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nhan1603/IoTsystem/api/internal/model"
	"github.com/nhan1603/IoTsystem/api/internal/pkg/kafka"
)

// sendMessage send the message to kafka
func sendMessage(ctx context.Context, message model.IoTDataMessage, topic string, producer *kafka.SyncProducer) error {
	b, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal input failed: %w", err)
	}

	log.Printf("Sending IoT sensor data to kafka: (topic: %s, payload: %s)\n", topic, string(b))
	_, _, err = producer.SendMessage(ctx, topic, b, kafka.ProducerMessageOption{})
	if err != nil {
		return err
	}

	return nil
}

// generateSensorReading generates a random sensor reading for a device
func generateSensorReading(device model.IoTDevice) model.IoTDataMessage {
	return model.IoTDataMessage{
		DeviceID:    device.DeviceID,
		DeviceName:  device.Name,
		DeviceType:  device.Type,
		Location:    device.Location,
		Floor:       device.Floor,
		Zone:        device.Zone,
		Temperature: 20 + rand.Float64()*10,   // 20-30°C
		Humidity:    30 + rand.Float64()*40,   // 30-70%
		CO2:         400 + rand.Float64()*200, // 400-600 ppm
		Timestamp:   time.Now(),
	}
}
