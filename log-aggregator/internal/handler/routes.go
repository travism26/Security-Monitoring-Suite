package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/travism26/log-aggregator/internal/middleware"
)

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

	// Create API group with version prefix
	api := r.Group("/api/v1")

	// Apply authentication and rate limiting to API routes
	api.Use(middleware.APIKeyAuth(config.APIKeys))
	api.Use(middleware.RateLimit(100, time.Minute)) // 100 requests per minute

	{
		// Log routes
		logs := api.Group("/logs")
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

		// Health check endpoint (no auth required)
		r.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"time":   time.Now().Format(time.RFC3339),
			})
		})
	}
}
