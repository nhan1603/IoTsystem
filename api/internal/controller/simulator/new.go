package simulator

import (
	"context"

	"github.com/nhan1603/IoTsystem/api/internal/pkg/kafka"
	"github.com/nhan1603/IoTsystem/api/internal/repository"
)

type Controller interface {
	Simulate(ctx context.Context)
}

// New initializes a new Controller instance and returns it
func New(
	repo repository.Registry,
	producer *kafka.SyncProducer,
	topic string,
) Controller {
	return impl{
		repo:     repo,
		producer: producer,
		topic:    topic,
	}
}

type impl struct {
	repo     repository.Registry
	producer *kafka.SyncProducer
	topic    string
}
