package main

import (
	"log"
	"net/http"
	"todo_app/config"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.New()

	// Global middleware (order matters!)
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	v1 := r.Group("/api/v1")
	v1.Use()

	// Health check endpoint (public)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"service": "todo-api",
			"version": "1.0.0",
		})
	})

	return r
}

func main() {
	cfg := config.Load()
	log.Printf("Starting Server on port %s [%s]", cfg.Server.Port, cfg.Server.Environment)
	r := setupRouter()
	r.Run(":" + cfg.Server.Port)
}