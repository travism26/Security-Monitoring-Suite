package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/travism26/log-aggregator/internal/domain"
	"github.com/travism26/log-aggregator/internal/service"

	"github.com/IBM/sarama"
)

type Consumer struct {
	consumer   sarama.ConsumerGroup
	logService *service.LogService
	topic      string
}

func NewConsumer(brokers []string, groupID string, topic string, logService *service.LogService) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:   consumer,
		logService: logService,
		topic:      topic,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	topics := []string{c.topic}
	handler := &ConsumerHandler{logService: c.logService}

	for {
		err := c.consumer.Consume(ctx, topics, handler)
		if err != nil {
			log.Printf("Error from consumer: %v", err)
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
	}
}

type ConsumerHandler struct {
	logService *service.LogService
}

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		var log domain.Log
		if err := json.Unmarshal(message.Value, &log); err != nil {
			continue
		}

		if err := h.logService.StoreLog(&log); err != nil {
			continue
		}

		session.MarkMessage(message, "")
	}
	return nil
}

// Required ConsumerGroupHandler interface methods
func (h *ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
