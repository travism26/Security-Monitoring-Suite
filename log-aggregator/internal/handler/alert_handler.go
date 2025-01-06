package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/travism26/log-aggregator/internal/domain"
	"github.com/travism26/log-aggregator/internal/service"
)

type AlertHandler struct {
	alertService *service.AlertService
}

func NewAlertHandler(alertService *service.AlertService) *AlertHandler {
	return &AlertHandler{
		alertService: alertService,
	}
}

// GetAlert godoc
// @Summary Get an alert by ID
// @Description Retrieve a specific alert by its ID
// @Tags alerts
// @Accept json
// @Produce json
// @Param id path string true "Alert ID"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Router /alerts/{id} [get]
func (h *AlertHandler) GetAlert(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Alert ID is required",
		})
		return
	}

	alert, err := h.alertService.GetAlert(id)
	if err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Alert not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    alert,
	})
}

// ListAlerts godoc
// @Summary List alerts with pagination
// @Description Retrieve a list of alerts with pagination support
// @Tags alerts
// @Accept json
// @Produce json
// @Param limit query int false "Number of items per page" default(10)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} PaginatedResponse
// @Failure 500 {object} Response
// @Router /alerts [get]
func (h *AlertHandler) ListAlerts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Validate pagination parameters
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	alerts, err := h.alertService.ListAlerts(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve alerts",
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    alerts,
		Meta: struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}{
			Limit:  limit,
			Offset: offset,
		},
	})
}

// ListAlertsByStatus godoc
// @Summary List alerts by status
// @Description Retrieve alerts filtered by status with pagination
// @Tags alerts
// @Accept json
// @Produce json
// @Param status path string true "Alert status (OPEN, RESOLVED, IGNORED)"
// @Param limit query int false "Number of items per page" default(10)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} PaginatedResponse
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /alerts/status/{status} [get]
func (h *AlertHandler) ListAlertsByStatus(c *gin.Context) {
	status := domain.AlertStatus(c.Param("status"))
	if !status.IsValid() {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid alert status. Must be one of: OPEN, RESOLVED, IGNORED",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Validate pagination parameters
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	alerts, err := h.alertService.ListAlertsByStatus(status, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve alerts",
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    alerts,
		Meta: struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}{
			Limit:  limit,
			Offset: offset,
		},
	})
}

// ListAlertsBySeverity godoc
// @Summary List alerts by severity
// @Description Retrieve alerts filtered by severity with pagination
// @Tags alerts
// @Accept json
// @Produce json
// @Param severity path string true "Alert severity (LOW, MEDIUM, HIGH, CRITICAL)"
// @Param limit query int false "Number of items per page" default(10)
// @Param offset query int false "Number of items to skip" default(0)
// @Success 200 {object} PaginatedResponse
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /alerts/severity/{severity} [get]
func (h *AlertHandler) ListAlertsBySeverity(c *gin.Context) {
	severity := domain.AlertSeverity(c.Param("severity"))
	if !severity.IsValid() {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid alert severity. Must be one of: LOW, MEDIUM, HIGH, CRITICAL",
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	// Validate pagination parameters
	if limit < 1 || limit > 100 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	alerts, err := h.alertService.ListAlertsBySeverity(severity, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve alerts",
		})
		return
	}

	c.JSON(http.StatusOK, PaginatedResponse{
		Success: true,
		Data:    alerts,
		Meta: struct {
			Limit  int `json:"limit"`
			Offset int `json:"offset"`
		}{
			Limit:  limit,
			Offset: offset,
		},
	})
}

// UpdateAlertStatus godoc
// @Summary Update alert status
// @Description Update the status of a specific alert
// @Tags alerts
// @Accept json
// @Produce json
// @Param id path string true "Alert ID"
// @Param status body object true "Status update object"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Router /alerts/{id}/status [put]
func (h *AlertHandler) UpdateAlertStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Alert ID is required",
		})
		return
	}

	var request struct {
		Status domain.AlertStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid request format",
		})
		return
	}

	if !request.Status.IsValid() {
		c.JSON(http.StatusBadRequest, Response{
			Success: false,
			Error:   "Invalid alert status. Must be one of: OPEN, RESOLVED, IGNORED",
		})
		return
	}

	if err := h.alertService.UpdateAlertStatus(id, request.Status); err != nil {
		c.JSON(http.StatusNotFound, Response{
			Success: false,
			Error:   "Alert not found",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    gin.H{"message": "Alert status updated successfully"},
	})
}

// GetAlertTrends godoc
// @Summary Get alert trends
// @Description Get alert trends over a specified time period
// @Tags alerts
// @Accept json
// @Produce json
// @Param start_time query string true "Start time (RFC3339)"
// @Param end_time query string true "End time (RFC3339)"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Router /alerts/trends [get]
func (h *AlertHandler) GetAlertTrends(c *gin.Context) {
	startStr := c.Query("start_time")
	endStr := c.Query("end_time")

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

	trends, err := h.alertService.GetAlertTrends(start, end)
	if err != nil {
		c.JSON(http.StatusInternalServerError, Response{
			Success: false,
			Error:   "Failed to retrieve alert trends",
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    trends,
	})
}

func RegisterAlertRoutes(r *gin.Engine, h *AlertHandler) {
	api := r.Group("/api/v1")
	{
		// GET endpoints
		api.GET("/alerts/:id", h.GetAlert)
		api.GET("/alerts", h.ListAlerts)
		api.GET("/alerts/status/:status", h.ListAlertsByStatus)
		api.GET("/alerts/severity/:severity", h.ListAlertsBySeverity)
		api.GET("/alerts/trends", h.GetAlertTrends)

		// PUT endpoints
		api.PUT("/alerts/:id/status", h.UpdateAlertStatus)
	}
}
