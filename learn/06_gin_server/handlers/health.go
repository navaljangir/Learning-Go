package handlers

import (
	"gin_server/utils"
	"time"

	"github.com/gin-gonic/gin"
)

var startTime = time.Now()

// HealthCheck returns server health status
func HealthCheck(c *gin.Context) {
	utils.Success(c, gin.H{
		"status":  "healthy",
		"uptime":  time.Since(startTime).String(),
		"version": "1.0.0",
	})
}
