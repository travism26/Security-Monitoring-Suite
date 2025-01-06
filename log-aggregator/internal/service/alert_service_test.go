package service

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/travism26/log-aggregator/internal/domain"
)

// MockAlertRepository implements domain.AlertRepository for testing
type MockAlertRepository struct {
	mock.Mock
}

func (m *MockAlertRepository) Store(alert *domain.Alert) error {
	args := m.Called(alert)
	return args.Error(0)
}

func (m *MockAlertRepository) Update(alert *domain.Alert) error {
	args := m.Called(alert)
	return args.Error(0)
}

func (m *MockAlertRepository) FindByID(id string) (*domain.Alert, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Alert), args.Error(1)
}

func (m *MockAlertRepository) List(limit, offset int) ([]*domain.Alert, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Alert), args.Error(1)
}

func (m *MockAlertRepository) FindByStatus(status domain.AlertStatus, limit, offset int) ([]*domain.Alert, error) {
	args := m.Called(status, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Alert), args.Error(1)
}

func (m *MockAlertRepository) FindBySeverity(severity domain.AlertSeverity, limit, offset int) ([]*domain.Alert, error) {
	args := m.Called(severity, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Alert), args.Error(1)
}

func TestNewAlertService(t *testing.T) {
	mockRepo := new(MockAlertRepository)
	service := NewAlertService(mockRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
	assert.Equal(t, DefaultThresholds, service.thresholds)
}

func TestAlertService_SetThresholds(t *testing.T) {
	service := NewAlertService(new(MockAlertRepository))
	customThresholds := AlertThresholds{
		CPUUsagePercent:    90.0,
		MemoryUsagePercent: 95.0,
		ProcessCount:       2000,
	}

	service.SetThresholds(customThresholds)
	assert.Equal(t, customThresholds, service.thresholds)
}

func TestAlertService_ProcessMetrics(t *testing.T) {
	tests := []struct {
		name          string
		log           *domain.Log
		thresholds    AlertThresholds
		setupMock     func(*MockAlertRepository)
		expectError   bool
		expectedError string
		alertCount    int
	}{
		{
			name: "No thresholds exceeded",
			log: &domain.Log{
				ID:               "test-id",
				Host:             "test-host",
				TotalCPUPercent:  50.0,
				TotalMemoryUsage: 8 * 1024 * 1024 * 1024, // 8GB
				ProcessCount:     500,
			},
			thresholds: DefaultThresholds,
			setupMock:  func(m *MockAlertRepository) {},
			alertCount: 0,
		},
		{
			name: "CPU threshold exceeded",
			log: &domain.Log{
				ID:               "test-id",
				Host:             "test-host",
				TotalCPUPercent:  90.0,
				TotalMemoryUsage: 8 * 1024 * 1024 * 1024,
				ProcessCount:     500,
			},
			thresholds: DefaultThresholds,
			setupMock: func(m *MockAlertRepository) {
				m.On("Store", mock.MatchedBy(func(alert *domain.Alert) bool {
					return alert.Severity == domain.SeverityHigh &&
						alert.Status == domain.StatusOpen
				})).Return(nil)
			},
			alertCount: 1,
		},
		{
			name: "All thresholds exceeded",
			log: &domain.Log{
				ID:               "test-id",
				Host:             "test-host",
				TotalCPUPercent:  90.0,
				TotalMemoryUsage: 15 * 1024 * 1024 * 1024,
				ProcessCount:     1500,
			},
			thresholds: DefaultThresholds,
			setupMock: func(m *MockAlertRepository) {
				m.On("Store", mock.MatchedBy(func(alert *domain.Alert) bool {
					return alert.Status == domain.StatusOpen
				})).Return(nil).Times(3)
			},
			alertCount: 3,
		},
		{
			name: "Repository error",
			log: &domain.Log{
				ID:              "test-id",
				Host:            "test-host",
				TotalCPUPercent: 90.0,
			},
			thresholds: DefaultThresholds,
			setupMock: func(m *MockAlertRepository) {
				m.On("Store", mock.Anything).Return(errors.New("db error"))
			},
			expectError:   true,
			expectedError: "failed to store alert",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAlertRepository)
			tt.setupMock(mockRepo)

			service := NewAlertService(mockRepo)
			service.SetThresholds(tt.thresholds)

			err := service.ProcessMetrics(tt.log)

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

func TestAlertService_UpdateAlertStatus(t *testing.T) {
	tests := []struct {
		name          string
		alertID       string
		newStatus     domain.AlertStatus
		setupMock     func(*MockAlertRepository)
		expectError   bool
		expectedError string
	}{
		{
			name:      "Successfully update to resolved",
			alertID:   "test-id",
			newStatus: domain.StatusResolved,
			setupMock: func(m *MockAlertRepository) {
				now := time.Now()
				m.On("FindByID", "test-id").Return(&domain.Alert{
					ID:        "test-id",
					Status:    domain.StatusOpen,
					CreatedAt: now,
					UpdatedAt: now,
				}, nil)
				m.On("Update", mock.MatchedBy(func(alert *domain.Alert) bool {
					return alert.Status == domain.StatusResolved &&
						alert.ResolvedAt != nil &&
						alert.UpdatedAt.After(now)
				})).Return(nil)
			},
		},
		{
			name:      "Alert not found",
			alertID:   "non-existent",
			newStatus: domain.StatusResolved,
			setupMock: func(m *MockAlertRepository) {
				m.On("FindByID", "non-existent").Return(nil, errors.New("not found"))
			},
			expectError:   true,
			expectedError: "failed to find alert",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAlertRepository)
			tt.setupMock(mockRepo)

			service := NewAlertService(mockRepo)
			err := service.UpdateAlertStatus(tt.alertID, tt.newStatus)

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

func TestAlertService_GetAlert(t *testing.T) {
	tests := []struct {
		name          string
		alertID       string
		setupMock     func(*MockAlertRepository)
		expectError   bool
		expectedError string
		validateAlert func(*testing.T, *domain.Alert)
	}{
		{
			name:    "Successfully get alert",
			alertID: "test-id",
			setupMock: func(m *MockAlertRepository) {
				m.On("FindByID", "test-id").Return(&domain.Alert{
					ID:       "test-id",
					Title:    "Test Alert",
					Severity: domain.SeverityHigh,
				}, nil)
			},
			validateAlert: func(t *testing.T, alert *domain.Alert) {
				assert.Equal(t, "test-id", alert.ID)
				assert.Equal(t, "Test Alert", alert.Title)
				assert.Equal(t, domain.SeverityHigh, alert.Severity)
			},
		},
		{
			name:    "Alert not found",
			alertID: "non-existent",
			setupMock: func(m *MockAlertRepository) {
				m.On("FindByID", "non-existent").Return(nil, errors.New("not found"))
			},
			expectError:   true,
			expectedError: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAlertRepository)
			tt.setupMock(mockRepo)

			service := NewAlertService(mockRepo)
			alert, err := service.GetAlert(tt.alertID)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, alert)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, alert)
				if tt.validateAlert != nil {
					tt.validateAlert(t, alert)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAlertService_ListAlerts(t *testing.T) {
	mockAlerts := []*domain.Alert{
		{ID: "1", Title: "Alert 1", Severity: domain.SeverityHigh},
		{ID: "2", Title: "Alert 2", Severity: domain.SeverityMedium},
	}

	tests := []struct {
		name           string
		limit          int
		offset         int
		setupMock      func(*MockAlertRepository)
		expectError    bool
		expectedError  string
		validateAlerts func(*testing.T, []*domain.Alert)
	}{
		{
			name:   "Successfully list alerts",
			limit:  10,
			offset: 0,
			setupMock: func(m *MockAlertRepository) {
				m.On("List", 10, 0).Return(mockAlerts, nil)
			},
			validateAlerts: func(t *testing.T, alerts []*domain.Alert) {
				assert.Len(t, alerts, 2)
				assert.Equal(t, "Alert 1", alerts[0].Title)
				assert.Equal(t, "Alert 2", alerts[1].Title)
			},
		},
		{
			name:   "Repository error",
			limit:  10,
			offset: 0,
			setupMock: func(m *MockAlertRepository) {
				m.On("List", 10, 0).Return(nil, errors.New("db error"))
			},
			expectError:   true,
			expectedError: "db error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAlertRepository)
			tt.setupMock(mockRepo)

			service := NewAlertService(mockRepo)
			alerts, err := service.ListAlerts(tt.limit, tt.offset)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, alerts)
				if tt.expectedError != "" {
					assert.Contains(t, err.Error(), tt.expectedError)
				}
			} else {
				assert.NoError(t, err)
				if tt.validateAlerts != nil {
					tt.validateAlerts(t, alerts)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAlertService_ListAlertsByStatus(t *testing.T) {
	mockAlerts := []*domain.Alert{
		{ID: "1", Title: "Alert 1", Status: domain.StatusOpen},
		{ID: "2", Title: "Alert 2", Status: domain.StatusOpen},
	}

	tests := []struct {
		name           string
		status         domain.AlertStatus
		limit          int
		offset         int
		setupMock      func(*MockAlertRepository)
		expectError    bool
		expectedError  string
		validateAlerts func(*testing.T, []*domain.Alert)
	}{
		{
			name:   "Successfully list alerts by status",
			status: domain.StatusOpen,
			limit:  10,
			offset: 0,
			setupMock: func(m *MockAlertRepository) {
				m.On("FindByStatus", domain.StatusOpen, 10, 0).Return(mockAlerts, nil)
			},
			validateAlerts: func(t *testing.T, alerts []*domain.Alert) {
				assert.Len(t, alerts, 2)
				for _, alert := range alerts {
					assert.Equal(t, domain.StatusOpen, alert.Status)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAlertRepository)
			tt.setupMock(mockRepo)

			service := NewAlertService(mockRepo)
			alerts, err := service.ListAlertsByStatus(tt.status, tt.limit, tt.offset)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, alerts)
			} else {
				assert.NoError(t, err)
				if tt.validateAlerts != nil {
					tt.validateAlerts(t, alerts)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAlertService_ListAlertsBySeverity(t *testing.T) {
	mockAlerts := []*domain.Alert{
		{ID: "1", Title: "Alert 1", Severity: domain.SeverityHigh},
		{ID: "2", Title: "Alert 2", Severity: domain.SeverityHigh},
	}

	tests := []struct {
		name           string
		severity       domain.AlertSeverity
		limit          int
		offset         int
		setupMock      func(*MockAlertRepository)
		expectError    bool
		expectedError  string
		validateAlerts func(*testing.T, []*domain.Alert)
	}{
		{
			name:     "Successfully list alerts by severity",
			severity: domain.SeverityHigh,
			limit:    10,
			offset:   0,
			setupMock: func(m *MockAlertRepository) {
				m.On("FindBySeverity", domain.SeverityHigh, 10, 0).Return(mockAlerts, nil)
			},
			validateAlerts: func(t *testing.T, alerts []*domain.Alert) {
				assert.Len(t, alerts, 2)
				for _, alert := range alerts {
					assert.Equal(t, domain.SeverityHigh, alert.Severity)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockAlertRepository)
			tt.setupMock(mockRepo)

			service := NewAlertService(mockRepo)
			alerts, err := service.ListAlertsBySeverity(tt.severity, tt.limit, tt.offset)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, alerts)
			} else {
				assert.NoError(t, err)
				if tt.validateAlerts != nil {
					tt.validateAlerts(t, alerts)
				}
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
