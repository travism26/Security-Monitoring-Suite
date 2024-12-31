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

func (c *Consumer) processMessage(msg *sarama.ConsumerMessage) error {
	// Log the raw message
	log.Printf("Raw message received: %s", string(msg.Value))

	rawMsg, err := c.unmarshalRawMessage(msg.Value)
	if err != nil {
		return fmt.Errorf("failed to unmarshal message: %w", err)
	}

	logEntry, err := c.createLogEntry(rawMsg)
	if err != nil {
		return fmt.Errorf("failed to create log entry: %w", err)
	}

	processes, err := c.extractProcesses(rawMsg, logEntry.ID)
	if err != nil {
		return fmt.Errorf("failed to extract processes: %w", err)
	}

	if err := c.storeData(logEntry, processes); err != nil {
		return fmt.Errorf("failed to store data: %w", err)
	}

	log.Printf("Successfully processed message from topic '%s', partition: %d, offset: %d",
		msg.Topic, msg.Partition, msg.Offset)
	return nil
}

func (c *Consumer) unmarshalRawMessage(msgValue []byte) (*struct {
	HostInfo         interface{} `json:"host_info"`
	Metrics          interface{} `json:"metrics"`
	ThreatIndicators interface{} `json:"threat_indicators"`
	Metadata         interface{} `json:"metadata"`
	Processes        interface{} `json:"processes"`
}, error) {
	var rawMsg struct {
		HostInfo         interface{} `json:"host_info"`
		Metrics          interface{} `json:"metrics"`
		ThreatIndicators interface{} `json:"threat_indicators"`
		Metadata         interface{} `json:"metadata"`
		Processes        interface{} `json:"processes"`
	}
	if err := json.Unmarshal(msgValue, &rawMsg); err != nil {
		return nil, err
	}
	return &rawMsg, nil
}

func (c *Consumer) createLogEntry(rawMsg *struct {
	HostInfo         interface{} `json:"host_info"`
	Metrics          interface{} `json:"metrics"`
	ThreatIndicators interface{} `json:"threat_indicators"`
	Metadata         interface{} `json:"metadata"`
	Processes        interface{} `json:"processes"`
}) (*domain.Log, error) {
	processes := rawMsg.Processes.(map[string]interface{})

	logEntry := &domain.Log{
		ID:        uuid.New().String(),
		Timestamp: time.Now(),
		Host:      rawMsg.HostInfo.(map[string]interface{})["hostname"].(string),
		Message: fmt.Sprintf("CPU Usage: %.2f%%, Memory Usage: %.2f%%",
			rawMsg.Metrics.(map[string]interface{})["cpu_usage"].(float64),
			rawMsg.Metrics.(map[string]interface{})["memory_usage_percent"].(float64)),
		Level:            "INFO",
		ProcessCount:     int(processes["total_count"].(float64)),
		TotalCPUPercent:  processes["total_cpu_percent"].(float64),
		TotalMemoryUsage: int64(processes["total_memory_usage"].(float64)),
	}

	if rawMsg.Metadata != nil {
		metadataJSON, err := json.Marshal(rawMsg.Metadata)
		if err != nil {
			return nil, fmt.Errorf("error marshaling metadata: %w", err)
		}
		logEntry.MetadataStr = string(metadataJSON)
	}

	return logEntry, nil
}

func (c *Consumer) extractProcesses(rawMsg *struct {
	HostInfo         interface{} `json:"host_info"`
	Metrics          interface{} `json:"metrics"`
	ThreatIndicators interface{} `json:"threat_indicators"`
	Metadata         interface{} `json:"metadata"`
	Processes        interface{} `json:"processes"`
}, logID string) ([]domain.Process, error) {
	processesData := rawMsg.Processes.(map[string]interface{})
	processList := processesData["process_list"].([]interface{})

	processes := make([]domain.Process, 0, len(processList))
	for _, p := range processList {
		proc := p.(map[string]interface{})
		processes = append(processes, domain.Process{
			ID:          uuid.New().String(),
			LogID:       logID,
			Name:        proc["name"].(string),
			PID:         int(proc["pid"].(float64)),
			CPUPercent:  proc["cpu_percent"].(float64),
			MemoryUsage: int64(proc["memory_usage"].(float64)),
			Status:      proc["status"].(string),
			Timestamp:   time.Now(),
		})
	}
	return processes, nil
}

func (c *Consumer) storeData(logEntry *domain.Log, processes []domain.Process) error {
	if err := c.logService.StoreLog(logEntry); err != nil {
		return fmt.Errorf("error storing log entry: %w", err)
	}

	if err := c.processRepository.StoreBatch(processes); err != nil {
		return fmt.Errorf("error storing processes: %w", err)
	}

	return nil
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
				if err := c.processMessage(msg); err != nil {
					log.Printf("Error processing message: %v", err)
				}
			}
		}(pc)
	}

	return nil
}
