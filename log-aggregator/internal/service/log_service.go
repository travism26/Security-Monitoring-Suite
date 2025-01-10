package service

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/travism26/log-aggregator/internal/cache"
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
	Cache       struct {
		Enabled      bool
		TTL          time.Duration
		TimeRangeTTL time.Duration
	}
}

type LogService struct {
	repo         domain.LogRepository
	config       LogServiceConfig
	cache        cache.Cache
	keyGenerator *cache.CacheKeyGenerator
}

func NewLogService(repo domain.LogRepository, config LogServiceConfig) *LogService {
	var cacheInstance cache.Cache
	if config.Cache.Enabled {
		cacheInstance = cache.NewInMemoryCache()
	}

	return &LogService{
		repo:         repo,
		config:       config,
		cache:        cacheInstance,
		keyGenerator: cache.NewCacheKeyGenerator(),
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

	err := s.retryOperation(func() error {
		return s.repo.Store(log)
	})

	// If store was successful and cache is enabled, invalidate related cache entries
	if err == nil && s.cache != nil {
		s.cache.Delete(s.keyGenerator.ForLog(log.ID))
		// Clear list caches as they might be affected
		s.cache.Clear() // TODO: Implement more granular cache invalidation
	}

	return err
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

	err := s.retryOperation(func() error {
		return s.repo.StoreBatch(logs)
	})

	// If store was successful and cache is enabled, invalidate affected cache entries
	if err == nil && s.cache != nil {
		for _, log := range logs {
			s.cache.Delete(s.keyGenerator.ForLog(log.ID))
		}
		// Clear list caches as they might be affected
		s.cache.Clear() // TODO: Implement more granular cache invalidation
	}

	return err
}

func (s *LogService) GetLog(id string) (*domain.Log, error) {
	if s.cache != nil {
		// Try to get from cache first
		if cached, found := s.cache.Get(s.keyGenerator.ForLog(id)); found {
			if log, ok := cached.(*domain.Log); ok {
				return log, nil
			}
		}
	}

	// If not in cache or cache disabled, get from repository
	var log *domain.Log
	err := s.retryOperation(func() error {
		var err error
		log, err = s.repo.FindByID(id)
		if err != nil {
			return err
		}

		// Cache the result if cache is enabled
		if s.cache != nil && log != nil {
			s.cache.Set(s.keyGenerator.ForLog(id), log, s.config.Cache.TTL)
		}
		return nil
	})
	return log, err
}

func (s *LogService) ListLogs(limit, offset int) ([]*domain.Log, error) {
	if s.cache != nil {
		// Try to get from cache first
		if cached, found := s.cache.Get(s.keyGenerator.ForLogList(limit, offset)); found {
			if logs, ok := cached.([]*domain.Log); ok {
				return logs, nil
			}
		}
	}

	// If not in cache or cache disabled, get from repository
	var logs []*domain.Log
	err := s.retryOperation(func() error {
		var err error
		logs, err = s.repo.List(limit, offset)
		if err != nil {
			return err
		}

		// Cache the result if cache is enabled
		if s.cache != nil && logs != nil {
			s.cache.Set(s.keyGenerator.ForLogList(limit, offset), logs, s.config.Cache.TTL)
		}
		return nil
	})
	return logs, err
}

func (s *LogService) ListByTimeRange(start, end time.Time, limit, offset int) ([]*domain.Log, error) {
	// Validate time range
	if start.After(end) {
		return nil, fmt.Errorf("invalid time range: start time %v is after end time %v", start, end)
	}

	if s.cache != nil {
		// Try to get from cache first
		if cached, found := s.cache.Get(s.keyGenerator.ForTimeRange(start, end, limit, offset)); found {
			if logs, ok := cached.([]*domain.Log); ok {
				return logs, nil
			}
		}
	}

	// If not in cache or cache disabled, get from repository
	var logs []*domain.Log
	err := s.retryOperation(func() error {
		var err error
		logs, err = s.repo.ListByTimeRange(start, end, limit, offset)
		if err != nil {
			return err
		}

		// Cache the result if cache is enabled
		if s.cache != nil && logs != nil {
			s.cache.Set(s.keyGenerator.ForTimeRange(start, end, limit, offset), logs, s.config.Cache.TimeRangeTTL)
		}
		return nil
	})
	return logs, err
}
