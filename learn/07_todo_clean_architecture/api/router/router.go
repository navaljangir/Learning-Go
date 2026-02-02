package router

import (
	"net/http"
	"todo_app/api/handler"
	"todo_app/api/middleware"
	"todo_app/pkg/utils"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all routes for the application
func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	todoHandler *handler.TodoHandler,
	jwtUtil *utils.JWTUtil,
) *gin.Engine {
	r := gin.New()

	// Global middleware
	r.Use(gin.Recovery())
	r.Use(middleware.ErrorRecoveryMiddleware())
	r.Use(middleware.LoggerMiddleware())
	r.Use(middleware.CORSMiddleware())

	// Health check endpoint (public)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "todo-api",
			"version": "1.0.0",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes (require authentication)
		authorized := v1.Group("")
		authorized.Use(middleware.AuthMiddleware(jwtUtil))
		{
			// User routes
			users := authorized.Group("/users")
			{
				users.GET("/profile", userHandler.GetProfile)
				users.PUT("/profile", userHandler.UpdateProfile)
			}

			// Todo routes
			todos := authorized.Group("/todos")
			{
				todos.GET("", todoHandler.List)
				todos.POST("", todoHandler.Create)
				todos.GET("/:id", todoHandler.GetByID)
				todos.PUT("/:id", todoHandler.Update)
				todos.PATCH("/:id/toggle", todoHandler.ToggleComplete)
				todos.DELETE("/:id", todoHandler.Delete)
			}
		}
	}

	return r
}
