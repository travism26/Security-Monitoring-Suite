package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/travism26/log-aggregator/internal/domain"
	"github.com/travism26/log-aggregator/internal/service"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
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
				// Log the raw message
				log.Printf("Raw message received: %s", string(msg.Value))

				// First unmarshal into a temporary struct that matches the JSON structure
				var rawMsg struct {
					HostInfo         interface{} `json:"host_info"`
					Metrics          interface{} `json:"metrics"`
					ThreatIndicators interface{} `json:"threat_indicators"`
					Metadata         interface{} `json:"metadata"`
				}

				if err := json.Unmarshal(msg.Value, &rawMsg); err != nil {
					log.Printf("Error unmarshaling message: %v", err)
					continue
				}

				// Create the log entry
				logEntry := domain.Log{
					ID:        uuid.New().String(),
					Timestamp: time.Now(),
					Host:      rawMsg.HostInfo.(map[string]interface{})["hostname"].(string),
					Message: fmt.Sprintf("CPU Usage: %.2f%%, Memory Usage: %.2f%%",
						rawMsg.Metrics.(map[string]interface{})["cpu_usage"].(float64),
						rawMsg.Metrics.(map[string]interface{})["memory_usage_percent"].(float64)),
					Level: "INFO",
				}

				// Log initial state
				log.Printf("Initial logEntry state - Metadata: %+v, MetadataStr: %q",
					rawMsg.Metadata, logEntry.MetadataStr)

				// Handle metadata separately
				if rawMsg.Metadata != nil {
					metadataJSON, err := json.Marshal(rawMsg.Metadata)
					if err != nil {
						log.Printf("Error marshaling metadata: %v", err)
						continue
					}
					logEntry.MetadataStr = string(metadataJSON)
					log.Printf("Metadata processed - JSON string: %s", logEntry.MetadataStr)
				}

				// Log final state before storage
				log.Printf("Final logEntry state before storage - Metadata: %+v, MetadataStr: %q",
					logEntry.Metadata, logEntry.MetadataStr)

				// Store the log
				if err := c.logService.StoreLog(&logEntry); err != nil {
					log.Printf("Error storing log: %v", err)
					log.Printf("LogEntry at time of error: %+v", logEntry)
					continue
				}

				log.Printf("Successfully processed message from topic '%s', partition: %d, offset: %d",
					msg.Topic, msg.Partition, msg.Offset)
			}
		}(pc)
	}

	return nil
}
