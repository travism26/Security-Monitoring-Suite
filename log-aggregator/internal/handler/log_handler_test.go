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

// MockLogRepository is a mock implementation of the LogRepository interface
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

func (m *MockLogRepository) FindByID(organization_id, id string) (*domain.Log, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Log), args.Error(1)
}

func (m *MockLogRepository) List(organization_id string, limit, offset int) ([]*domain.Log, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]*domain.Log), args.Error(1)
}

func (m *MockLogRepository) ListByTimeRange(organization_id string, start, end time.Time, limit, offset int) ([]*domain.Log, error) {
	args := m.Called(start, end, limit, offset)
	return args.Get(0).([]*domain.Log), args.Error(1)
}

func (m *MockLogRepository) CountByTimeRange(orgID string, start, end time.Time) (int64, error) {
	return 1, nil
}

func (m *MockLogRepository) ListByHost(orgID string, host string, limit, offset int) ([]*domain.Log, error) {
	args := m.Called(orgID, host, limit, offset)
	return args.Get(0).([]*domain.Log), args.Error(1)
}

func (m *MockLogRepository) ListByLevel(orgID string, level string, limit, offset int) ([]*domain.Log, error) {
	args := m.Called(orgID, level, limit, offset)
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

func setupTest() (*gin.Engine, *MockLogRepository) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	mockRepo := new(MockLogRepository)
	logService := service.NewLogService(mockRepo, service.LogServiceConfig{
		Environment:         "test",
		Application:         "log-aggregator",
		Component:           "api",
		MultiTenancyEnabled: false,
	})
	handler := NewLogHandler(logService)
	RegisterRoutes(r, handler)
	return r, mockRepo
}

func TestGetLog(t *testing.T) {
	r, mockRepo := setupTest()

	t.Run("Success", func(t *testing.T) {
		expectedLog := &domain.Log{
			ID:        "123",
			Message:   "Test log",
			Timestamp: time.Now(),
			Level:     "INFO",
			Host:      "test-host",
		}

		mockRepo.On("FindByID", "123").Return(expectedLog, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/logs/123", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)

		logData, err := json.Marshal(response.Data)
		assert.NoError(t, err)
		var returnedLog domain.Log
		err = json.Unmarshal(logData, &returnedLog)
		assert.NoError(t, err)
		assert.Equal(t, expectedLog.ID, returnedLog.ID)
	})

	t.Run("Not Found", func(t *testing.T) {
		mockRepo.On("FindByID", "999").Return(nil, fmt.Errorf("not found"))

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/logs/999", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Log not found", response.Error)
	})
}

func TestListLogs(t *testing.T) {
	r, mockRepo := setupTest()

	t.Run("Success with default pagination", func(t *testing.T) {
		expectedLogs := []*domain.Log{
			{ID: "1", Message: "Log 1"},
			{ID: "2", Message: "Log 2"},
		}

		mockRepo.On("List", 10, 0).Return(expectedLogs, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/logs", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, 10, response.Meta.Limit)
		assert.Equal(t, 0, response.Meta.Offset)
	})

	t.Run("Success with custom pagination", func(t *testing.T) {
		expectedLogs := []*domain.Log{
			{ID: "3", Message: "Log 3"},
		}

		mockRepo.On("List", 5, 10).Return(expectedLogs, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/logs?limit=5&offset=10", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.Equal(t, 5, response.Meta.Limit)
		assert.Equal(t, 10, response.Meta.Offset)
	})
}

func TestListLogsByTimeRange(t *testing.T) {
	r, mockRepo := setupTest()

	t.Run("Success", func(t *testing.T) {
		start := time.Now().Add(-24 * time.Hour)
		end := time.Now()
		expectedLogs := []*domain.Log{
			{ID: "1", Message: "Log 1", Timestamp: start.Add(time.Hour)},
			{ID: "2", Message: "Log 2", Timestamp: start.Add(2 * time.Hour)},
		}

		mockRepo.On("ListByTimeRange", mock.Anything, mock.Anything, 10, 0).Return(expectedLogs, nil)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", fmt.Sprintf("/logs/range?start_time=%s&end_time=%s",
			start.Format(time.RFC3339), end.Format(time.RFC3339)), nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response PaginatedResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("Invalid Time Format", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/logs/range?start_time=invalid&end_time=invalid", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Contains(t, response.Error, "Invalid start_time format")
	})
}

func TestStoreLog(t *testing.T) {
	r, mockRepo := setupTest()

	t.Run("Success", func(t *testing.T) {
		log := domain.Log{
			Message:   "Test log",
			Level:     "INFO",
			Host:      "test-host",
			Timestamp: time.Now(),
		}

		mockRepo.On("Store", mock.Anything).Return(nil)

		logJSON, _ := json.Marshal(log)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/logs", bytes.NewBuffer(logJSON))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/logs", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid log format", response.Error)
	})
}

func TestStoreBatchLogs(t *testing.T) {
	r, mockRepo := setupTest()

	t.Run("Success", func(t *testing.T) {
		logs := []*domain.Log{
			{Message: "Log 1", Level: "INFO", Host: "test-host-1", Timestamp: time.Now()},
			{Message: "Log 2", Level: "ERROR", Host: "test-host-2", Timestamp: time.Now()},
		}

		mockRepo.On("StoreBatch", mock.Anything).Return(nil)

		logsJSON, _ := json.Marshal(logs)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/logs/batch", bytes.NewBuffer(logsJSON))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/logs/batch", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response Response
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "Invalid logs format", response.Error)
	})
}
