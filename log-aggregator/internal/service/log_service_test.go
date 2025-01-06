package service

import (
	"encoding/json"
	"errors"
	"testing"

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

func (m *MockLogRepository) FindByID(id string) (*domain.Log, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Log), args.Error(1)
}

func (m *MockLogRepository) List(limit, offset int) ([]*domain.Log, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Log), args.Error(1)
}

func TestNewLogService(t *testing.T) {
	mockRepo := new(MockLogRepository)
	service := NewLogService(mockRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
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
				m.On("Store", mock.AnythingOfType("*domain.Log")).Return(nil)
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
					return err == nil && metadata["key"] == "value"
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
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockLogRepository)
			tt.setupMock(mockRepo)

			service := NewLogService(mockRepo)
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

func TestLogService_GetLog(t *testing.T) {
	tests := []struct {
		name          string
		id            string
		setupMock     func(*MockLogRepository)
		expectError   bool
		expectedError string
		validateLog   func(*testing.T, *domain.Log)
	}{
		{
			name: "Successfully get log",
			id:   "test-id",
			setupMock: func(m *MockLogRepository) {
				m.On("FindByID", "test-id").Return(&domain.Log{
					ID:   "test-id",
					Host: "test-host",
				}, nil)
			},
			expectError: false,
			validateLog: func(t *testing.T, log *domain.Log) {
				assert.Equal(t, "test-id", log.ID)
				assert.Equal(t, "test-host", log.Host)
			},
		},
		{
			name: "Log not found",
			id:   "non-existent-id",
			setupMock: func(m *MockLogRepository) {
				m.On("FindByID", "non-existent-id").Return(nil, errors.New("log not found"))
			},
			expectError:   true,
			expectedError: "log not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockLogRepository)
			tt.setupMock(mockRepo)

			service := NewLogService(mockRepo)
			log, err := service.GetLog(tt.id)

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
		limit         int
		offset        int
		setupMock     func(*MockLogRepository)
		expectError   bool
		expectedError string
		validateLogs  func(*testing.T, []*domain.Log)
	}{
		{
			name:   "Successfully list logs",
			limit:  10,
			offset: 0,
			setupMock: func(m *MockLogRepository) {
				m.On("List", 10, 0).Return([]*domain.Log{
					{ID: "1", Host: "host-1"},
					{ID: "2", Host: "host-2"},
				}, nil)
			},
			expectError: false,
			validateLogs: func(t *testing.T, logs []*domain.Log) {
				assert.Len(t, logs, 2)
				assert.Equal(t, "1", logs[0].ID)
				assert.Equal(t, "2", logs[1].ID)
			},
		},
		{
			name:   "Empty result",
			limit:  10,
			offset: 100,
			setupMock: func(m *MockLogRepository) {
				m.On("List", 10, 100).Return([]*domain.Log{}, nil)
			},
			expectError: false,
			validateLogs: func(t *testing.T, logs []*domain.Log) {
				assert.Empty(t, logs)
			},
		},
		{
			name:   "Repository error",
			limit:  10,
			offset: 0,
			setupMock: func(m *MockLogRepository) {
				m.On("List", 10, 0).Return(nil, errors.New("db error"))
			},
			expectError:   true,
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockLogRepository)
			tt.setupMock(mockRepo)

			service := NewLogService(mockRepo)
			logs, err := service.ListLogs(tt.limit, tt.offset)

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
