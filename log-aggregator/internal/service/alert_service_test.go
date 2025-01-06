package service

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/travism26/log-aggregator/internal/domain"
)

type MockAlertRepository struct {
	mock.Mock
}

func (m *MockAlertRepository) Store(alert *domain.Alert) error {
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
	if args.Get(0) == nil || args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Alert), nil
}

func (m *MockAlertRepository) Update(alert *domain.Alert) error {
	args := m.Called(alert)
	return args.Error(0)
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

func TestProcessMetrics(t *testing.T) {
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	mockRepo := new(MockAlertRepository)
	service := NewAlertService(mockRepo, &AlertServiceConfig{
		SystemMemory: 16 * 1024 * 1024 * 1024,
		TimeNowFn: func() time.Time {
			return fixedTime
		},
	})

	t.Run("High CPU Usage", func(t *testing.T) {
		log := &domain.Log{
			ID:              "log1",
			Host:            "test-host",
			TotalCPUPercent: 90.0, // Above default threshold of 80%
		}

		mockRepo.On("Store", mock.MatchedBy(func(alert *domain.Alert) bool {
			return alert.Source == "test-host" &&
				alert.Severity == domain.SeverityHigh &&
				alert.Status == domain.StatusOpen &&
				len(alert.RelatedLogs) == 1 &&
				alert.RelatedLogs[0] == "log1" &&
				alert.CreatedAt == fixedTime &&
				alert.UpdatedAt == fixedTime
		})).Return(nil)

		err := service.ProcessMetrics(log)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("High Memory Usage", func(t *testing.T) {
		log := &domain.Log{
			ID:               "log2",
			Host:             "test-host",
			TotalMemoryUsage: 14 * 1024 * 1024 * 1024, // ~87.5% of 16GB
		}

		mockRepo.On("Store", mock.MatchedBy(func(alert *domain.Alert) bool {
			return alert.Source == "test-host" &&
				alert.Severity == domain.SeverityHigh &&
				alert.Status == domain.StatusOpen &&
				len(alert.RelatedLogs) == 1 &&
				alert.RelatedLogs[0] == "log2" &&
				alert.CreatedAt == fixedTime &&
				alert.UpdatedAt == fixedTime
		})).Return(nil)

		err := service.ProcessMetrics(log)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("High Process Count", func(t *testing.T) {
		log := &domain.Log{
			ID:           "log3",
			Host:         "test-host",
			ProcessCount: 1200, // Above default threshold of 1000
		}

		mockRepo.On("Store", mock.MatchedBy(func(alert *domain.Alert) bool {
			return alert.Source == "test-host" &&
				alert.Severity == domain.SeverityMedium &&
				alert.Status == domain.StatusOpen &&
				len(alert.RelatedLogs) == 1 &&
				alert.RelatedLogs[0] == "log3" &&
				alert.CreatedAt == fixedTime &&
				alert.UpdatedAt == fixedTime
		})).Return(nil)

		err := service.ProcessMetrics(log)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("No Thresholds Exceeded", func(t *testing.T) {
		log := &domain.Log{
			ID:               "log4",
			Host:             "test-host",
			TotalCPUPercent:  70.0,
			TotalMemoryUsage: 12 * 1024 * 1024 * 1024,
			ProcessCount:     800,
		}

		err := service.ProcessMetrics(log)
		assert.NoError(t, err)
		// No Store calls expected as no thresholds were exceeded
		mockRepo.AssertNotCalled(t, "Store")
	})
}

func TestSetThresholds(t *testing.T) {
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	mockRepo := new(MockAlertRepository)
	service := NewAlertService(mockRepo, &AlertServiceConfig{
		SystemMemory: 16 * 1024 * 1024 * 1024,
		TimeNowFn: func() time.Time {
			return fixedTime
		},
	})

	customThresholds := AlertThresholds{
		CPUUsagePercent:    70.0,
		MemoryUsagePercent: 75.0,
		ProcessCount:       800,
	}

	t.Run("Custom CPU Threshold", func(t *testing.T) {
		service.SetThresholds(customThresholds)

		log := &domain.Log{
			ID:              "log1",
			Host:            "test-host",
			TotalCPUPercent: 75.0, // Above custom threshold of 70%
		}

		mockRepo.On("Store", mock.MatchedBy(func(alert *domain.Alert) bool {
			return alert.Source == "test-host" &&
				alert.Severity == domain.SeverityHigh &&
				alert.Status == domain.StatusOpen &&
				alert.CreatedAt == fixedTime &&
				alert.UpdatedAt == fixedTime
		})).Return(nil)

		err := service.ProcessMetrics(log)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateAlertStatus(t *testing.T) {
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	mockRepo := new(MockAlertRepository)
	service := NewAlertService(mockRepo, &AlertServiceConfig{
		SystemMemory: 16 * 1024 * 1024 * 1024,
		TimeNowFn: func() time.Time {
			return fixedTime
		},
	})

	t.Run("Resolve Alert", func(t *testing.T) {
		alert := &domain.Alert{
			ID:          "123",
			Title:       "Test Alert",
			Description: "Test Description",
			Severity:    domain.SeverityHigh,
			Status:      domain.StatusOpen,
			Source:      "test-host",
			CreatedAt:   fixedTime.Add(-1 * time.Hour),
			UpdatedAt:   fixedTime.Add(-1 * time.Hour),
		}

		mockRepo.On("FindByID", "123").Return(alert, nil)
		var capturedAlert *domain.Alert
		mockRepo.On("Update", mock.MatchedBy(func(a *domain.Alert) bool {
			capturedAlert = a
			return true
		})).Return(nil)

		err := service.UpdateAlertStatus("123", domain.StatusResolved)
		assert.NoError(t, err)
		assert.Equal(t, alert.ID, capturedAlert.ID)
		assert.Equal(t, alert.Title, capturedAlert.Title)
		assert.Equal(t, alert.Description, capturedAlert.Description)
		assert.Equal(t, alert.Severity, capturedAlert.Severity)
		assert.Equal(t, alert.Source, capturedAlert.Source)
		assert.Equal(t, domain.StatusResolved, capturedAlert.Status)
		assert.Equal(t, fixedTime, *capturedAlert.ResolvedAt)
		assert.Equal(t, fixedTime, capturedAlert.UpdatedAt)
		assert.Equal(t, alert.CreatedAt, capturedAlert.CreatedAt)
	})

	t.Run("Alert Not Found", func(t *testing.T) {
		mockRepo.On("FindByID", "999").Return(nil, fmt.Errorf("not found"))

		err := service.UpdateAlertStatus("999", domain.StatusResolved)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAlertTrends(t *testing.T) {
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	mockRepo := new(MockAlertRepository)
	service := NewAlertService(mockRepo, &AlertServiceConfig{
		SystemMemory: 16 * 1024 * 1024 * 1024,
		TimeNowFn: func() time.Time {
			return fixedTime
		},
	})

	t.Run("Success", func(t *testing.T) {
		start := fixedTime.Add(-24 * time.Hour)
		end := fixedTime

		alerts := []*domain.Alert{
			{
				ID:        "1",
				Severity:  domain.SeverityHigh,
				Status:    domain.StatusOpen,
				Source:    "host1",
				CreatedAt: start.Add(time.Hour),
			},
			{
				ID:        "2",
				Severity:  domain.SeverityMedium,
				Status:    domain.StatusResolved,
				Source:    "host2",
				CreatedAt: start.Add(2 * time.Hour),
			},
			{
				ID:        "3",
				Severity:  domain.SeverityLow,
				Status:    domain.StatusIgnored,
				Source:    "host1",
				CreatedAt: start.Add(3 * time.Hour),
			},
		}

		mockRepo.On("List", 1000, 0).Return(alerts, nil)

		trends, err := service.GetAlertTrends(start, end)
		assert.NoError(t, err)
		assert.NotNil(t, trends)

		assert.Equal(t, 3, trends.TotalAlerts)
		assert.Equal(t, 1, trends.AlertsBySeverity["HIGH"])
		assert.Equal(t, 1, trends.AlertsBySeverity["MEDIUM"])
		assert.Equal(t, 1, trends.AlertsBySeverity["LOW"])
		assert.Equal(t, 1, trends.AlertsByStatus["OPEN"])
		assert.Equal(t, 1, trends.AlertsByStatus["RESOLVED"])
		assert.Equal(t, 1, trends.AlertsByStatus["IGNORED"])
		assert.Equal(t, 2, trends.TopSources["host1"])
		assert.Equal(t, 1, trends.TopSources["host2"])

		mockRepo.AssertExpectations(t)
	})

	// t.Run("Repository Error", func(t *testing.T) {
	// 	start := fixedTime.Add(-24 * time.Hour)
	// 	end := fixedTime

	// 	expectedErr := fmt.Errorf("database error")
	// 	mockRepo.On("List", 1000, 0).Return([]*domain.Alert(nil), expectedErr).Once()

	// 	trends, err := service.GetAlertTrends(start, end)
	// 	assert.Error(t, err)
	// 	assert.Equal(t, fmt.Sprintf("failed to retrieve alerts: %v", expectedErr), err.Error())
	// 	assert.Nil(t, trends)
	// 	mockRepo.AssertExpectations(t)
	// })
}
