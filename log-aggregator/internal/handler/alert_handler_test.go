package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/travism26/log-aggregator/internal/domain"
	"github.com/travism26/log-aggregator/internal/service"
)

type MockAlertRepository struct {
	mock.Mock
}

func (m *MockAlertRepository) Store(alert *domain.Alert) error {
	args := m.Called(alert)
	return args.Error(0)
}

func (m *MockAlertRepository) FindByID(orgID, id string) (*domain.Alert, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Alert), args.Error(1)
}

func (m *MockAlertRepository) List(orgID string, limit, offset int) ([]*domain.Alert, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]*domain.Alert), args.Error(1)
}

func (m *MockAlertRepository) Update(alert *domain.Alert) error {
	args := m.Called(alert)
	return args.Error(0)
}

func (m *MockAlertRepository) FindByStatus(orgID string, status domain.AlertStatus, limit, offset int) ([]*domain.Alert, error) {
	args := m.Called(status, limit, offset)
	return args.Get(0).([]*domain.Alert), args.Error(1)
}

func (m *MockAlertRepository) FindBySeverity(orgID string, severity domain.AlertSeverity, limit, offset int) ([]*domain.Alert, error) {
	args := m.Called(severity, limit, offset)
	return args.Get(0).([]*domain.Alert), args.Error(1)
}

func (m *MockAlertRepository) CountBySeverity(orgID string, severity domain.AlertSeverity) (int64, error) {
	args := m.Called(orgID, severity)
	return int64(1), args.Error(1)
}

func (m *MockAlertRepository) ListByTimeRange(orgID string, start, end time.Time, limit, offset int) ([]*domain.Alert, error) {
	args := m.Called(orgID, start, end, limit, offset)
	return args.Get(0).([]*domain.Alert), args.Error(1)
}

func (m *MockAlertRepository) CountBySource(orgID string, source string) (int64, error) {
	return 1, nil
}

func (m *MockAlertRepository) CountByStatus(orgID string, status domain.AlertStatus) (int64, error) {
	return 1, nil
}

func (m *MockAlertRepository) CountByTimeRange(orgID string, start, end time.Time) (int64, error) {
	return 1, nil
}

func (m *MockAlertRepository) Delete(orgID, id string) error {
	return nil
}

func (m *MockAlertRepository) FindBySource(orgID string, source string, limit, offset int) ([]*domain.Alert, error) {
	args := m.Called(orgID, source, limit, offset)
	return args.Get(0).([]*domain.Alert), args.Error(1)
}

/*
AlertRepo

	Store(alert *Alert) error
	Update(alert *Alert) error
	Delete(orgID, id string) error

	// Retrieval operations
	FindByID(orgID, id string) (*Alert, error)
	List(orgID string, limit, offset int) ([]*Alert, error)

	// Status-based queries
	FindByStatus(orgID string, status AlertStatus, limit, offset int) ([]*Alert, error)
	CountByStatus(orgID string, status AlertStatus) (int64, error)

	// Severity-based queries
	FindBySeverity(orgID string, severity AlertSeverity, limit, offset int) ([]*Alert, error)
	CountBySeverity(orgID string, severity AlertSeverity) (int64, error)

	// Time-based queries
	ListByTimeRange(orgID string, start, end time.Time, limit, offset int) ([]*Alert, error)
	CountByTimeRange(orgID string, start, end time.Time) (int64, error)

	// Source-based queries
	FindBySource(orgID string, source string, limit, offset int) ([]*Alert, error)
	CountBySource(orgID string, source string) (int64, error)
*/
func setupTestRouter() (*gin.Engine, *MockAlertRepository, *service.AlertService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)
	mockRepo := new(MockAlertRepository)
	alertService := service.NewAlertService(mockRepo, &service.AlertServiceConfig{
		SystemMemory: 16 * 1024 * 1024 * 1024,
		TimeNowFn: func() time.Time {
			return fixedTime
		},
	})

	handler := NewAlertHandler(alertService)
	RegisterAlertRoutes(router, handler)
	return router, mockRepo, alertService
}

func TestGetAlert(t *testing.T) {
	router, mockRepo, _ := setupTestRouter()

	t.Run("Success", func(t *testing.T) {
		alert := &domain.Alert{
			ID:          "123",
			Title:       "Test Alert",
			Description: "Test Description",
			Status:      domain.StatusOpen,
		}

		mockRepo.On("FindByID", "123").Return(alert, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/alerts/123", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool                   `json:"success"`
			Data    map[string]interface{} `json:"data"`
			Error   string                 `json:"error,omitempty"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		assert.Equal(t, alert.ID, response.Data["id"])
		assert.Equal(t, alert.Title, response.Data["title"])
		assert.Equal(t, alert.Description, response.Data["description"])
		assert.Equal(t, string(alert.Status), response.Data["status"])
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.On("FindByID", "999").Return(nil, fmt.Errorf("not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/alerts/999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestListAlerts(t *testing.T) {
	router, mockRepo, _ := setupTestRouter()

	t.Run("Success with default pagination", func(t *testing.T) {
		alerts := []*domain.Alert{
			{
				ID:          "1",
				Title:       "Alert 1",
				Description: "Description 1",
				Status:      domain.StatusOpen,
			},
			{
				ID:          "2",
				Title:       "Alert 2",
				Description: "Description 2",
				Status:      domain.StatusResolved,
			},
		}

		mockRepo.On("List", 10, 0).Return(alerts, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/alerts", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool          `json:"success"`
			Data    []interface{} `json:"data"`
			Meta    struct {
				Limit  int `json:"limit"`
				Offset int `json:"offset"`
			} `json:"meta"`
			Error string `json:"error,omitempty"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Len(t, response.Data, 2)

		alert1 := response.Data[0].(map[string]interface{})
		assert.Equal(t, alerts[0].ID, alert1["id"])
		assert.Equal(t, alerts[0].Title, alert1["title"])

		alert2 := response.Data[1].(map[string]interface{})
		assert.Equal(t, alerts[1].ID, alert2["id"])
		assert.Equal(t, alerts[1].Title, alert2["title"])
	})
}

func TestListAlertsByStatus(t *testing.T) {
	router, mockRepo, _ := setupTestRouter()

	t.Run("Success", func(t *testing.T) {
		alerts := []*domain.Alert{
			{
				ID:          "1",
				Title:       "Alert 1",
				Description: "Description 1",
				Status:      domain.StatusOpen,
			},
		}

		mockRepo.On("FindByStatus", domain.StatusOpen, 10, 0).Return(alerts, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/alerts/status/OPEN", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool          `json:"success"`
			Data    []interface{} `json:"data"`
			Meta    struct {
				Limit  int `json:"limit"`
				Offset int `json:"offset"`
			} `json:"meta"`
			Error string `json:"error,omitempty"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Len(t, response.Data, 1)

		alert1 := response.Data[0].(map[string]interface{})
		assert.Equal(t, alerts[0].ID, alert1["id"])
		assert.Equal(t, alerts[0].Title, alert1["title"])
		assert.Equal(t, string(alerts[0].Status), alert1["status"])
	})

	t.Run("Invalid Status", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/alerts/status/INVALID", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestListAlertsBySeverity(t *testing.T) {
	router, mockRepo, _ := setupTestRouter()

	t.Run("Success", func(t *testing.T) {
		alerts := []*domain.Alert{
			{
				ID:          "1",
				Title:       "Alert 1",
				Description: "Description 1",
				Severity:    domain.SeverityHigh,
			},
		}

		mockRepo.On("FindBySeverity", domain.SeverityHigh, 10, 0).Return(alerts, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/alerts/severity/HIGH", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool          `json:"success"`
			Data    []interface{} `json:"data"`
			Meta    struct {
				Limit  int `json:"limit"`
				Offset int `json:"offset"`
			} `json:"meta"`
			Error string `json:"error,omitempty"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Len(t, response.Data, 1)

		alert1 := response.Data[0].(map[string]interface{})
		assert.Equal(t, alerts[0].ID, alert1["id"])
		assert.Equal(t, alerts[0].Title, alert1["title"])
		assert.Equal(t, string(alerts[0].Severity), alert1["severity"])
	})

	t.Run("Invalid Severity", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/alerts/severity/INVALID", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestUpdateAlertStatus(t *testing.T) {
	router, mockRepo, _ := setupTestRouter()

	t.Run("Success - Resolve Alert", func(t *testing.T) {
		alert := &domain.Alert{
			ID:     "123",
			Status: domain.StatusOpen,
		}

		mockRepo.On("FindByID", "123").Return(alert, nil)
		mockRepo.On("Update", mock.AnythingOfType("*domain.Alert")).Return(nil)

		// Create request body
		body := map[string]interface{}{
			"status": "RESOLVED",
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/alerts/123/status", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool        `json:"success"`
			Data    interface{} `json:"data"`
			Error   string      `json:"error,omitempty"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("Invalid Status", func(t *testing.T) {
		// Create request body with invalid status
		body := map[string]interface{}{
			"status": "INVALID",
		}
		jsonBody, _ := json.Marshal(body)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("PUT", "/api/v1/alerts/123/status", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestGetAlertTrends(t *testing.T) {
	router, mockRepo, _ := setupTestRouter()
	fixedTime := time.Date(2023, 1, 1, 12, 0, 0, 0, time.UTC)

	t.Run("Success", func(t *testing.T) {
		alerts := []*domain.Alert{
			{
				ID:        "1",
				Severity:  domain.SeverityHigh,
				Status:    domain.StatusOpen,
				Source:    "host1",
				CreatedAt: fixedTime.Add(-1 * time.Hour),
			},
		}

		mockRepo.On("List", 1000, 0).Return(alerts, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/alerts/trends?start_time=2023-01-01T11:00:00Z&end_time=2023-01-01T13:00:00Z", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool                   `json:"success"`
			Data    map[string]interface{} `json:"data"`
			Error   string                 `json:"error,omitempty"`
		}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		trends := response.Data
		assert.Equal(t, float64(1), trends["total_alerts"])

		severityMap := trends["alerts_by_severity"].(map[string]interface{})
		assert.Equal(t, float64(1), severityMap["HIGH"])

		statusMap := trends["alerts_by_status"].(map[string]interface{})
		assert.Equal(t, float64(1), statusMap["OPEN"])
	})

	t.Run("Invalid Time Format", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/v1/alerts/trends?start_time=invalid&end_time=invalid", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
