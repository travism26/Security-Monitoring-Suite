package handler

import (
	"github.com/gin-gonic/gin"
)

// RegisterAPIRoutes registers all API routes
func RegisterAPIRoutes(r *gin.Engine, logHandler *LogHandler, alertHandler *AlertHandler) {
	api := r.Group("/api/v1")
	{
		// Log routes
		logs := api.Group("/logs")
		{
			logs.GET("", logHandler.ListLogs)
			logs.GET("/:id", logHandler.GetLog)
		}

		// Alert routes
		alerts := api.Group("/alerts")
		{
			alerts.GET("", alertHandler.ListAlerts)
			alerts.GET("/:id", alertHandler.GetAlert)
			alerts.GET("/status/:status", alertHandler.ListAlertsByStatus)
			alerts.GET("/severity/:severity", alertHandler.ListAlertsBySeverity)
			alerts.PUT("/:id/status", alertHandler.UpdateAlertStatus)
		}
	}
}
