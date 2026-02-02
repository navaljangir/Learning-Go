package handler

import (
	"runtime"
	"todo_concurrency/internal/service"
	"todo_concurrency/pkg/utils"

	"github.com/gin-gonic/gin"
)

// StatsHandler handles statistics endpoints
//
// MUTEX LEARNING:
// StatsService uses mutex to protect counters from concurrent access
type StatsHandler struct {
	statsService *service.StatsService
}

// NewStatsHandler creates a new stats handler
func NewStatsHandler(statsService *service.StatsService) *StatsHandler {
	return &StatsHandler{
		statsService: statsService,
	}
}

// GetStats handles GET /api/v1/stats
//
// LEARNING:
// Try hitting this endpoint from multiple terminals simultaneously:
// while true; do curl http://localhost:8080/api/v1/stats; sleep 0.1; done
//
// The mutex ensures request_count increments correctly!
func (h *StatsHandler) GetStats(c *gin.Context) {
	stats, err := h.statsService.GetStats(c.Request.Context())
	if err != nil {
		utils.RespondInternalError(c, err.Error())
		return
	}

	utils.RespondOK(c, stats)
}

// GetDetailedStats handles GET /api/v1/stats/detailed
//
// MUTEX LEARNING:
// Multiple fields read atomically under read lock
func (h *StatsHandler) GetDetailedStats(c *gin.Context) {
	stats := h.statsService.GetDetailedStats()
	utils.RespondOK(c, stats)
}

// GetStorageStats handles GET /api/v1/stats/storage
//
// INTERFACE LEARNING:
// Returns storage-specific stats if the repository implements StorageInfo
func (h *StatsHandler) GetStorageStats(c *gin.Context) {
	stats := h.statsService.GetStorageStats()
	utils.RespondOK(c, stats)
}

// GetGoroutineCount handles GET /api/v1/stats/goroutines
//
// GOROUTINE LEARNING:
// Shows how many goroutines are currently running
// Try hitting the batch endpoint and check this - you'll see goroutines increase!
func (h *StatsHandler) GetGoroutineCount(c *gin.Context) {
	count := runtime.NumGoroutine()

	utils.RespondOK(c, gin.H{
		"goroutine_count": count,
		"explanation":     "This shows how many goroutines are currently running. Try the batch endpoint to see this number increase!",
	})
}

// ResetStats handles POST /api/v1/stats/reset
func (h *StatsHandler) ResetStats(c *gin.Context) {
	h.statsService.ResetStats()
	utils.RespondSuccess(c, 200, nil, "Statistics reset successfully")
}
