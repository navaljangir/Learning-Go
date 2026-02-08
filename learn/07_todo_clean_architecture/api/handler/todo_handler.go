package handler

import (
	"strconv"
	"todo_app/domain/service"
	"todo_app/internal/dto"
	"todo_app/pkg/constants"
	"todo_app/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TodoHandler handles todo-related HTTP requests
type TodoHandler struct {
	todoService service.TodoService
}

// NewTodoHandler creates a new todo handler
// Returns TodoHandlerInterface to enforce dependency on abstraction
func NewTodoHandler(todoService service.TodoService) TodoHandlerInterface {
	return &TodoHandler{todoService: todoService}
}

// Create handles creating a new todo
func (h *TodoHandler) Create(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	var req dto.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	response, err := h.todoService.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.Created(c, response)
}

// List handles listing todos with pagination
func (h *TodoHandler) List(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > constants.MaxPageSize {
		pageSize = constants.DefaultPageSize
	}

	response, err := h.todoService.List(c.Request.Context(), userID, page, pageSize)
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, response)
}

// GetByID handles getting a specific todo
func (h *TodoHandler) GetByID(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid todo ID")
		return
	}

	response, err := h.todoService.GetByID(c.Request.Context(), todoID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, response)
}

// Update handles updating a todo
func (h *TodoHandler) Update(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid todo ID")
		return
	}

	var req dto.UpdateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	response, err := h.todoService.Update(c.Request.Context(), todoID, userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, response)
}

// ToggleComplete handles toggling the completion status of a todo
func (h *TodoHandler) ToggleComplete(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid todo ID")
		return
	}

	response, err := h.todoService.ToggleComplete(c.Request.Context(), todoID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, response)
}

// Delete handles deleting a todo
func (h *TodoHandler) Delete(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid todo ID")
		return
	}

	if err := h.todoService.Delete(c.Request.Context(), todoID, userID); err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, gin.H{"message": "todo deleted successfully"})
}

// MoveTodos handles moving multiple todos to a list or to global
func (h *TodoHandler) MoveTodos(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	var req dto.MoveTodosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	if err := h.todoService.MoveTodos(c.Request.Context(), userID, req); err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, gin.H{"message": "todos moved successfully"})
}
