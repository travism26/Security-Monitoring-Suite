package mock

import (
	"time"

	"github.com/travism26/log-aggregator/internal/domain"
)

// MockLogRepository is a mock implementation of the LogRepository interface
type MockLogRepository struct {
	logs []*domain.Log
}

func NewMockLogRepository() *MockLogRepository {
	return &MockLogRepository{
		logs: make([]*domain.Log, 0),
	}
}

func (m *MockLogRepository) Store(log *domain.Log) error {
	m.logs = append(m.logs, log)
	return nil
}

func (m *MockLogRepository) StoreBatch(logs []*domain.Log) error {
	m.logs = append(m.logs, logs...)
	return nil
}

func (m *MockLogRepository) FindByID(userID, id string) (*domain.Log, error) {
	for _, log := range m.logs {
		if log.ID == id && log.UserID == userID {
			return log, nil
		}
	}
	return nil, nil
}

func (m *MockLogRepository) List(userID string, limit, offset int) ([]*domain.Log, error) {
	var result []*domain.Log
	for _, log := range m.logs {
		if log.UserID == userID {
			result = append(result, log)
		}
	}
	return paginateResults(result, limit, offset), nil
}

func (m *MockLogRepository) ListByTimeRange(userID string, start, end time.Time, limit, offset int) ([]*domain.Log, error) {
	var result []*domain.Log
	for _, log := range m.logs {
		if log.UserID == userID &&
			log.Timestamp.After(start) && log.Timestamp.Before(end) {
			result = append(result, log)
		}
	}
	return paginateResults(result, limit, offset), nil
}

func (m *MockLogRepository) CountByTimeRange(userID string, start, end time.Time) (int64, error) {
	var count int64
	for _, log := range m.logs {
		if log.UserID == userID &&
			log.Timestamp.After(start) && log.Timestamp.Before(end) {
			count++
		}
	}
	return count, nil
}

func (m *MockLogRepository) ListByHost(userID string, host string, limit, offset int) ([]*domain.Log, error) {
	var result []*domain.Log
	for _, log := range m.logs {
		if log.UserID == userID && log.Host == host {
			result = append(result, log)
		}
	}
	return paginateResults(result, limit, offset), nil
}

func (m *MockLogRepository) ListByLevel(userID string, level string, limit, offset int) ([]*domain.Log, error) {
	var result []*domain.Log
	for _, log := range m.logs {
		if log.UserID == userID && log.Level == level {
			result = append(result, log)
		}
	}
	return paginateResults(result, limit, offset), nil
}

func (m *MockLogRepository) ListByAPIKey(apiKey string, limit, offset int) ([]*domain.Log, error) {
	var result []*domain.Log
	for _, log := range m.logs {
		if log.APIKey == apiKey {
			result = append(result, log)
		}
	}
	return paginateResults(result, limit, offset), nil
}

func (m *MockLogRepository) CountByAPIKey(apiKey string) (int64, error) {
	var count int64
	for _, log := range m.logs {
		if log.APIKey == apiKey {
			count++
		}
	}
	return count, nil
}

func (m *MockLogRepository) ListByUserID(userID string, limit, offset int) ([]*domain.Log, error) {
	var result []*domain.Log
	for _, log := range m.logs {
		if log.UserID == userID {
			result = append(result, log)
		}
	}
	return paginateResults(result, limit, offset), nil
}

func (m *MockLogRepository) CountByUserID(userID string) (int64, error) {
	var count int64
	for _, log := range m.logs {
		if log.UserID == userID {
			count++
		}
	}
	return count, nil
}

func (m *MockLogRepository) ListByUserIDAndTimeRange(userID string, start, end time.Time, limit, offset int) ([]*domain.Log, error) {
	var result []*domain.Log
	for _, log := range m.logs {
		if log.UserID == userID &&
			log.Timestamp.After(start) && log.Timestamp.Before(end) {
			result = append(result, log)
		}
	}
	return paginateResults(result, limit, offset), nil
}

// Helper function to paginate results
func paginateResults(logs []*domain.Log, limit, offset int) []*domain.Log {
	if offset >= len(logs) {
		return []*domain.Log{}
	}

	end := offset + limit
	if end > len(logs) {
		end = len(logs)
	}

	return logs[offset:end]
}
