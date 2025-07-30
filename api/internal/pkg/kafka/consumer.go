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

// Consumer is the kafka consumer
type Consumer struct {
	client   sarama.Client
	consumer sarama.ConsumerGroup
	topic    string
	handler  messageHandler
}

// ConsumeHandler handles message consuming.
type ConsumeHandler func(ctx context.Context, msg ConsumerMessage) error

type consumerConfig struct {
	*sarama.Config
	disablePayloadLogging bool
	maxRetriesPerMsg      int
	groupID               string
	batchSize             int
	batchTimeout          time.Duration
}

// NewConsumer creates a new consumer using the given broker addresses and configuration.
func NewConsumer(
	ctx context.Context,
	topic string,
	broker string,
	handler ConsumeHandler,
) (*Consumer, error) {
	log.Printf("Initializing kafka consumer for topic: [%s]", topic)

	if topic == "" {
		return nil, errors.New("topic is empty")
	}

	baseCfg := sarama.NewConfig()

	cfg := &consumerConfig{Config: baseCfg, groupID: "iot"}
	cfg.maxRetriesPerMsg = 35 // This evaluates to around 13hrs with the current backoff config

	client, err := sarama.NewClient([]string{broker}, cfg.Config)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "client init failed")
	}

	cg, err := sarama.NewConsumerGroupFromClient(cfg.groupID, client)
	if err != nil {
		return nil, pkgerrors.Wrap(err, "creating new consumer")
	}

	msgHandler := messageHandler{
		handler:               handler,
		maxRetriesPerMsg:      cfg.maxRetriesPerMsg,
		disablePayloadLogging: cfg.disablePayloadLogging,
	}

	return &Consumer{
		client:   client,
		consumer: cg,
		topic:    topic,
		handler:  msgHandler,
	}, nil
}

// Consume consumes messages in a loop
func (c *Consumer) Consume(ctx context.Context) error {
	consumeErr := make(chan error, 1)
	go func() {
		for {
			if ctx.Err() != nil {
				return
			}
			if err := c.consumer.Consume(ctx, []string{c.topic}, c.handler); err != nil {
				consumeErr <- pkgerrors.WithStack(fmt.Errorf("consuming failed. err: %w", err))
				return
			}
		}
	}()

	select {
	case err := <-consumeErr:
		return err
	case <-ctx.Done():
		log.Printf("[kafka_consumer] Closing consumer group....")
		if err := c.consumer.Close(); err != nil {
			return pkgerrors.Wrap(err, "could not stop consumer")
		}
		if !c.client.Closed() {
			if err := c.client.Close(); err != nil {
				return pkgerrors.Wrap(err, "could not stop consumer client")
			}
		}
		log.Printf("[kafka_consumer] Consumer group closed")
		return nil
	}
}

// ConsumerMessage encapsulates a Kafka message returned by the consumer.
type ConsumerMessage struct {
	ID      ConsumerMessageID
	Value   []byte
	Headers map[string]string
}

// ConsumerMessageID is the unique identifier of the message
type ConsumerMessageID struct {
	Topic     string
	Partition int32 // Cannot use int because of potential force conversion failures
	Offset    int64
	Key       string
}

type messageHandler struct {
	handler               ConsumeHandler
	disablePayloadLogging bool
	maxRetriesPerMsg      int
}

func (h messageHandler) Setup(s sarama.ConsumerGroupSession) error {
	return nil
}

func (h messageHandler) Cleanup(sarama.ConsumerGroupSession) error {
	log.Println("[Kafka Consumer] Cleaning up...")
	return nil
}

func (h messageHandler) ConsumeClaim(s sarama.ConsumerGroupSession, c sarama.ConsumerGroupClaim) error {
	for msg := range c.Messages() {
		h.consume(s, msg)
	}
	return nil
}

func (h messageHandler) consume(s sarama.ConsumerGroupSession, cm *sarama.ConsumerMessage) {
	ctx := context.Background()

	var msgKey string
	if cm.Key != nil {
		msgKey = string(cm.Key)
	}

	var err error
	defer func() {
		if rcv := recover(); rcv != nil {
			err = pkgerrors.WithStack(fmt.Errorf("panic err: %s", rcv))
		}
	}()

	msg := ConsumerMessage{
		ID: ConsumerMessageID{
			Topic:     cm.Topic,
			Partition: cm.Partition,
			Offset:    cm.Offset,
			Key:       msgKey,
		},
		Value:   cm.Value,
		Headers: make(map[string]string, len(cm.Headers)), // It's ok to possibly over provision in case of duplicate.
	}
	if cm.Headers != nil {
		for _, r := range cm.Headers {
			msg.Headers[string(r.Key[:])] = string(r.Value[:])
		}
	}

	var attempts int

	if err = backoff.Retry(func() error {
		attempts++
		log.Printf("[Kafka Consumer] Consuming Attempt: [%d]", attempts)

		if err = h.handler(ctx, msg); err != nil {
			log.Printf("consume message failed, err %v", err)
			return err
		}
		return nil
	}, backoff.WithContext(consumeBackoff(h.maxRetriesPerMsg), ctx)); err != nil {
		log.Printf("[Kafka Consumer] Giving up on processing. Partition: [%d], Offset: [%d] after [%d] attempts. Will just commit and move on", cm.Partition, cm.Offset, attempts)
	}

	h.commitMessageOffset(ctx, s, msg.ID)
}

// CommitMessageOffset commits the message's offset+1 for the topic & partition
func (h messageHandler) commitMessageOffset(
	ctx context.Context,
	cgs sarama.ConsumerGroupSession,
	msgID ConsumerMessageID,
) {
	offsetToCommit := msgID.Offset + 1 // Should always commit next offset as best practice
	cgs.MarkOffset(msgID.Topic, msgID.Partition, offsetToCommit, "")
}

func consumeBackoff(maxRetries int) backoff.BackOff {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = 5 * time.Second
	b.RandomizationFactor = 0
	b.Multiplier = 1.25
	b.MaxInterval = 30 * time.Minute
	b.MaxElapsedTime = 12 * time.Hour
	return backoff.WithMaxRetries(b, uint64(maxRetries))
}
