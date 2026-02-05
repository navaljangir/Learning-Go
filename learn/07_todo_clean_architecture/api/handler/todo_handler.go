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
// @Summary Create a new todo
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTodoRequest true "Todo details"
// @Success 201 {object} dto.TodoResponse
// @Failure 400 {object} utils.Response
// @Router /api/v1/todos [post]
func (h *TodoHandler) Create(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	var req dto.CreateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	response, err := h.todoService.Create(c.Request.Context(), userID, req)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Created(c, response)
}

// List handles listing todos with pagination
// @Summary List user's todos
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(10)
// @Success 200 {object} dto.TodoListResponse
// @Failure 500 {object} utils.Response
// @Router /api/v1/todos [get]
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
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, response)
}

// GetByID handles getting a specific todo
// @Summary Get a todo by ID
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 200 {object} dto.TodoResponse
// @Failure 404 {object} utils.Response
// @Router /api/v1/todos/{id} [get]
func (h *TodoHandler) GetByID(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid todo ID")
		return
	}

	response, err := h.todoService.GetByID(c.Request.Context(), todoID, userID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.Success(c, response)
}

// Update handles updating a todo
// @Summary Update a todo
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Param request body dto.UpdateTodoRequest true "Todo update details"
// @Success 200 {object} dto.TodoResponse
// @Failure 400 {object} utils.Response
// @Router /api/v1/todos/{id} [put]
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
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, response)
}

// ToggleComplete handles toggling the completion status of a todo
// @Summary Toggle todo completion status
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 200 {object} dto.TodoResponse
// @Failure 400 {object} utils.Response
// @Router /api/v1/todos/{id}/toggle [patch]
func (h *TodoHandler) ToggleComplete(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid todo ID")
		return
	}

	response, err := h.todoService.ToggleComplete(c.Request.Context(), todoID, userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, response)
}

// Delete handles deleting a todo
// @Summary Delete a todo
// @Tags todos
// @Produce json
// @Security BearerAuth
// @Param id path string true "Todo ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/todos/{id} [delete]
func (h *TodoHandler) Delete(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	todoID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid todo ID")
		return
	}

	if err := h.todoService.Delete(c.Request.Context(), todoID, userID); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "todo deleted successfully"})
}

// MoveTodos handles moving multiple todos to a list or to global
// @Summary Move todos to a list or to global
// @Description Move multiple todos to a specific list (list_id) or to global (list_id = null)
// @Tags todos
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.MoveTodosRequest true "Move todos request"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/todos/move [patch]
func (h *TodoHandler) MoveTodos(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	var req dto.MoveTodosRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if err := h.todoService.MoveTodos(c.Request.Context(), userID, req); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.Success(c, gin.H{"message": "todos moved successfully"})
}
