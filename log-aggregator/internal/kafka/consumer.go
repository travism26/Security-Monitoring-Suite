package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/travism26/log-aggregator/internal/domain"
	"github.com/travism26/log-aggregator/internal/repository/postgres"
	"github.com/travism26/log-aggregator/internal/service"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type Consumer struct {
	consumer          sarama.Consumer
	topic             string
	logService        *service.LogService
	processRepository *postgres.ProcessRepository
}

func NewConsumer(brokers []string, groupID, topic string, logService *service.LogService, processRepo *postgres.ProcessRepository) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		consumer:          consumer,
		topic:             topic,
		logService:        logService,
		processRepository: processRepo,
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
					Processes        interface{} `json:"processes"`
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

				if rawMsg.Metrics != nil {
					metrics := rawMsg.Metrics.(map[string]interface{})
					if processesData, ok := metrics["processes"].(map[string]interface{}); ok {
						if processList, ok := processesData["process_list"].([]interface{}); ok {
							var processes []domain.Process
							for _, p := range processList {
								proc := p.(map[string]interface{})
								process := domain.Process{
									ID:          uuid.New().String(),
									LogID:       logEntry.ID,
									Name:        proc["name"].(string),
									PID:         int(proc["pid"].(float64)),
									CPUPercent:  proc["cpu_percent"].(float64),
									MemoryUsage: int64(proc["memory_usage"].(float64)),
									Status:      proc["status"].(string),
								}
								processes = append(processes, process)
							}

							if err := c.processRepository.StoreBatch(processes); err != nil {
								log.Printf("Error storing processes: %v", err)
							}
						}
					}
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
