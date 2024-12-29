package service

import (
	"github.com/travism26/log-aggregator/internal/domain"
)

type ProcessService struct {
	repo domain.ProcessRepository
}

func NewProcessService(repo domain.ProcessRepository) *ProcessService {
	return &ProcessService{
		repo: repo,
	}
}

func (s *ProcessService) StoreProcess(process *domain.Process) error {
	return s.repo.Store(process)
}

func (s *ProcessService) GetProcess(id string) (*domain.Process, error) {
	return s.repo.FindByID(id)
}

func (s *ProcessService) ListProcesses(limit, offset int) ([]*domain.Process, error) {
	return s.repo.List(limit, offset)
}
