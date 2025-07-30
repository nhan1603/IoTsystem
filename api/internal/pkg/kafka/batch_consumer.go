package kafka

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/cenkalti/backoff/v4"
	pkgerrors "github.com/pkg/errors"
)

// BatchConsumer wraps a Kafka consumer with batch processing capabilities.
type BatchConsumer struct {
	client   sarama.Client
	consumer sarama.ConsumerGroup
	topic    string
	handler  batchMessageHandler
}

// BatchConsumeHandler is called with a slice of messages to process as a batch.
type BatchConsumeHandler func(ctx context.Context, msgs []ConsumerMessage) error

// NewBatchConsumer creates a new consumer supporting batch processing.
func NewBatchConsumer(
	ctx context.Context,
	topic string,
	broker string,
	handler BatchConsumeHandler,
	groupID string,
	batchSize int,
	batchTimeout time.Duration,
) (*BatchConsumer, error) {
	if topic == "" {
		return nil, errors.New("topic is empty")
	}

	baseCfg := sarama.NewConfig()
	cfg := &consumerConfig{Config: baseCfg, groupID: groupID}
	cfg.maxRetriesPerMsg = 35
	cfg.batchSize = batchSize
	cfg.batchTimeout = batchTimeout

	client, err := sarama.NewClient([]string{broker}, cfg.Config)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "client init failed")
	}

	cg, err := sarama.NewConsumerGroupFromClient(cfg.groupID, client)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "creating new consumer group")
	}

	h := batchMessageHandler{
		handler:               handler,
		maxRetriesPerMsg:      cfg.maxRetriesPerMsg,
		disablePayloadLogging: cfg.disablePayloadLogging,
		batchSize:             cfg.batchSize,
		batchTimeout:          cfg.batchTimeout,
	}

	return &BatchConsumer{
		client:   client,
		consumer: cg,
		topic:    topic,
		handler:  h,
	}, nil
}

// Consume starts consuming messages in batches.
func (c *BatchConsumer) Consume(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			err := c.consumer.Consume(ctx, []string{c.topic}, c.handler)
			if err != nil {
				errCh <- pkgerrors.WithStack(fmt.Errorf("consuming failed: %w", err))
				return
			}
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		log.Println("[Kafka BatchConsumer] shutting down")
		c.consumer.Close()
		c.client.Close()
		return nil
	}
}

// batchMessageHandler implements sarama.ConsumerGroupHandler
type batchMessageHandler struct {
	handler               BatchConsumeHandler
	disablePayloadLogging bool
	maxRetriesPerMsg      int
	batchSize             int
	batchTimeout          time.Duration
}

func (h batchMessageHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }
func (h batchMessageHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	log.Println("[Kafka Consumer] Cleaning up...")
	return nil
}

func (h batchMessageHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	batch := make([]ConsumerMessage, 0, h.batchSize)
	timer := time.NewTimer(h.batchTimeout)
	defer timer.Stop()

	commitBatch := func() {
		if len(batch) == 0 {
			return
		}
		ctx := context.Background()
		if err := h.processBatch(ctx, batch); err != nil {
			log.Printf("batch processing failed: %v", err)
		}
		// mark offsets
		for _, msg := range batch {
			sess.MarkOffset(msg.ID.Topic, msg.ID.Partition, msg.ID.Offset+1, "")
		}
		batch = batch[:0]
	}

	for {
		select {
		case cm, ok := <-claim.Messages():
			if !ok {
				commitBatch()
				return nil
			}
			// collect message
			batch = append(batch, toConsumerMessage(cm))
			if len(batch) >= h.batchSize {
				commitBatch()
				timer.Reset(h.batchTimeout)
			}

		case <-timer.C:
			commitBatch()
			timer.Reset(h.batchTimeout)
		}
	}
}

func (h batchMessageHandler) processBatch(ctx context.Context, batch []ConsumerMessage) error {
	// wrap in retry
	return backoff.Retry(func() error {
		return h.handler(ctx, batch)
	}, backoff.WithContext(consumeBackoff(h.maxRetriesPerMsg), ctx))
}

func toConsumerMessage(cm *sarama.ConsumerMessage) ConsumerMessage {
	m := ConsumerMessage{
		ID: ConsumerMessageID{
			Topic:     cm.Topic,
			Partition: cm.Partition,
			Offset:    cm.Offset,
			Key:       string(cm.Key),
		},
		Value:   cm.Value,
		Headers: make(map[string]string, len(cm.Headers)),
	}
	for _, h := range cm.Headers {
		m.Headers[string(h.Key)] = string(h.Value)
	}
	return m
}
