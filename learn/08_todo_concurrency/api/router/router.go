package router

import (
	"todo_concurrency/api/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes
func SetupRouter(
	todoHandler *handler.TodoHandler,
	batchHandler *handler.BatchHandler,
	statsHandler *handler.StatsHandler,
	notifyHandler *handler.NotifyHandler,
	adminHandler *handler.AdminHandler,
) *gin.Engine {
	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Todo Concurrency Learning API",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Basic CRUD operations (demonstrates interfaces)
		todos := v1.Group("/todos")
		{
			todos.POST("", todoHandler.Create)           // Create todo
			todos.GET("", todoHandler.GetAll)             // List todos
			todos.GET("/:id", todoHandler.GetByID)        // Get specific todo
			todos.PUT("/:id", todoHandler.Update)         // Update todo
			todos.DELETE("/:id", todoHandler.Delete)      // Delete todo
			todos.PATCH("/:id/toggle", todoHandler.ToggleComplete) // Toggle completion

			// Notification endpoint (goroutines + channels)
			todos.POST("/:id/notify", notifyHandler.SendNotification)

			// Batch operations (goroutines + channels + worker pool)
			todos.POST("/batch", batchHandler.ProcessBatch)
			todos.POST("/batch-v2", batchHandler.ProcessBatchV2)
		}

		// Statistics endpoints (mutex for thread-safe counters)
		stats := v1.Group("/stats")
		{
			stats.GET("", statsHandler.GetStats)                 // Basic stats
			stats.GET("/detailed", statsHandler.GetDetailedStats) // Detailed stats
			stats.GET("/storage", statsHandler.GetStorageStats)  // Storage-specific stats
			stats.GET("/goroutines", statsHandler.GetGoroutineCount) // Goroutine count
			stats.POST("/reset", statsHandler.ResetStats)        // Reset stats
		}

		// Notification statistics
		notifications := v1.Group("/notifications")
		{
			notifications.GET("/stats", notifyHandler.GetNotificationStats)
		}

		// Admin endpoints (interface switching)
		admin := v1.Group("/admin")
		{
			admin.POST("/switch-storage", adminHandler.SwitchStorage)  // Switch storage backend
			admin.GET("/storage-info", adminHandler.GetCurrentStorage) // Get current storage info
		}
	}

	return r
}
