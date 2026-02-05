package sqlc_impl

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"todo_app/domain/entity"
	"todo_app/domain/repository"
	"todo_app/internal/repository/sqlc"

	"github.com/google/uuid"
)

type todoRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewTodoRepository creates a new sqlc-based todo repository
func NewTodoRepository(db *sql.DB) repository.TodoRepository {
	return &todoRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

// convertTodoRowToEntity is a helper to convert any todo row type to entity
func convertTodoRowToEntity(id, userID string, listID sql.NullString, title string, description sql.NullString,
	completed bool, priority sql.NullString, dueDate sql.NullTime, createdAt, updatedAt time.Time,
	completedAt, deletedAt sql.NullTime) *entity.Todo {

	todoID, _ := uuid.Parse(id)
	todoUserID, _ := uuid.Parse(userID)

	var todoListID *uuid.UUID
	if listID.Valid {
		parsed, _ := uuid.Parse(listID.String)
		todoListID = &parsed
	}

	desc := ""
	if description.Valid {
		desc = description.String
	}

	prio := entity.PriorityMedium
	if priority.Valid {
		prio = entity.Priority(priority.String)
	}

	return &entity.Todo{
		ID:          todoID,
		UserID:      todoUserID,
		ListID:      todoListID,
		Title:       title,
		Description: desc,
		Completed:   completed,
		Priority:    prio,
		DueDate:     fromNullTimePtr(dueDate),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		CompletedAt: fromNullTimePtr(completedAt),
		DeletedAt:   fromNullTimePtr(deletedAt),
	}
}

// todoRowByIDToEntity converts sqlc.GetTodoByIDRow to entity.Todo
func todoRowByIDToEntity(t sqlc.GetTodoByIDRow) *entity.Todo {
	return convertTodoRowToEntity(
		t.ID, t.UserID, t.ListID, t.Title, t.Description,
		t.Completed, t.Priority, t.DueDate, t.CreatedAt, t.UpdatedAt,
		t.CompletedAt, t.DeletedAt,
	)
}

// todoRowByUserIDToEntity converts sqlc.GetTodosByUserIDRow to entity.Todo
func todoRowByUserIDToEntity(t sqlc.GetTodosByUserIDRow) *entity.Todo {
	return convertTodoRowToEntity(
		t.ID, t.UserID, t.ListID, t.Title, t.Description,
		t.Completed, t.Priority, t.DueDate, t.CreatedAt, t.UpdatedAt,
		t.CompletedAt, t.DeletedAt,
	)
}

// todoRowFilteredToEntity converts sqlc.GetTodosFilteredRow to entity.Todo
func todoRowFilteredToEntity(t sqlc.GetTodosFilteredRow) *entity.Todo {
	return convertTodoRowToEntity(
		t.ID, t.UserID, t.ListID, t.Title, t.Description,
		t.Completed, t.Priority, t.DueDate, t.CreatedAt, t.UpdatedAt,
		t.CompletedAt, t.DeletedAt,
	)
}

func (r *todoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	var listID sql.NullString
	if todo.ListID != nil {
		listID = sql.NullString{String: todo.ListID.String(), Valid: true}
	}

	params := sqlc.CreateTodoParams{
		ID:          todo.ID.String(),
		UserID:      todo.UserID.String(),
		ListID:      listID,
		Title:       todo.Title,
		Description: toNullString(todo.Description),
		Completed:   todo.Completed,
		Priority:    toNullString(string(todo.Priority)),
		DueDate:     toNullTimePtr(todo.DueDate),
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.UpdatedAt,
	}
	return r.queries.CreateTodo(ctx, params)
}

func (r *todoRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Todo, error) {
	todo, err := r.queries.GetTodoByID(ctx, id.String())
	if err == sql.ErrNoRows {
		return nil, errors.New("todo not found")
	}
	if err != nil {
		return nil, err
	}
	return todoRowByIDToEntity(todo), nil
}

func (r *todoRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Todo, error) {
	params := sqlc.GetTodosByUserIDParams{
		UserID: userID.String(),
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	todos, err := r.queries.GetTodosByUserID(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Todo, len(todos))
	for i, t := range todos {
		result[i] = todoRowByUserIDToEntity(t)
	}

	return result, nil
}

func (r *todoRepository) FindWithFilters(ctx context.Context, filter repository.TodoFilter, limit, offset int) ([]*entity.Todo, error) {
	params := sqlc.GetTodosFilteredParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	// Convert pointer types to nullable types
	if filter.UserID != nil {
		params.UserID = sql.NullString{String: filter.UserID.String(), Valid: true}
	}

	if filter.Completed != nil {
		params.Completed = sql.NullBool{Bool: *filter.Completed, Valid: true}
	}

	if filter.Priority != nil {
		params.Priority = sql.NullString{String: string(*filter.Priority), Valid: true}
	}

	if filter.FromDate != nil {
		params.DueDateFrom = sql.NullTime{Time: *filter.FromDate, Valid: true}
	}

	if filter.ToDate != nil {
		params.DueDateTo = sql.NullTime{Time: *filter.ToDate, Valid: true}
	}

	todos, err := r.queries.GetTodosFiltered(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Todo, len(todos))
	for i, t := range todos {
		result[i] = todoRowFilteredToEntity(t)
	}

	return result, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	var listID sql.NullString
	if todo.ListID != nil {
		listID = sql.NullString{String: todo.ListID.String(), Valid: true}
	}

	params := sqlc.UpdateTodoParams{
		ID:          todo.ID.String(),
		Title:       todo.Title,
		Description: toNullString(todo.Description),
		Completed:   todo.Completed,
		Priority:    toNullString(string(todo.Priority)),
		DueDate:     toNullTimePtr(todo.DueDate),
		UpdatedAt:   todo.UpdatedAt,
		CompletedAt: toNullTimePtr(todo.CompletedAt),
		ListID:      listID,
	}
	return r.queries.UpdateTodo(ctx, params)
}

func (r *todoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.SoftDeleteTodo(ctx, id.String())
}

func (r *todoRepository) Count(ctx context.Context, filter repository.TodoFilter) (int64, error) {
	params := sqlc.CountTodosFilteredParams{}

	// Convert pointer types to nullable types
	if filter.UserID != nil {
		params.UserID = sql.NullString{String: filter.UserID.String(), Valid: true}
	}

	if filter.Completed != nil {
		params.Completed = sql.NullBool{Bool: *filter.Completed, Valid: true}
	}

	if filter.Priority != nil {
		params.Priority = sql.NullString{String: string(*filter.Priority), Valid: true}
	}

	if filter.FromDate != nil {
		params.DueDateFrom = sql.NullTime{Time: *filter.FromDate, Valid: true}
	}

	if filter.ToDate != nil {
		params.DueDateTo = sql.NullTime{Time: *filter.ToDate, Valid: true}
	}

	return r.queries.CountTodosFiltered(ctx, params)
}

func (r *todoRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.CountTodosByUser(ctx, userID.String())
}

func (r *todoRepository) UpdateListID(ctx context.Context, todoIDs []uuid.UUID, listID *uuid.UUID, userID uuid.UUID) error {
	if len(todoIDs) == 0 {
		return nil
	}

	// Build placeholders for IN clause
	placeholders := make([]interface{}, 0, len(todoIDs)+2)

	// First parameter is list_id (can be NULL)
	if listID != nil {
		placeholders = append(placeholders, listID.String())
	} else {
		placeholders = append(placeholders, nil)
	}

	// Add todo IDs
	query := "UPDATE todos SET list_id = ?, updated_at = CURRENT_TIMESTAMP WHERE id IN ("
	for i, id := range todoIDs {
		if i > 0 {
			query += ","
		}
		query += "?"
		placeholders = append(placeholders, id.String())
	}
	query += ") AND user_id = ? AND deleted_at IS NULL"

	// Last parameter is user_id
	placeholders = append(placeholders, userID.String())

	_, err := r.db.ExecContext(ctx, query, placeholders...)
	return err
}
