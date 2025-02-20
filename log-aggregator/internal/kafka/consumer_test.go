package kafka

import (
	"encoding/json"
	"testing"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/travism26/log-aggregator/internal/domain"
)

// Interfaces for testing
type LogServiceInterface interface {
	StoreLog(log *domain.Log) error
	GetLog(userID, id string) (*domain.Log, error)
	ListLogs(userID string, limit, offset int) ([]*domain.Log, error)
}

type AlertServiceInterface interface {
	ProcessMetrics(log *domain.Log) error
	GetAlert(id string) (*domain.Alert, error)
	ListAlerts(limit, offset int) ([]*domain.Alert, error)
	ListAlertsByStatus(status domain.AlertStatus, limit, offset int) ([]*domain.Alert, error)
	ListAlertsBySeverity(severity domain.AlertSeverity, limit, offset int) ([]*domain.Alert, error)
	UpdateAlertStatus(id string, status domain.AlertStatus) error
}

type ProcessRepositoryInterface interface {
	StoreBatch(processes []domain.Process) error
	FindByLogID(logID string) ([]domain.Process, error)
}

// Mock implementations
type MockLogService struct {
	mock.Mock
}

func (m *MockLogService) StoreLog(log *domain.Log) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockLogService) GetLog(userID, id string) (*domain.Log, error) {
	args := m.Called(userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Log), args.Error(1)
}

func (m *MockLogService) ListLogs(userID string, limit, offset int) ([]*domain.Log, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Log), args.Error(1)
}

type MockAlertService struct {
	mock.Mock
}

func (m *MockAlertService) ProcessMetrics(log *domain.Log) error {
	args := m.Called(log)
	return args.Error(0)
}

func (m *MockAlertService) GetAlert(id string) (*domain.Alert, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Alert), args.Error(1)
}

func (m *MockAlertService) ListAlerts(limit, offset int) ([]*domain.Alert, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Alert), args.Error(1)
}

func (m *MockAlertService) ListAlertsByStatus(status domain.AlertStatus, limit, offset int) ([]*domain.Alert, error) {
	args := m.Called(status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Alert), args.Error(1)
}

func (m *MockAlertService) ListAlertsBySeverity(severity domain.AlertSeverity, limit, offset int) ([]*domain.Alert, error) {
	args := m.Called(severity, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Alert), args.Error(1)
}

func (m *MockAlertService) UpdateAlertStatus(id string, status domain.AlertStatus) error {
	args := m.Called(id, status)
	return args.Error(0)
}

type MockProcessRepository struct {
	mock.Mock
}

func (m *MockProcessRepository) StoreBatch(processes []domain.Process) error {
	args := m.Called(processes)
	return args.Error(0)
}

func (m *MockProcessRepository) FindByLogID(logID string) ([]domain.Process, error) {
	args := m.Called(logID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Process), args.Error(1)
}

func TestConsumer_ProcessMessage(t *testing.T) {
	tests := []struct {
		name          string
		messageJSON   string
		expectError   bool
		expectedError string
	}{
		{
			name:        "Valid message with processes",
			messageJSON: `{"tenant_id":"67a5da7f9f3f88e40759e219","api_key":"sms_123456","host":{"os":"darwin","arch":"arm64","hostname":"Traviss-MacBook-Pro.local","cpu_cores":12,"go_version":"go1.23.2"},"metrics":{"cpu_usage":13.13563381573259,"disk":{"free":269959933952,"total":494384795648,"usage_percent":45.39477420656553,"used":224424861696},"memory_usage":13316669440,"memory_usage_percent":68.90063815646701,"network":{"BytesReceived":47879701332,"BytesSent":7329810320},"processes":{"process_list":[{"cpu_percent":0.22119553205535347,"memory_usage":13090816,"name":"launchd","pid":1,"status":"S"}],"total_count":1,"total_cpu_percent":0.22119553205535347,"total_memory_usage":13090816}},"threat_indicators":[{"type":"high_cpu_usage","description":"CPU usage exceeds threshold","severity":"low","score":23.635782970117063,"timestamp":"2024-12-28T08:00:11.665024-05:00","metadata":{"tags":["performance","resource_usage"]}}],"metadata":{"collection_duration":"8.466914875s","collector_count":5}}`,
			expectError: false,
		},
		{
			name: "Valid message without processes",
			messageJSON: `{
				"tenant_id": "67a5da7f9f3f88e40759e219",
				"api_key": "sms_123456",
				"host": {
					"hostname": "test-host",
					"os": "linux"
				},
				"metrics": {
					"cpu_usage": 50.5,
					"memory_usage_percent": 75.0
				},
				"processes": null
			}`,
			expectError: false,
		},
		{
			name: "Invalid message - missing required fields",
			messageJSON: `{
				"tenant_id": "67a5da7f9f3f88e40759e219",
				"api_key": "sms_123456",
				"host": {
					"os": "linux"
				},
				"metrics": {
					"cpu_usage": 50.5
				}
			}`,
			expectError:   true,
			expectedError: "invalid hostname format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mocks
			mockLogService := new(MockLogService)
			mockAlertService := new(MockAlertService)
			mockProcessRepo := new(MockProcessRepository)

			// Setup expectations
			mockLogService.On("StoreLog", mock.Anything).Return(nil)
			mockAlertService.On("ProcessMetrics", mock.Anything).Return(nil)
			mockProcessRepo.On("StoreBatch", mock.Anything).Return(nil)

			// Create consumer with interface implementations
			consumer := &Consumer{
				logService:        LogServiceInterface(mockLogService),
				alertService:      AlertServiceInterface(mockAlertService),
				processRepository: ProcessRepositoryInterface(mockProcessRepo),
			}

			// Create message
			msg := &sarama.ConsumerMessage{
				Value: []byte(tt.messageJSON),
			}

			// Process message
			err := consumer.processMessage(msg)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
				mockLogService.AssertExpectations(t)
				mockAlertService.AssertExpectations(t)
				mockProcessRepo.AssertExpectations(t)
			}
		})
	}
}

func TestConsumer_ExtractProcesses(t *testing.T) {
	tests := []struct {
		name           string
		processes      interface{}
		expectedCount  int
		expectError    bool
		expectedError  string
		validateFields bool
	}{
		{
			name:          "Nil processes",
			processes:     nil,
			expectedCount: 0,
			expectError:   false,
		},
		{
			name: "Valid processes",
			processes: map[string]interface{}{
				"process_list": []interface{}{
					map[string]interface{}{
						"name":         "test-process",
						"pid":          float64(123),
						"cpu_percent":  float64(10.5),
						"memory_usage": float64(512),
						"status":       "running",
					},
				},
			},
			expectedCount:  1,
			expectError:    false,
			validateFields: true,
		},
		{
			name:          "Invalid processes format",
			processes:     "invalid",
			expectError:   true,
			expectedError: "invalid processes data format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumer := &Consumer{}
			logID := uuid.New().String()

			rawMsg := &struct {
				Host             interface{} `json:"host"`
				Metrics          interface{} `json:"metrics"`
				ThreatIndicators interface{} `json:"threat_indicators"`
				Metadata         interface{} `json:"metadata"`
				Processes        interface{} `json:"processes"`
				TenantID         string      `json:"tenant_id"`
				APIKey           string      `json:"api_key"`
			}{
				Processes: tt.processes,
				TenantID:  "67a5da7f9f3f88e40759e219",
				APIKey:    "sms_123456",
			}

			processes, err := consumer.extractProcesses(rawMsg, logID)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, processes, tt.expectedCount)

				if tt.validateFields && len(processes) > 0 {
					process := processes[0]
					assert.Equal(t, "test-process", process.Name)
					assert.Equal(t, 123, process.PID)
					assert.Equal(t, 10.5, process.CPUPercent)
					assert.Equal(t, int64(512), process.MemoryUsage)
					assert.Equal(t, "running", process.Status)
					assert.Equal(t, logID, process.LogID)
					assert.NotEmpty(t, process.ID)
					assert.False(t, process.Timestamp.IsZero())
				}
			}
		})
	}
}

func TestConsumer_CreateLogEntry(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectError   bool
		expectedError string
		validate      func(*testing.T, *domain.Log)
	}{
		{
			name: "Valid input with processes",
			input: `{
				"tenant_id": "67a5da7f9f3f88e40759e219",
				"api_key": "sms_123456",
				"host": {"hostname": "test-host"},
				"metrics": {
					"cpu_usage": 50.5,
					"memory_usage_percent": 75.0
				},
				"processes": {
					"total_count": 10,
					"total_cpu_percent": 80.5,
					"total_memory_usage": 1024
				}
			}`,
			expectError: false,
			validate: func(t *testing.T, log *domain.Log) {
				assert.Equal(t, "test-host", log.Host)
				assert.Equal(t, "67a5da7f9f3f88e40759e219", log.OrganizationID)
				assert.Equal(t, 10, log.ProcessCount)
				assert.Equal(t, 80.5, log.TotalCPUPercent)
				assert.Equal(t, int64(1024), log.TotalMemoryUsage)
			},
		},
		{
			name: "Valid input without processes",
			input: `{
				"tenant_id": "67a5da7f9f3f88e40759e219",
				"api_key": "sms_123456",
				"host": {"hostname": "test-host"},
				"metrics": {
					"cpu_usage": 50.5,
					"memory_usage_percent": 75.0
				},
				"processes": null
			}`,
			expectError: false,
			validate: func(t *testing.T, log *domain.Log) {
				assert.Equal(t, "test-host", log.Host)
				assert.Equal(t, "67a5da7f9f3f88e40759e219", log.OrganizationID)
				assert.Equal(t, 0, log.ProcessCount)
				assert.Equal(t, 0.0, log.TotalCPUPercent)
				assert.Equal(t, int64(0), log.TotalMemoryUsage)
			},
		},
		{
			name: "Missing hostname",
			input: `{
				"tenant_id": "67a5da7f9f3f88e40759e219",
				"api_key": "sms_123456",
				"host": {},
				"metrics": {
					"cpu_usage": 50.5,
					"memory_usage_percent": 75.0
				}
			}`,
			expectError:   true,
			expectedError: "invalid hostname format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumer := &Consumer{}
			var rawMsg struct {
				Host             interface{} `json:"host"`
				Metrics          interface{} `json:"metrics"`
				ThreatIndicators interface{} `json:"threat_indicators"`
				Metadata         interface{} `json:"metadata"`
				Processes        interface{} `json:"processes"`
				TenantID         string      `json:"tenant_id"`
				APIKey           string      `json:"api_key"`
			}

			err := json.Unmarshal([]byte(tt.input), &rawMsg)
			assert.NoError(t, err, "Failed to unmarshal test input")

			logEntry, err := consumer.createLogEntry(&rawMsg)

			if tt.expectError {
				assert.Error(t, err)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, logEntry)
				assert.NotEmpty(t, logEntry.ID)
				assert.False(t, logEntry.Timestamp.IsZero())
				if tt.validate != nil {
					tt.validate(t, logEntry)
				}
			}
		})
	}
}
