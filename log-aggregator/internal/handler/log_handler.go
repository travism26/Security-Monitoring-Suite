package handler

import (
	"net/http"
	"strconv"

	"github.com/travism26/log-aggregator/internal/service"

	"github.com/gin-gonic/gin"
)

type LogHandler struct {
	logService *service.LogService
}

func NewLogHandler(logService *service.LogService) *LogHandler {
	return &LogHandler{
		logService: logService,
	}
}

func (h *LogHandler) GetLog(c *gin.Context) {
	id := c.Param("id")
	log, err := h.logService.GetLog(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Log not found"})
		return
	}
	c.JSON(http.StatusOK, log)
}

func (h *LogHandler) ListLogs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	logs, err := h.logService.ListLogs(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve logs"})
		return
	}
	c.JSON(http.StatusOK, logs)
}

func RegisterRoutes(r *gin.Engine, h *LogHandler) {
	r.GET("/logs/:id", h.GetLog)
	r.GET("/logs", h.ListLogs)
}
