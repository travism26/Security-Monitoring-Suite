package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/travism26/log-aggregator/internal/domain"
	"github.com/travism26/log-aggregator/internal/service"
)

// Response structures for consistent API responses
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
	} `json:"meta"`
}

type LogHandler struct {
	logService *service.LogService
}

func NewLogHandler(logService *service.LogService) *LogHandler {
	return &LogHandler{
		logService: logService,
	}
}

// GetLog godoc
// @Summary Get a log by ID
// @Description Retrieve a specific log entry by its ID
// @Tags logs
// @Accept json
// @Produce json
// @Param id path string true "Log ID"
// @Success 200 {object} Response
// @Failure 404 {object} Response
// @Router /logs/{id} [get]
func (h *LogHandler) GetLog(c *gin.Context) {
	id := c.Param("id")
	userID := c.GetString("user_id") // Get user_id from context (set by auth middleware)
	if userID == "" {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	log, err := h.logService.GetLog(userID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Log not found",
		})
		return
	}
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    log,
	})
}

// ListLogs godoc
// @Summary List logs with pagination
// @Description Retrieve a list of logs with pagination support
// @Tags logs
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page" default(10)
// @Param offset query int false "Number of items to skip" default(0)
// @Param user_id query string false "Filter logs by user ID"
// @Success 200 {object} PaginatedResponse
// @Failure 500 {object} Response
// @Router /logs [get]
func (h *LogHandler) ListLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	userID := c.Query("user_id")

	// Validate pagination parameters
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	logs, err := h.logService.ListLogs(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve logs",
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    logs,
		Meta: struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}{
			Limit:  limit,
			Offset: offset,
		},
	})
}

// ListLogsByTimeRange godoc
// @Summary List logs within a time range
// @Description Retrieve logs within a specified time range with pagination
// @Tags logs
// @Accept json
// @Produce json
// @Param start_time query string true "Start time (RFC3339)"
// @Param end_time query string true "End time (RFC3339)"
// @Param limit query int false "Number of items per page" default(10)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} PaginatedResponse
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /logs/range [get]
func (h *LogHandler) ListLogsByTimeRange(c *gin.Context) {
	startStr := c.Query("start_time")
	endStr := c.Query("end_time")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	userID := c.GetString("user_id") // Get user_id from context

	if userID == "" {
		c.JSON(http.StatusUnauthorized, Response{
			Success: false,
			Error:   "User not authenticated",
		})
		return
	}

	// Validate pagination parameters
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	// Parse time parameters
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid start_time format. Expected RFC3339",
		})
		return
	}

	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid end_time format. Expected RFC3339",
		})
		return
	}

	logs, err := h.logService.ListByTimeRange(userID, start, end, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve logs",
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    logs,
		Meta: struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}{
			Limit:  limit,
			Offset: offset,
		},
	})
}

// StoreLog godoc
// @Summary Store a new log
// @Description Store a single log entry
// @Tags logs
// @Accept json
// @Produce json
// @Param log body domain.Log true "Log object"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Router /logs [post]
func (h *LogHandler) StoreLog(c *gin.Context) {
	var log domain.Log
	if err := c.ShouldBindJSON(&log); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid log format",
		})
		return
	}

	if err := h.logService.StoreLog(&log); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to store log",
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    log,
	})
}

// StoreBatchLogs godoc
// @Summary Store multiple logs
// @Description Store multiple log entries in a single request
// @Tags logs
// @Accept json
// @Produce json
// @Param logs body []domain.Log true "Array of log objects"
// @Success 201 {object} Response
// @Failure 400 {object} Response
// @Router /logs/batch [post]
func (h *LogHandler) StoreBatchLogs(c *gin.Context) {
	var logs []*domain.Log
	if err := c.ShouldBindJSON(&logs); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid logs format",
		})
		return
	}

	if err := h.logService.StoreBatch(logs); err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to store logs",
		})
		return
	}

	c.JSON(http.StatusCreated, Response{
		Success: true,
		Data:    logs,
	})
}

func RegisterRoutes(r *gin.Engine, h *LogHandler) {
	// GET endpoints
	r.GET("/logs/:id", h.GetLog)
	r.GET("/logs", h.ListLogs)
	r.GET("/logs/range", h.ListLogsByTimeRange)

	// POST endpoints
	r.POST("/logs", h.StoreLog)
	r.POST("/logs/batch", h.StoreBatchLogs)
}
