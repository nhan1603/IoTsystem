package kafka

import (
	"context"
	"errors"
	"log"

	"github.com/IBM/sarama"
	"github.com/nhan1603/IoTsystem/api/internal/pkg/uid"
	pkgerrors "github.com/pkg/errors"
)

// SyncProducer publishes Kafka message.
type SyncProducer struct {
	client   sarama.Client
	producer sarama.SyncProducer
}

// NewSyncProducer creates a newsync producer using the given broker addresses and configuration.
func NewSyncProducer(ctx context.Context,
	broker string) (*SyncProducer, error) {
	log.Println("Initializing Kafka SyncProducer")

	baseCfg := sarama.NewConfig()
	cfg := &producerConfig{baseCfg}
	cfg.Producer.Return.Successes = true // Mandatory forsync producer

	client, err := sarama.NewClient([]string{broker}, cfg.Config)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "client init failed")
	}

	p, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "starting sync producer")
	}

	sp := &SyncProducer{
		client:   client,
		producer: p,
	}

	return sp, nil
}

// SendMessage sends a given message and returns err if encountered error else the partition & offset of the message
func (p *SyncProducer) SendMessage(ctx context.Context, topic string, payload []byte, opt ProducerMessageOption) (int32, int64, error) {
	pm, err := prepareProducerMessage(topic, payload, opt)
	if err != nil {
		return 0, 0, err
	}

	var partition int32
	var offset int64

	partition, offset, err = p.producer.SendMessage(pm)
	if err != nil {
		return 0, 0, pkgerrors.WithStack(err)
	}

	return partition, offset, nil
}

// Close shuts down the producer.
func (p *SyncProducer) Close() error {
	if err := p.producer.Close(); err != nil {
		return pkgerrors.Wrap(err, "could not stop producer")
	}
	if !p.client.Closed() { // Just in case
		if err := p.client.Close(); err != nil {
			return pkgerrors.Wrap(err, "could not stop producer client")
		}
	}
	return nil
}

type producerConfig struct {
	*sarama.Config
}

// ProducerOption overrides the properties of a producer
type ProducerOption func(*producerConfig)

// ProducerMessageOption specifies options for the message
type ProducerMessageOption struct {
	Key                   string
	Partition             *int32 // Cannot use int because of forced conversion
	Headers               map[string]string
	DisablePayloadLogging bool
}

var generateUIDFunc = uid.Generate

func prepareProducerMessage(topic string, payload []byte, opt ProducerMessageOption) (*sarama.ProducerMessage, error) {
	if topic == "" {
		return nil, errors.New("topic is empty")
	}
	if len(payload) == 0 {
		return nil, errors.New("no payload provided")
	}

	if opt.Key == "" {
		id, err := generateUIDFunc()
		if err != nil {
			return nil, err
		}
		opt.Key = id
	}

	pm := &sarama.ProducerMessage{Topic: topic, Value: sarama.ByteEncoder(payload), Key: sarama.StringEncoder(opt.Key)}

	if opt.Partition != nil {
		pm.Partition = *opt.Partition
	}

	if l := len(opt.Headers); l > 0 {
		pm.Headers = make([]sarama.RecordHeader, 0, l)
		for k, v := range opt.Headers {
			pm.Headers = append(pm.Headers, sarama.RecordHeader{
				Key:   []byte(k),
				Value: []byte(v),
			})
		}
	}

	return pm, nil
}
