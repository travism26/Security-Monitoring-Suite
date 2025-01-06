package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/travism26/log-aggregator/internal/domain"
)

// MockProcessRepository implements domain.ProcessRepository for testing
type MockProcessRepository struct {
	mock.Mock
}

func (m *MockProcessRepository) Store(process *domain.Process) error {
	args := m.Called(process)
	return args.Error(0)
}

func (m *MockProcessRepository) FindByID(id string) (*domain.Process, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Process), args.Error(1)
}

func (m *MockProcessRepository) List(limit, offset int) ([]*domain.Process, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Process), args.Error(1)
}

func TestNewProcessService(t *testing.T) {
	mockRepo := new(MockProcessRepository)
	service := NewProcessService(mockRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
}

func TestProcessService_StoreProcess(t *testing.T) {
	tests := []struct {
		name          string
		process       *domain.Process
		setupMock     func(*MockProcessRepository)
		expectError   bool
		expectedError string
	}{
		{
			name: "Successfully store process",
			process: &domain.Process{
				ID:          "test-id",
				Name:        "test-process",
				PID:         1234,
				CPUPercent:  5.0,
				MemoryUsage: 1024 * 1024,
				Status:      "running",
				LogID:       "log-1",
				Timestamp:   time.Now(),
			},
			setupMock: func(m *MockProcessRepository) {
				m.On("Store", mock.MatchedBy(func(p *domain.Process) bool {
					return p.ID == "test-id" && p.Name == "test-process"
				})).Return(nil)
			},
			expectError: false,
		},
		{
			name: "Repository error",
			process: &domain.Process{
				ID:   "test-id",
				Name: "test-process",
			},
			setupMock: func(m *MockProcessRepository) {
				m.On("Store", mock.Anything).Return(errors.New("db error"))
			},
			expectError:   true,
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProcessRepository)
			tt.setupMock(mockRepo)

			service := NewProcessService(mockRepo)
			err := service.StoreProcess(tt.process)

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

func TestProcessService_GetProcess(t *testing.T) {
	tests := []struct {
		name            string
		processID       string
		setupMock       func(*MockProcessRepository)
		expectError     bool
		expectedError   string
		validateProcess func(*testing.T, *domain.Process)
	}{
		{
			name:      "Successfully get process",
			processID: "test-id",
			setupMock: func(m *MockProcessRepository) {
				m.On("FindByID", "test-id").Return(&domain.Process{
					ID:          "test-id",
					Name:        "test-process",
					PID:         1234,
					CPUPercent:  5.0,
					MemoryUsage: 1024 * 1024,
					Status:      "running",
				}, nil)
			},
			validateProcess: func(t *testing.T, process *domain.Process) {
				assert.Equal(t, "test-id", process.ID)
				assert.Equal(t, "test-process", process.Name)
				assert.Equal(t, 1234, process.PID)
				assert.Equal(t, 5.0, process.CPUPercent)
				assert.Equal(t, int64(1024*1024), process.MemoryUsage)
				assert.Equal(t, "running", process.Status)
			},
		},
		{
			name:      "Process not found",
			processID: "non-existent",
			setupMock: func(m *MockProcessRepository) {
				m.On("FindByID", "non-existent").Return(nil, errors.New("not found"))
			},
			expectError:   true,
			expectedError: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProcessRepository)
			tt.setupMock(mockRepo)

			service := NewProcessService(mockRepo)
			process, err := service.GetProcess(tt.processID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, process)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, process)
				if tt.validateProcess != nil {
					tt.validateProcess(t, process)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestProcessService_ListProcesses(t *testing.T) {
	mockProcesses := []*domain.Process{
		{
			ID:          "1",
			Name:        "process-1",
			PID:         1234,
			CPUPercent:  5.0,
			MemoryUsage: 1024 * 1024,
			Status:      "running",
		},
		{
			ID:          "2",
			Name:        "process-2",
			PID:         5678,
			CPUPercent:  10.0,
			MemoryUsage: 2048 * 1024,
			Status:      "sleeping",
		},
	}

	tests := []struct {
		name              string
		limit             int
		offset            int
		setupMock         func(*MockProcessRepository)
		expectError       bool
		expectedError     string
		validateProcesses func(*testing.T, []*domain.Process)
	}{
		{
			name:   "Successfully list processes",
			limit:  10,
			offset: 0,
			setupMock: func(m *MockProcessRepository) {
				m.On("List", 10, 0).Return(mockProcesses, nil)
			},
			validateProcesses: func(t *testing.T, processes []*domain.Process) {
				assert.Len(t, processes, 2)
				assert.Equal(t, "process-1", processes[0].Name)
				assert.Equal(t, "process-2", processes[1].Name)
				assert.Equal(t, 1234, processes[0].PID)
				assert.Equal(t, 5678, processes[1].PID)
			},
		},
		{
			name:   "Empty result",
			limit:  10,
			offset: 100,
			setupMock: func(m *MockProcessRepository) {
				m.On("List", 10, 100).Return([]*domain.Process{}, nil)
			},
			validateProcesses: func(t *testing.T, processes []*domain.Process) {
				assert.Empty(t, processes)
			},
		},
		{
			name:   "Repository error",
			limit:  10,
			offset: 0,
			setupMock: func(m *MockProcessRepository) {
				m.On("List", 10, 0).Return(nil, errors.New("db error"))
			},
			expectError:   true,
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockProcessRepository)
			tt.setupMock(mockRepo)

			service := NewProcessService(mockRepo)
			processes, err := service.ListProcesses(tt.limit, tt.offset)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, processes)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
				if tt.validateProcesses != nil {
					tt.validateProcesses(t, processes)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
