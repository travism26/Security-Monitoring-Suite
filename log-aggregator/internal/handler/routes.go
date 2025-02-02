package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/travism26/log-aggregator/internal/middleware"
)

// metricsHandler returns basic metrics in Prometheus format
func metricsHandler(c *gin.Context) {
	metrics := `
# HELP log_aggregator_up Indicates if the log-aggregator is up
# TYPE log_aggregator_up gauge
log_aggregator_up 1
`
	c.String(http.StatusOK, metrics)
}

// readinessHandler indicates service readiness
func readinessHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ready"})
}

// APIConfig holds API configuration
type APIConfig struct {
	APIKeys []string
}

// RegisterAPIRoutes registers all API routes with middleware
func RegisterAPIRoutes(r *gin.Engine, logHandler *LogHandler, alertHandler *AlertHandler, config APIConfig) {
	// Apply global middleware
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.RequestID())
	r.Use(middleware.CORS())

	// Register health and metrics endpoints
	r.GET("/metrics", metricsHandler)
	r.GET("/readiness", readinessHandler)

	// Create API group with version prefix
	api := r.Group("/logs/api/v1")

	// Apply authentication and rate limiting to API routes
	api.Use(middleware.APIKeyAuth(config.APIKeys))
	api.Use(middleware.RateLimit(100, time.Minute)) // 100 requests per minute

	{
		// Log routes - removed /logs since it's already in the base path
		logs := api.Group("")
		{
			logs.GET("", logHandler.ListLogs)
			logs.GET("/:id", logHandler.GetLog)
			logs.GET("/range", logHandler.ListLogsByTimeRange)
			logs.POST("", logHandler.StoreLog)
			logs.POST("/batch", logHandler.StoreBatchLogs)
		}

		// Alert routes
		alerts := api.Group("/alerts")
		{
			alerts.GET("", alertHandler.ListAlerts)
			alerts.GET("/:id", alertHandler.GetAlert)
			alerts.GET("/status/:status", alertHandler.ListAlertsByStatus)
			alerts.GET("/severity/:severity", alertHandler.ListAlertsBySeverity)
			alerts.GET("/trends", alertHandler.GetAlertTrends)
			alerts.PUT("/:id/status", alertHandler.UpdateAlertStatus)
		}

	}
}
