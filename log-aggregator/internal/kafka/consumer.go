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
	consumer   sarama.Consumer
	topic      string
	logService *service.LogService
}

func NewConsumer(brokers []string, groupID, topic string, logService *service.LogService) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:   consumer,
		topic:      topic,
		logService: logService,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	partitions, err := c.consumer.Partitions(c.topic)
	if err != nil {
		return err
	}

	for _, partition := range partitions {
		pc, err := c.consumer.ConsumePartition(c.topic, partition, sarama.OffsetNewest)
		if err != nil {
			return err
		}

		go func(pc sarama.PartitionConsumer) {
			defer pc.Close()
			for msg := range pc.Messages() {
				var logEntry domain.Log
				if err := json.Unmarshal(msg.Value, &logEntry); err != nil {
					log.Printf("Error unmarshaling message: %v", err)
					continue
				}

				if err := c.logService.StoreLog(&logEntry); err != nil {
					log.Printf("Error storing log: %v", err)
					continue
				}

				log.Printf("Successfully processed message from topic '%s', partition: %d, offset: %d",
					msg.Topic, msg.Partition, msg.Offset)
				log.Printf("Log Entry Details - ID: %s, Host: %s, Level: %s, Timestamp: %s",
					logEntry.ID, logEntry.Host, logEntry.Level, logEntry.Timestamp)

				if len(logEntry.Metadata) > 0 {
					log.Printf("Message contains %d metadata fields", len(logEntry.Metadata))
				}
			}
		}(pc)
	}

	return nil
}
