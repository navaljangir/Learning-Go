package routes

import (
	"gin_server/handlers"
	"gin_server/middlewares"

	"github.com/gin-gonic/gin"
)

// SetupRoutes configures all routes for the application
func SetupRoutes(router *gin.Engine) {
	// Health check (public)
	router.GET("/health", handlers.HealthCheck)

	// API group
	api := router.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// User routes (protected)
		users := api.Group("/users")
		users.Use(middlewares.AuthMiddleware()) // Apply auth middleware
		{
			users.GET("/profile", handlers.GetProfile)
			users.GET("", handlers.GetAllUsers)
		}
	}
}
