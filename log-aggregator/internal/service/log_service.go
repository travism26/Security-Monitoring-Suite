package service

import (
	"encoding/json"
	"fmt"

	"github.com/travism26/log-aggregator/internal/domain"
)

type LogService struct {
	repo domain.LogRepository
}

func NewLogService(repo domain.LogRepository) *LogService {
	return &LogService{
		repo: repo,
	}
}

func (s *LogService) StoreLog(log *domain.Log) error {
	if len(log.Metadata) > 0 && log.MetadataStr == "" {
		metadataJSON, err := json.Marshal(log.Metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		log.MetadataStr = string(metadataJSON)
	}

	return s.repo.Store(log)
}

func (s *LogService) GetLog(id string) (*domain.Log, error) {
	return s.repo.FindByID(id)
}

func (s *LogService) ListLogs(limit, offset int) ([]*domain.Log, error) {
	return s.repo.List(limit, offset)
}
