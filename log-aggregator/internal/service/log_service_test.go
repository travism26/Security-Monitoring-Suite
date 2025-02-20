package service

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/travism26/log-aggregator/internal/domain"
)

// MockLogRepository implements domain.LogRepository for testing
type MockLogRepository struct {
	mock.Mock
}

func (m *MockLogRepository) Store(log *domain.Log) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockLogRepository) StoreBatch(logs []*domain.Log) error {
	args := m.Called(logs)
	return args.Error(0)
}

func (m *MockLogRepository) FindByID(userID, id string) (*domain.Log, error) {
	args := m.Called(userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Log), args.Error(1)
}

func (m *MockLogRepository) List(userID string, limit, offset int) ([]*domain.Log, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Log), args.Error(1)
}

func (m *MockLogRepository) ListByTimeRange(userID string, start, end time.Time, limit, offset int) ([]*domain.Log, error) {
	args := m.Called(userID, start, end, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Log), args.Error(1)
}

func (m *MockLogRepository) CountByTimeRange(userID string, start, end time.Time) (int64, error) {
	args := m.Called(userID, start, end)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLogRepository) ListByHost(userID string, host string, limit, offset int) ([]*domain.Log, error) {
	args := m.Called(userID, host, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Log), args.Error(1)
}

func (m *MockLogRepository) ListByLevel(userID string, level string, limit, offset int) ([]*domain.Log, error) {
	args := m.Called(userID, level, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Log), args.Error(1)
}

func (m *MockLogRepository) ListByUserID(userID string, limit, offset int) ([]*domain.Log, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Log), args.Error(1)
}

func (m *MockLogRepository) CountByUserID(userID string) (int64, error) {
	args := m.Called(userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockLogRepository) ListByUserIDAndTimeRange(userID string, start, end time.Time, limit, offset int) ([]*domain.Log, error) {
	args := m.Called(userID, start, end, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Log), args.Error(1)
}

func (m *MockLogRepository) ListByAPIKey(apiKey string, limit, offset int) ([]*domain.Log, error) {
	args := m.Called(apiKey, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Log), args.Error(1)
}

func (m *MockLogRepository) CountByAPIKey(apiKey string) (int64, error) {
	args := m.Called(apiKey)
	return args.Get(0).(int64), args.Error(1)
}

// getTestConfig returns a standard test configuration
func getTestConfig() LogServiceConfig {
	return LogServiceConfig{
		Environment: "test",
		Application: "log-aggregator",
		Component:   "test-component",
	}
}

func TestNewLogService(t *testing.T) {
	mockRepo := new(MockLogRepository)
	config := getTestConfig()
	service := NewLogService(mockRepo, config)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
	assert.Equal(t, config, service.config)
}

func TestLogService_StoreLog(t *testing.T) {
	tests := []struct {
		name          string
		log           *domain.Log
		setupMock     func(*MockLogRepository)
		expectError   bool
		expectedError string
	}{
		{
			name: "Successfully store log without metadata",
			log: &domain.Log{
				ID:   "test-id",
				Host: "test-host",
			},
			setupMock: func(m *MockLogRepository) {
				m.On("Store", mock.MatchedBy(func(log *domain.Log) bool {
					return log.Environment == "test" && log.Application == "log-aggregator"
				})).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Successfully store log with metadata",
			log: &domain.Log{
				ID:   "test-id",
				Host: "test-host",
				Metadata: map[string]interface{}{
					"key": "value",
				},
			},
			setupMock: func(m *MockLogRepository) {
				m.On("Store", mock.MatchedBy(func(log *domain.Log) bool {
					var metadata map[string]interface{}
					err := json.Unmarshal([]byte(log.MetadataStr), &metadata)
					return err == nil && metadata["key"] == "value" &&
						log.Environment == "test" && log.Application == "log-aggregator"
				})).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Repository error",
			log: &domain.Log{
				ID:   "test-id",
				Host: "test-host",
			},
			setupMock: func(m *MockLogRepository) {
				m.On("Store", mock.AnythingOfType("*domain.Log")).Return(errors.New("db error"))
			},
			expectError:   true,
			expectedError: "operation failed after 3 attempts: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockLogRepository)
			tt.setupMock(mockRepo)

			service := NewLogService(mockRepo, getTestConfig())
			err := service.StoreLog(tt.log)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogService_StoreBatch(t *testing.T) {
	tests := []struct {
		name          string
		logs          []*domain.Log
		setupMock     func(*MockLogRepository)
		expectError   bool
		expectedError string
	}{
		{
			name: "Successfully store batch of logs",
			logs: []*domain.Log{
				{
					ID:   "test-id-1",
					Host: "test-host-1",
					Metadata: map[string]interface{}{
						"key1": "value1",
					},
				},
				{
					ID:   "test-id-2",
					Host: "test-host-2",
					Metadata: map[string]interface{}{
						"key2": "value2",
					},
				},
			},
			setupMock: func(m *MockLogRepository) {
				m.On("StoreBatch", mock.MatchedBy(func(logs []*domain.Log) bool {
					return len(logs) == 2 &&
						logs[0].ID == "test-id-1" &&
						logs[1].ID == "test-id-2" &&
						logs[0].Environment == "test" &&
						logs[1].Environment == "test"
				})).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Repository error",
			logs: []*domain.Log{
				{
					ID:   "test-id-1",
					Host: "test-host-1",
				},
			},
			setupMock: func(m *MockLogRepository) {
				m.On("StoreBatch", mock.Anything).Return(errors.New("db error"))
			},
			expectError:   true,
			expectedError: "operation failed after 3 attempts: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockLogRepository)
			tt.setupMock(mockRepo)

			service := NewLogService(mockRepo, getTestConfig())
			err := service.StoreBatch(tt.logs)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogService_GetLog(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		id            string
		setupMock     func(*MockLogRepository)
		expectError   bool
		expectedError string
		validateLog   func(*testing.T, *domain.Log)
	}{
		{
			name:   "Successfully get log",
			userID: "test-user",
			id:     "test-id",
			setupMock: func(m *MockLogRepository) {
				m.On("FindByID", "test-user", "test-id").Return(&domain.Log{
					ID:     "test-id",
					UserID: "test-user",
					Host:   "test-host",
				}, nil)
			},
			validateLog: func(t *testing.T, log *domain.Log) {
				assert.Equal(t, "test-id", log.ID)
				assert.Equal(t, "test-user", log.UserID)
				assert.Equal(t, "test-host", log.Host)
			},
		},
		{
			name:   "Log not found",
			userID: "test-user",
			id:     "non-existent-id",
			setupMock: func(m *MockLogRepository) {
				m.On("FindByID", "test-user", "non-existent-id").Return(nil, errors.New("log not found"))
			},
			expectError:   true,
			expectedError: "operation failed after 3 attempts: log not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockLogRepository)
			tt.setupMock(mockRepo)

			service := NewLogService(mockRepo, getTestConfig())
			log, err := service.GetLog(tt.userID, tt.id)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, log)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, log)
				if tt.validateLog != nil {
					tt.validateLog(t, log)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogService_ListLogs(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		limit         int
		offset        int
		setupMock     func(*MockLogRepository)
		expectError   bool
		expectedError string
		validateLogs  func(*testing.T, []*domain.Log)
	}{
		{
			name:   "Successfully list logs",
			userID: "test-user",
			limit:  10,
			offset: 0,
			setupMock: func(m *MockLogRepository) {
				m.On("List", "test-user", 10, 0).Return([]*domain.Log{
					{ID: "1", UserID: "test-user", Host: "host-1"},
					{ID: "2", UserID: "test-user", Host: "host-2"},
				}, nil)
			},
			validateLogs: func(t *testing.T, logs []*domain.Log) {
				assert.Len(t, logs, 2)
				assert.Equal(t, "1", logs[0].ID)
				assert.Equal(t, "2", logs[1].ID)
				assert.Equal(t, "test-user", logs[0].UserID)
				assert.Equal(t, "test-user", logs[1].UserID)
			},
		},
		{
			name:   "Repository error",
			userID: "test-user",
			limit:  10,
			offset: 0,
			setupMock: func(m *MockLogRepository) {
				m.On("List", "test-user", 10, 0).Return(nil, errors.New("db error"))
			},
			expectError:   true,
			expectedError: "operation failed after 3 attempts: db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockLogRepository)
			tt.setupMock(mockRepo)

			service := NewLogService(mockRepo, getTestConfig())
			logs, err := service.ListLogs(tt.userID, tt.limit, tt.offset)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, logs)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
				if tt.validateLogs != nil {
					tt.validateLogs(t, logs)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestLogService_ListByTimeRange(t *testing.T) {
	now := time.Now()
	start := now.Add(-1 * time.Hour)
	end := now

	tests := []struct {
		name          string
		userID        string
		start         time.Time
		end           time.Time
		limit         int
		offset        int
		setupMock     func(*MockLogRepository)
		expectError   bool
		expectedError string
		validateLogs  func(*testing.T, []*domain.Log)
	}{
		{
			name:   "Successfully list logs by time range",
			userID: "test-user",
			start:  start,
			end:    end,
			limit:  10,
			offset: 0,
			setupMock: func(m *MockLogRepository) {
				m.On("ListByTimeRange", "test-user", start, end, 10, 0).Return([]*domain.Log{
					{ID: "1", UserID: "test-user", Host: "host-1", Timestamp: start.Add(15 * time.Minute)},
					{ID: "2", UserID: "test-user", Host: "host-2", Timestamp: start.Add(30 * time.Minute)},
				}, nil)
			},
			validateLogs: func(t *testing.T, logs []*domain.Log) {
				assert.Len(t, logs, 2)
				assert.Equal(t, "1", logs[0].ID)
				assert.Equal(t, "2", logs[1].ID)
				assert.Equal(t, "test-user", logs[0].UserID)
				assert.Equal(t, "test-user", logs[1].UserID)
				assert.True(t, logs[0].Timestamp.After(start) && logs[0].Timestamp.Before(end))
				assert.True(t, logs[1].Timestamp.After(start) && logs[1].Timestamp.Before(end))
			},
		},
		{
			name:   "Repository error",
			userID: "test-user",
			start:  start,
			end:    end,
			limit:  10,
			offset: 0,
			setupMock: func(m *MockLogRepository) {
				m.On("ListByTimeRange", "test-user", start, end, 10, 0).Return(nil, errors.New("db error"))
			},
			expectError:   true,
			expectedError: "operation failed after 3 attempts: db error",
		},
		{
			name:   "Invalid time range",
			userID: "test-user",
			start:  end,   // start time is after end time
			end:    start, // end time is before start time
			limit:  10,
			offset: 0,
			setupMock: func(m *MockLogRepository) {
				// Mock should not be called
			},
			expectError:   true,
			expectedError: "invalid time range",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockLogRepository)
			tt.setupMock(mockRepo)

			service := NewLogService(mockRepo, getTestConfig())
			logs, err := service.ListByTimeRange(tt.userID, tt.start, tt.end, tt.limit, tt.offset)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, logs)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
				if tt.validateLogs != nil {
					tt.validateLogs(t, logs)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
