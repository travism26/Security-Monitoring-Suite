package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/IBM/sarama"
	"github.com/rickjms/security-monitoring-suite/mini-xdr/internal/domain"
	"go.uber.org/zap"
)

// EventHandler defines the interface for handling normalized events
type EventHandler interface {
	HandleEvent(ctx context.Context, event *domain.Event) error
}

// Consumer represents a Kafka consumer for security events
type Consumer struct {
	client    sarama.ConsumerGroup
	topics    []string
	handler   EventHandler
	logger    *zap.Logger
	ready     chan bool
	closeOnce sync.Once
}

// NewConsumer creates a new Kafka consumer instance
func NewConsumer(brokers []string, groupID string, topics []string, handler EventHandler, logger *zap.Logger) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	return &Consumer{
		client:  client,
		topics:  topics,
		handler: handler,
		logger:  logger,
		ready:   make(chan bool),
	}, nil
}

// Start begins consuming messages from Kafka
func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			err := c.client.Consume(ctx, c.topics, c)
			if err != nil {
				c.logger.Error("Error from consumer", zap.Error(err))
				return err
			}
			c.ready = make(chan bool)
		}
	}
}

// Setup is run at the beginning of a new session
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from a partition
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			event, err := c.normalizeMessage(message)
			if err != nil {
				c.logger.Error("Failed to normalize message",
					zap.Error(err),
					zap.String("topic", message.Topic),
					zap.Int32("partition", message.Partition),
					zap.Int64("offset", message.Offset))
				continue
			}

			if err := c.handler.HandleEvent(session.Context(), event); err != nil {
				c.logger.Error("Failed to handle event",
					zap.Error(err),
					zap.String("event_id", event.ID),
					zap.String("event_type", string(event.Type)))
				continue
			}

			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}

// normalizeMessage converts a Kafka message into a normalized Event
func (c *Consumer) normalizeMessage(msg *sarama.ConsumerMessage) (*domain.Event, error) {
	var rawEvent map[string]interface{}
	if err := json.Unmarshal(msg.Value, &rawEvent); err != nil {
		return nil, fmt.Errorf("failed to unmarshal message: %w", err)
	}

	// Extract common fields with type assertions
	eventType, ok := rawEvent["type"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid event type")
	}

	source, ok := rawEvent["source"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid source")
	}

	// Create normalized event
	event := &domain.Event{
		ID:          fmt.Sprintf("%s-%d-%d", msg.Topic, msg.Partition, msg.Offset),
		Type:        domain.EventType(eventType),
		Source:      source,
		RawData:     msg.Value,
		Timestamp:   msg.Timestamp,
		Description: rawEvent["description"].(string),
		Metadata:    make(map[string]string),
	}

	// Extract severity if present
	if severity, ok := rawEvent["severity"].(string); ok {
		event.Severity = domain.Severity(severity)
	} else {
		event.Severity = domain.SeverityLow
	}

	// Extract metadata fields
	if metadata, ok := rawEvent["metadata"].(map[string]interface{}); ok {
		for k, v := range metadata {
			if strVal, ok := v.(string); ok {
				event.Metadata[k] = strVal
			}
		}
	}

	if err := event.Validate(); err != nil {
		return nil, fmt.Errorf("event validation failed: %w", err)
	}

	return event, nil
}

// Close gracefully shuts down the consumer
func (c *Consumer) Close() error {
	var err error
	c.closeOnce.Do(func() {
		err = c.client.Close()
	})
	return err
}
