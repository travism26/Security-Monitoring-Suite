package service

import (
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
	return s.repo.Store(log)
}

func (s *LogService) GetLog(id string) (*domain.Log, error) {
	return s.repo.FindByID(id)
}

func (s *LogService) ListLogs(limit, offset int) ([]*domain.Log, error) {
	return s.repo.List(limit, offset)
}
