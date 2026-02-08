package handler

import (
	"todo_app/domain/service"
	"todo_app/internal/dto"
	"todo_app/pkg/constants"
	"todo_app/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TodoListHandler handles todo list-related HTTP requests
type TodoListHandler struct {
	listService service.TodoListService
}

// NewTodoListHandler creates a new todo list handler
// Returns TodoListHandlerInterface to enforce dependency on abstraction
func NewTodoListHandler(listService service.TodoListService) TodoListHandlerInterface {
	return &TodoListHandler{listService: listService}
}

// Create handles creating a new list
// @Summary Create a new list
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateListRequest true "List details"
// @Success 201 {object} dto.ListResponse
// @Failure 400 {object} utils.Response
// @Router /api/v1/lists [post]
func (h *TodoListHandler) Create(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	var req dto.CreateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	response, err := h.listService.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.Created(c, response)
}

// List handles listing all lists for a user
// @Summary List user's lists
// @Tags lists
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.ListsResponse
// @Failure 500 {object} utils.Response
// @Router /api/v1/lists [get]
func (h *TodoListHandler) List(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	response, err := h.listService.List(c.Request.Context(), userID)
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, response)
}

// GetByID handles getting a specific list with its todos
// @Summary Get a list by ID with todos
// @Tags lists
// @Produce json
// @Security BearerAuth
// @Param id path string true "List ID"
// @Success 200 {object} dto.ListWithTodosResponse
// @Failure 404 {object} utils.Response
// @Router /api/v1/lists/{id} [get]
func (h *TodoListHandler) GetByID(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	listID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid list ID")
		return
	}

	response, err := h.listService.GetByID(c.Request.Context(), listID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, response)
}

// Update handles updating a list (rename)
// @Summary Update a list
// @Tags lists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "List ID"
// @Param request body dto.UpdateListRequest true "List update details"
// @Success 200 {object} dto.ListResponse
// @Failure 400 {object} utils.Response
// @Router /api/v1/lists/{id} [put]
func (h *TodoListHandler) Update(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	listID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid list ID")
		return
	}

	var req dto.UpdateListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	response, err := h.listService.Update(c.Request.Context(), listID, userID, req)
	if err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, response)
}

// Delete handles deleting a list
// @Summary Delete a list
// @Description Permanently deletes a list and all todos within it
// @Tags lists
// @Produce json
// @Security BearerAuth
// @Param id path string true "List ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Router /api/v1/lists/{id} [delete]
func (h *TodoListHandler) Delete(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	listID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid list ID")
		return
	}

	if err := h.listService.Delete(c.Request.Context(), listID, userID); err != nil {
		c.Error(err)
		return
	}

	utils.Success(c, gin.H{"message": "list deleted successfully"})
}

// Duplicate handles duplicating a list with all its todos
// @Summary Duplicate a list
// @Description Creates a copy of a list with all its todos
// @Tags lists
// @Produce json
// @Security BearerAuth
// @Param id path string true "List ID"
// @Success 201 {object} dto.ListWithTodosResponse
// @Failure 400 {object} utils.Response
// @Router /api/v1/lists/{id}/duplicate [post]
func (h *TodoListHandler) Duplicate(c *gin.Context) {
	userID := c.MustGet(constants.ContextUserID).(uuid.UUID)

	listID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "invalid list ID")
		return
	}

	response, err := h.listService.Duplicate(c.Request.Context(), listID, userID)
	if err != nil {
		c.Error(err)
		return
	}

	utils.Created(c, response)
}
