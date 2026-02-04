package handler

import (
	"fmt"
	"todo_concurrency/internal/dto"
	"todo_concurrency/internal/service"
	"todo_concurrency/pkg/utils"

	"github.com/gin-gonic/gin"
)

// BatchHandler handles batch operations using goroutines and channels
//
// GOROUTINES & CHANNELS LEARNING:
// This demonstrates concurrent processing of multiple items
type BatchHandler struct {
	batchProcessor *service.BatchProcessor
}

// NewBatchHandler creates a new batch handler
func NewBatchHandler(batchProcessor *service.BatchProcessor) *BatchHandler {
	return &BatchHandler{
		batchProcessor: batchProcessor,
	}
}

// ProcessBatch handles POST /api/v1/todos/batch
//
// LEARNING DEMONSTRATION:
// Try creating 10 todos - you'll see worker output showing concurrent processing!
// Workers process items in parallel, making this faster than sequential processing.
func (h *BatchHandler) ProcessBatch(c *gin.Context) {
	var req dto.BatchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	if len(req.Todos) == 0 {
		utils.RespondBadRequest(c, "at least one todo is required")
		return
	}

	fmt.Printf("\n🚀 Starting batch processing of %d todos...\n", len(req.Todos))
	fmt.Println("Watch the console to see workers processing concurrently!")
	fmt.Println("═══════════════════════════════════════════════════════")

	// Process batch concurrently using worker pool
	// GOROUTINES: Multiple workers process jobs simultaneously
	// CHANNELS: Jobs sent to workers, results collected
	response := h.batchProcessor.ProcessBatch(c.Request.Context(), req.Todos)

	fmt.Println("\n═══════════════════════════════════════════════════════")
	fmt.Printf("✅ Batch complete! Success: %d, Failed: %d, Time: %s\n\n",
		response.SuccessCount, response.FailureCount, response.TimeElapsed)

	utils.RespondOK(c, response)
}

// ProcessBatchV2 handles POST /api/v1/todos/batch-v2
//
// ADVANCED PATTERN:
// Uses semaphore pattern instead of worker pool
// One goroutine per item, but limited concurrency
func (h *BatchHandler) ProcessBatchV2(c *gin.Context) {
	var req dto.BatchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	fmt.Printf("\n🚀 Starting batch processing V2 (semaphore pattern) of %d todos...\n", len(req.Todos))

	// Use alternative processing method
	response := h.batchProcessor.ProcessBatchV2(c.Request.Context(), req.Todos)

	fmt.Printf("✅ Batch V2 complete! Success: %d, Failed: %d, Time: %s\n\n",
		response.SuccessCount, response.FailureCount, response.TimeElapsed)

	utils.RespondOK(c, response)
}
