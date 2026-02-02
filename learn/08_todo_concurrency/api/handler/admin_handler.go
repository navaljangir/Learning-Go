package handler

import (
	"todo_concurrency/domain/repository"
	"todo_concurrency/internal/dto"
	"todo_concurrency/internal/repository/cache"
	"todo_concurrency/internal/repository/memory"
	"todo_concurrency/internal/service"
	"todo_concurrency/pkg/utils"

	"github.com/gin-gonic/gin"
)

// AdminHandler handles administrative operations
//
// INTERFACE LEARNING:
// Demonstrates runtime switching of interface implementations
type AdminHandler struct {
	todoService  *service.TodoService
	statsService *service.StatsService
}

// NewAdminHandler creates a new admin handler
func NewAdminHandler(todoService *service.TodoService, statsService *service.StatsService) *AdminHandler {
	return &AdminHandler{
		todoService:  todoService,
		statsService: statsService,
	}
}

// SwitchStorage handles POST /api/v1/admin/switch-storage
//
// INTERFACE LEARNING:
// This demonstrates the power of interfaces!
// We can swap the entire storage backend at runtime.
//
// WARNING: This is for educational purposes only!
// In production, you wouldn't switch storage at runtime.
func (h *AdminHandler) SwitchStorage(c *gin.Context) {
	var req dto.SwitchStorageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	var newRepo repository.TodoRepository

	switch req.Backend {
	case "memory":
		newRepo = memory.NewInMemoryTodoRepository()
	case "cache":
		newRepo = cache.NewCachedTodoRepository(100)
	default:
		utils.RespondBadRequest(c, "unknown backend type")
		return
	}

	// This is a simplified example - in reality, you'd need to:
	// 1. Migrate existing data
	// 2. Handle the switch atomically
	// 3. Update all services that use the repository

	// For educational purposes, we just show that different implementations can be created
	storageType := "unknown"
	if si, ok := newRepo.(repository.StorageInfo); ok {
		storageType = si.GetStorageType()
	}

	utils.RespondSuccess(c, 200, gin.H{
		"new_backend":  req.Backend,
		"storage_type": storageType,
		"note":         "Storage backend created (educational example - not actually switched in running system)",
	}, "Storage backend demonstration")
}

// GetCurrentStorage handles GET /api/v1/admin/storage-info
//
// INTERFACE LEARNING:
// Uses type assertion to check if repository implements StorageInfo
func (h *AdminHandler) GetCurrentStorage(c *gin.Context) {
	repo := h.todoService.GetRepository()

	info := gin.H{
		"implements_storage_info": false,
	}

	// Type assertion - check if repo implements StorageInfo interface
	// INTERFACE LEARNING:
	// This is how you check for optional interface methods in Go
	if si, ok := repo.(repository.StorageInfo); ok {
		info["implements_storage_info"] = true
		info["storage_type"] = si.GetStorageType()
		info["stats"] = si.GetStats()
	}

	utils.RespondOK(c, info)
}
