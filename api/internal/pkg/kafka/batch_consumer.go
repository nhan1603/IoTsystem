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
	log.Println("[Kafka Consumer] Starting to consume messages...")
	timer := time.NewTimer(h.batchTimeout)
	defer timer.Stop()

	resetTimer := func() {
		if !timer.Stop() {
			select {
			case <-timer.C:
			default:
			}
		}
		timer.Reset(h.batchTimeout)
	}

	commitSuccesses := func(succeeded []ConsumerMessage) {
		for _, m := range succeeded {
			// Mark only the offsets that truly succeeded
			sess.MarkOffset(m.ID.Topic, m.ID.Partition, m.ID.Offset+1, "")
		}
	}

	processAndCommit := func(b []ConsumerMessage) {
		if len(b) == 0 {
			return
		}

		// respect rebalance/cancel
		ctx, cancel := context.WithTimeout(sess.Context(), h.batchTimeout) // e.g., 5–10s
		defer cancel()

		log.Printf("[Kafka Consumer] Processing batch of %d messages...\n", len(b))

		// Try whole batch with retry
		if err := h.processBatch(ctx, b); err == nil {
			commitSuccesses(b)
			sess.Commit()
			return
		}

		// Fallback: bisect to isolate bad records
		succeeded, dropped := h.processWithBisection(sess.Context(), b)

		// Commit successful records to make progress
		commitSuccesses(succeeded)

		// For drops: log + metrics + commit so the partition doesn't stall
		for _, d := range dropped {
			log.Printf("[Kafka Consumer] Dropping poison message at %s[%d]@%d key=%q",
				d.ID.Topic, d.ID.Partition, d.ID.Offset, d.ID.Key)
			sess.MarkOffset(d.ID.Topic, d.ID.Partition, d.ID.Offset+1, "")
		}

		// flush offsets promptly
		// only commit if we have marked offsets
		// this avoids committing empty batches which can lead to confusion
		if len(succeeded)+len(dropped) > 0 {
			sess.Commit()
		}
	}

	for {
		select {
		case cm, ok := <-claim.Messages():
			if !ok {
				processAndCommit(batch)
				return nil
			}
			// collect message
			batch = append(batch, toConsumerMessage(cm))
			if len(batch) >= h.batchSize {
				processAndCommit(batch)
				batch = batch[:0] // reset batch
				resetTimer()      // reset timer after processing
			}

		case <-timer.C:
			processAndCommit(batch)
			batch = batch[:0]
			resetTimer()
		case <-sess.Context().Done():
			if len(batch) > 0 {
				processAndCommit(batch)
			}
			return nil
		}
	}
}

// Bisection fallback (no need for a single-message handler)
func (h batchMessageHandler) processWithBisection(ctx context.Context, batch []ConsumerMessage) (succeeded, dropped []ConsumerMessage) {
	var ok, bad []ConsumerMessage

	var bisect func([]ConsumerMessage)
	bisect = func(b []ConsumerMessage) {
		if len(b) == 0 {
			return
		}
		// try processing the batch
		if err := h.handler(ctx, b); err == nil {
			ok = append(ok, b...)
			return
		}
		if len(b) == 1 {
			// Cannot be processed even alone → drop
			bad = append(bad, b[0])
			return
		}
		mid := len(b) / 2
		bisect(b[:mid])
		bisect(b[mid:])
	}

	bisect(batch)
	return ok, bad
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
