package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/travism26/log-aggregator/internal/domain"
)

const (
	maxRetries = 3
	retryDelay = 100 * time.Millisecond
)

type LogServiceConfig struct {
	Environment string
	Application string
	Component   string
}

type LogService struct {
	repo   domain.LogRepository
	config LogServiceConfig
}

func NewLogService(repo domain.LogRepository, config LogServiceConfig) *LogService {
	return &LogService{
		repo:   repo,
		config: config,
	}
}

// retryOperation executes an operation with retries and exponential backoff
func (s *LogService) retryOperation(operation func() error) error {
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if err := operation(); err != nil {
			lastErr = err
			time.Sleep(retryDelay * time.Duration(attempt+1))
			continue
		}
		return nil
	}
	return fmt.Errorf("operation failed after %d attempts: %w", maxRetries, lastErr)
}

func (s *LogService) StoreLog(log *domain.Log) error {
	// Enrich log with environment information
	log.EnrichLog(s.config.Environment, s.config.Application, s.config.Component)
	log.ProcessedCount++

	if len(log.Metadata) > 0 && log.MetadataStr == "" {
		metadataJSON, err := json.Marshal(log.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		log.MetadataStr = string(metadataJSON)
	}

	return s.retryOperation(func() error {
		return s.repo.Store(log)
	})
}

func (s *LogService) StoreBatch(logs []*domain.Log) error {
	// Process metadata and enrich each log in the batch
	for _, log := range logs {
		// Enrich log with environment information
		log.EnrichLog(s.config.Environment, s.config.Application, s.config.Component)
		log.ProcessedCount++

		if len(log.Metadata) > 0 && log.MetadataStr == "" {
			metadataJSON, err := json.Marshal(log.Metadata)
			if err != nil {
				return fmt.Errorf("failed to marshal metadata for log %s: %w", log.ID, err)
			}
			log.MetadataStr = string(metadataJSON)
		}
	}

	return s.retryOperation(func() error {
		return s.repo.StoreBatch(logs)
	})
}

func (s *LogService) GetLog(id string) (*domain.Log, error) {
	var log *domain.Log
	err := s.retryOperation(func() error {
		var err error
		log, err = s.repo.FindByID(id)
		return err
	})
	return log, err
}

func (s *LogService) ListLogs(limit, offset int) ([]*domain.Log, error) {
	var logs []*domain.Log
	err := s.retryOperation(func() error {
		var err error
		logs, err = s.repo.List(limit, offset)
		return err
	})
	return logs, err
}

func (s *LogService) ListByTimeRange(start, end time.Time, limit, offset int) ([]*domain.Log, error) {
	// Validate time range
	if start.After(end) {
		return nil, fmt.Errorf("invalid time range: start time %v is after end time %v", start, end)
	}

	var logs []*domain.Log
	err := s.retryOperation(func() error {
		var err error
		logs, err = s.repo.ListByTimeRange(start, end, limit, offset)
		return err
	})
	return logs, err
}
