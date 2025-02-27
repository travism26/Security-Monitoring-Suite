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
	fmt.Printf("[DEBUG] Storing log with ID: %s\n", log.ID)
	fmt.Printf("[DEBUG] Raw log data: %+v\n", log)

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
		fmt.Printf("[DEBUG] Attempting to store log with enriched data: %+v\n", log)
		if err := s.repo.Store(log); err != nil {
			fmt.Printf("[ERROR] Failed to store log: %v\n", err)
			return err
		}
		return nil
	})

	// If store was successful and cache is enabled, invalidate related cache entries
	if err == nil && s.cache != nil {
		s.cache.Delete(s.keyGenerator.ForLog(log.UserID, log.ID))
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
			s.cache.Delete(s.keyGenerator.ForLog(log.UserID, log.ID))
		}
		// Clear list caches as they might be affected
		s.cache.Clear() // TODO: Implement more granular cache invalidation
	}

	return err
}

func (s *LogService) GetLog(userID, id string) (*domain.Log, error) {
	if s.cache != nil {
		// Try to get from cache first
		cacheKey := s.keyGenerator.ForLog(userID, id)
		if cached, found := s.cache.Get(cacheKey); found {
			if log, ok := cached.(*domain.Log); ok {
				return log, nil
			}
		}
	}

	// If not in cache or cache disabled, get from repository
	var log *domain.Log
	err := s.retryOperation(func() error {
		var err error
		log, err = s.repo.FindByID(userID, id)
		if err != nil {
			return err
		}

		// Cache the result if cache is enabled
		if s.cache != nil && log != nil {
			s.cache.Set(s.keyGenerator.ForLog(userID, id), log, s.config.Cache.TTL)
		}
		return nil
	})
	return log, err
}

func (s *LogService) ListLogs(userID string, limit, offset int) ([]*domain.Log, error) {
	if s.cache != nil {
		// Try to get from cache first
		cacheKey := s.keyGenerator.ForLogList(userID, limit, offset)
		if cached, found := s.cache.Get(cacheKey); found {
			if logs, ok := cached.([]*domain.Log); ok {
				return logs, nil
			}
		}
	}

	// If not in cache or cache disabled, get from repository
	var logs []*domain.Log
	err := s.retryOperation(func() error {
		var err error
		logs, err = s.repo.List(userID, limit, offset)
		if err != nil {
			return err
		}

		// Cache the result if cache is enabled
		if s.cache != nil && logs != nil {
			s.cache.Set(s.keyGenerator.ForLogList(userID, limit, offset), logs, s.config.Cache.TTL)
		}
		return nil
	})
	return logs, err
}

func (s *LogService) ListByTimeRange(userID string, start, end time.Time, limit, offset int) ([]*domain.Log, error) {
	// Validate time range
	if start.After(end) {
		return nil, fmt.Errorf("invalid time range: start time %v is after end time %v", start, end)
	}

	if s.cache != nil {
		// Try to get from cache first
		cacheKey := s.keyGenerator.ForTimeRange(userID, start, end, limit, offset)
		if cached, found := s.cache.Get(cacheKey); found {
			if logs, ok := cached.([]*domain.Log); ok {
				return logs, nil
			}
		}
	}

	// If not in cache or cache disabled, get from repository
	var logs []*domain.Log
	err := s.retryOperation(func() error {
		var err error
		logs, err = s.repo.ListByTimeRange(userID, start, end, limit, offset)
		if err != nil {
			return err
		}

		// Cache the result if cache is enabled
		if s.cache != nil && logs != nil {
			s.cache.Set(s.keyGenerator.ForTimeRange(userID, start, end, limit, offset), logs, s.config.Cache.TimeRangeTTL)
		}
		return nil
	})
	return logs, err
}

// CountByTimeRange returns the total number of logs within a time range
func (s *LogService) CountByTimeRange(userID string, start, end time.Time) (int64, error) {
	if start.After(end) {
		return 0, fmt.Errorf("invalid time range: start time %v is after end time %v", start, end)
	}

	if s.cache != nil {
		cacheKey := s.keyGenerator.ForTimeRangeCount(userID, start, end)
		if cached, found := s.cache.Get(cacheKey); found {
			if count, ok := cached.(int64); ok {
				return count, nil
			}
		}
	}

	var count int64
	err := s.retryOperation(func() error {
		var err error
		count, err = s.repo.CountByTimeRange(userID, start, end)
		if err != nil {
			return err
		}

		if s.cache != nil {
			s.cache.Set(s.keyGenerator.ForTimeRangeCount(userID, start, end), count, s.config.Cache.TTL)
		}
		return nil
	})
	return count, err
}

// ListByHost retrieves logs for a specific host
func (s *LogService) ListByHost(userID, host string, limit, offset int) ([]*domain.Log, error) {
	if s.cache != nil {
		cacheKey := s.keyGenerator.ForHostLogs(userID, host, limit, offset)
		if cached, found := s.cache.Get(cacheKey); found {
			if logs, ok := cached.([]*domain.Log); ok {
				return logs, nil
			}
		}
	}

	var logs []*domain.Log
	err := s.retryOperation(func() error {
		var err error
		logs, err = s.repo.ListByHost(userID, host, limit, offset)
		if err != nil {
			return err
		}

		if s.cache != nil && logs != nil {
			s.cache.Set(s.keyGenerator.ForHostLogs(userID, host, limit, offset), logs, s.config.Cache.TTL)
		}
		return nil
	})
	return logs, err
}

// ListByLevel retrieves logs of a specific level
func (s *LogService) ListByLevel(userID, level string, limit, offset int) ([]*domain.Log, error) {
	if s.cache != nil {
		cacheKey := s.keyGenerator.ForLevelLogs(userID, level, limit, offset)
		if cached, found := s.cache.Get(cacheKey); found {
			if logs, ok := cached.([]*domain.Log); ok {
				return logs, nil
			}
		}
	}

	var logs []*domain.Log
	err := s.retryOperation(func() error {
		var err error
		logs, err = s.repo.ListByLevel(userID, level, limit, offset)
		if err != nil {
			return err
		}

		if s.cache != nil && logs != nil {
			s.cache.Set(s.keyGenerator.ForLevelLogs(userID, level, limit, offset), logs, s.config.Cache.TTL)
		}
		return nil
	})
	return logs, err
}
