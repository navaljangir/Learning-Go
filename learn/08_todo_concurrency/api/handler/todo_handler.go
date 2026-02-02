package handler

import (
	"time"
	"todo_concurrency/internal/dto"
	"todo_concurrency/internal/service"
	"todo_concurrency/pkg/utils"

	"github.com/gin-gonic/gin"
)

// TodoHandler handles basic CRUD operations for todos
//
// INTERFACE LEARNING:
// This handler works with any TodoService, which in turn works with any
// TodoRepository implementation. The handler doesn't care about storage details!
type TodoHandler struct {
	todoService  *service.TodoService
	statsService *service.StatsService
}

// NewTodoHandler creates a new todo handler
func NewTodoHandler(todoService *service.TodoService, statsService *service.StatsService) *TodoHandler {
	return &TodoHandler{
		todoService:  todoService,
		statsService: statsService,
	}
}

// Create handles POST /api/v1/todos
func (h *TodoHandler) Create(c *gin.Context) {
	startTime := time.Now()
	defer func() {
		h.statsService.RecordRequest("create", time.Since(startTime))
	}()

	var req dto.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	todo, err := h.todoService.Create(c.Request.Context(), req)
	if err != nil {
		utils.RespondInternalError(c, err.Error())
		return
	}

	utils.RespondCreated(c, todo)
}

// GetAll handles GET /api/v1/todos
func (h *TodoHandler) GetAll(c *gin.Context) {
	startTime := time.Now()
	defer func() {
		h.statsService.RecordRequest("read", time.Since(startTime))
	}()

	todos, err := h.todoService.GetAll(c.Request.Context())
	if err != nil {
		utils.RespondInternalError(c, err.Error())
		return
	}

	utils.RespondOK(c, todos)
}

// GetByID handles GET /api/v1/todos/:id
func (h *TodoHandler) GetByID(c *gin.Context) {
	startTime := time.Now()
	defer func() {
		h.statsService.RecordRequest("read", time.Since(startTime))
	}()

	id := c.Param("id")

	todo, err := h.todoService.GetByID(c.Request.Context(), id)
	if err != nil {
		utils.RespondNotFound(c, err.Error())
		return
	}

	utils.RespondOK(c, todo)
}

// Update handles PUT /api/v1/todos/:id
func (h *TodoHandler) Update(c *gin.Context) {
	startTime := time.Now()
	defer func() {
		h.statsService.RecordRequest("update", time.Since(startTime))
	}()

	id := c.Param("id")

	var req dto.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondBadRequest(c, err.Error())
		return
	}

	todo, err := h.todoService.Update(c.Request.Context(), id, req)
	if err != nil {
		utils.RespondNotFound(c, err.Error())
		return
	}

	utils.RespondOK(c, todo)
}

// Delete handles DELETE /api/v1/todos/:id
func (h *TodoHandler) Delete(c *gin.Context) {
	startTime := time.Now()
	defer func() {
		h.statsService.RecordRequest("delete", time.Since(startTime))
	}()

	id := c.Param("id")

	if err := h.todoService.Delete(c.Request.Context(), id); err != nil {
		utils.RespondNotFound(c, err.Error())
		return
	}

	utils.RespondSuccess(c, 200, nil, "Todo deleted successfully")
}

// ToggleComplete handles PATCH /api/v1/todos/:id/toggle
func (h *TodoHandler) ToggleComplete(c *gin.Context) {
	startTime := time.Now()
	defer func() {
		h.statsService.RecordRequest("update", time.Since(startTime))
	}()

	id := c.Param("id")

	todo, err := h.todoService.ToggleComplete(c.Request.Context(), id)
	if err != nil {
		utils.RespondNotFound(c, err.Error())
		return
	}

	utils.RespondOK(c, todo)
}
