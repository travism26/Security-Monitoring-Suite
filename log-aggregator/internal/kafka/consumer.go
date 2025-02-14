package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/travism26/log-aggregator/internal/domain"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type Consumer struct {
	consumer          sarama.Consumer
	topic             string
	logService        LogService
	alertService      AlertService
	processRepository ProcessRepository // rename to Service
}

func NewConsumer(brokers []string, groupID, topic string, logService LogService, alertService AlertService, processRepo ProcessRepository) (*Consumer, error) {
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
		alertService:      alertService,
		processRepository: processRepo,
	}, nil
}

func (c *Consumer) processMessage(msg *sarama.ConsumerMessage) error {
	// Log the raw message for debugging
	log.Printf("[DEBUG] Raw message received: %s", string(msg.Value))

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

	// Process metrics for alerts
	if err := c.alertService.ProcessMetrics(logEntry); err != nil {
		return fmt.Errorf("failed to process metrics for alerts: %w", err)
	}

	log.Printf("Successfully processed message from topic '%s', partition: %d, offset: %d",
		msg.Topic, msg.Partition, msg.Offset)
	return nil
}

func (c *Consumer) unmarshalRawMessage(msgValue []byte) (*struct {
	Host             interface{} `json:"host"`
	Metrics          interface{} `json:"metrics"`
	ThreatIndicators interface{} `json:"threat_indicators"`
	Metadata         interface{} `json:"metadata"`
	Processes        interface{} `json:"processes"`
	TenantID         string      `json:"tenant_id"`
	APIKey           string      `json:"api_key"`
}, error) {
	var rawMsg struct {
		Host             interface{} `json:"host"`
		Metrics          interface{} `json:"metrics"`
		ThreatIndicators interface{} `json:"threat_indicators"`
		Metadata         interface{} `json:"metadata"`
		Processes        interface{} `json:"processes"`
		TenantID         string      `json:"tenant_id"`
		APIKey           string      `json:"api_key"`
	}
	if err := json.Unmarshal(msgValue, &rawMsg); err != nil {
		return nil, err
	}
	return &rawMsg, nil
}

func (c *Consumer) createLogEntry(rawMsg *struct {
	Host             interface{} `json:"host"`
	Metrics          interface{} `json:"metrics"`
	ThreatIndicators interface{} `json:"threat_indicators"`
	Metadata         interface{} `json:"metadata"`
	Processes        interface{} `json:"processes"`
	TenantID         string      `json:"tenant_id"`
	APIKey           string      `json:"api_key"`
}) (*domain.Log, error) {
	// Debug log for host data
	log.Printf("[DEBUG] Host data type: %T", rawMsg.Host)
	log.Printf("[DEBUG] Host data value: %+v", rawMsg.Host)

	hostInfo, ok := rawMsg.Host.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid host format: expected map[string]interface{}, got %T", rawMsg.Host)
	}

	// Debug log for hostname field
	log.Printf("[DEBUG] Hostname field type: %T", hostInfo["hostname"])
	log.Printf("[DEBUG] Hostname field value: %+v", hostInfo["hostname"])

	hostname, ok := hostInfo["hostname"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid hostname format: expected string, got %T", hostInfo["hostname"])
	}

	metrics, ok := rawMsg.Metrics.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid metrics format")
	}

	cpuUsage, ok := metrics["cpu_usage"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid cpu_usage format")
	}

	memoryUsagePercent, ok := metrics["memory_usage_percent"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid memory_usage_percent format")
	}

	tenantID := rawMsg.TenantID
	if tenantID == "" {
		// Set a default UUID for the default tenant
		tenantID = "00000000-0000-4000-a000-000000000000"
	}

	logEntry := &domain.Log{
		ID:             uuid.New().String(),
		OrganizationID: tenantID,
		APIKey:         rawMsg.APIKey,
		Timestamp:      time.Now(),
		Host:           hostname,
		Message:        fmt.Sprintf("CPU Usage: %.2f%%, Memory Usage: %.2f%%", cpuUsage, memoryUsagePercent),
		Level:          "INFO",
	}

	// Handle processes data if available
	if rawMsg.Processes != nil {
		if processes, ok := rawMsg.Processes.(map[string]interface{}); ok {
			if totalCount, ok := processes["total_count"].(float64); ok {
				logEntry.ProcessCount = int(totalCount)
			}
			if totalCPU, ok := processes["total_cpu_percent"].(float64); ok {
				logEntry.TotalCPUPercent = totalCPU
			}
			if totalMemory, ok := processes["total_memory_usage"].(float64); ok {
				logEntry.TotalMemoryUsage = int64(totalMemory)
			}
		}
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
	Host             interface{} `json:"host"`
	Metrics          interface{} `json:"metrics"`
	ThreatIndicators interface{} `json:"threat_indicators"`
	Metadata         interface{} `json:"metadata"`
	Processes        interface{} `json:"processes"`
	TenantID         string      `json:"tenant_id"`
	APIKey           string      `json:"api_key"`
}, logID string) ([]domain.Process, error) {
	// Handle case where Processes is null
	if rawMsg.Processes == nil {
		log.Printf("[DEBUG] Processes data is nil")
		return []domain.Process{}, nil
	}

	// Debug log for processes data structure
	processesBytes, _ := json.MarshalIndent(rawMsg.Processes, "", "  ")
	log.Printf("[DEBUG] Raw processes data structure:\n%s", string(processesBytes))

	processesData, ok := rawMsg.Processes.(map[string]interface{})
	if !ok {
		log.Printf("[ERROR] Failed to cast processes data to map[string]interface{}, got type: %T", rawMsg.Processes)
		return nil, fmt.Errorf("invalid processes data format")
	}

	processList, ok := processesData["list"].([]interface{})
	if !ok {
		log.Printf("[DEBUG] Process list is nil or not an array. Raw data: %+v", processesData)
		return []domain.Process{}, nil
	}

	log.Printf("[DEBUG] Found %d processes in list", len(processList))

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
