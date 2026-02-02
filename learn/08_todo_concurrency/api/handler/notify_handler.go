package handler

import (
	"fmt"
	"todo_concurrency/internal/dto"
	"todo_concurrency/internal/service"
	"todo_concurrency/pkg/utils"

	"github.com/gin-gonic/gin"
)

// NotifyHandler handles notification endpoints
//
// GOROUTINES LEARNING:
// Demonstrates non-blocking asynchronous operations
type NotifyHandler struct {
	notifier *service.Notifier
}

// NewNotifyHandler creates a new notify handler
func NewNotifyHandler(notifier *service.Notifier) *NotifyHandler {
	return &NotifyHandler{
		notifier: notifier,
	}
}

// SendNotification handles POST /api/v1/todos/:id/notify
//
// LEARNING DEMONSTRATION:
// This endpoint returns IMMEDIATELY, but the notification is sent later!
// This is non-blocking async operation using goroutines.
//
// Try this:
// 1. Send a notification with 5 second delay
// 2. The HTTP response returns immediately
// 3. Watch the console - notification prints after 5 seconds
//
// This demonstrates: API doesn't block waiting for notification!
func (h *NotifyHandler) SendNotification(c *gin.Context) {
	todoID := c.Param("id")

	var req dto.NotifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	// Send async - this returns immediately!
	// GOROUTINE: The actual notification happens in background
	if err := h.notifier.SendAsync(c.Request.Context(), todoID, req.Message, req.DelaySeconds); err != nil {
		utils.RespondNotFound(c, err.Error())
		return
	}

	// Response sent immediately, notification happens later
	utils.RespondSuccess(c, 202, gin.H{
		"todo_id":       todoID,
		"message":       "Notification queued successfully",
		"delay_seconds": req.DelaySeconds,
		"explanation":   fmt.Sprintf("Notification will be sent in %d seconds. Watch the console!", req.DelaySeconds),
	}, "Notification queued")
}

// GetNotificationStats handles GET /api/v1/notifications/stats
//
// MUTEX & GOROUTINES:
// Shows stats about background notification processing
func (h *NotifyHandler) GetNotificationStats(c *gin.Context) {
	stats := h.notifier.GetStats()
	utils.RespondOK(c, stats)
}
