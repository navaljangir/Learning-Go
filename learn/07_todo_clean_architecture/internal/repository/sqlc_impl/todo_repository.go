package sqlc_impl

import (
	"context"
	"database/sql"
	"errors"
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

// sqlcTodoToEntity converts sqlc.Todo to entity.Todo
func sqlcTodoToEntity(t sqlc.Todo) *entity.Todo {
	description := ""
	if t.Description.Valid {
		description = t.Description.String
	}

	priority := entity.PriorityMedium
	if t.Priority.Valid {
		priority = entity.Priority(t.Priority.String)
	}

	return &entity.Todo{
		ID:          t.ID,
		UserID:      t.UserID,
		Title:       t.Title,
		Description: description,
		Completed:   t.Completed.Bool,
		Priority:    priority,
		DueDate:     fromNullTimePtr(t.DueDate),
		CreatedAt:   fromNullTime(t.CreatedAt),
		UpdatedAt:   fromNullTime(t.UpdatedAt),
		CompletedAt: fromNullTimePtr(t.CompletedAt),
		DeletedAt:   fromNullTimePtr(t.DeletedAt),
	}
}

func (r *todoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	params := sqlc.CreateTodoParams{
		ID:          todo.ID,
		UserID:      todo.UserID,
		Title:       todo.Title,
		Description: toNullString(todo.Description),
		Completed:   sql.NullBool{Bool: todo.Completed, Valid: true},
		Priority:    toNullString(string(todo.Priority)),
		DueDate:     toNullTimePtr(todo.DueDate),
		CreatedAt:   toNullTime(todo.CreatedAt),
		UpdatedAt:   toNullTime(todo.UpdatedAt),
	}
	return r.queries.CreateTodo(ctx, params)
}

func (r *todoRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Todo, error) {
	todo, err := r.queries.GetTodoByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, errors.New("todo not found")
	}
	if err != nil {
		return nil, err
	}
	return sqlcTodoToEntity(todo), nil
}

func (r *todoRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Todo, error) {
	params := sqlc.GetTodosByUserIDParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	todos, err := r.queries.GetTodosByUserID(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Todo, len(todos))
	for i, t := range todos {
		result[i] = sqlcTodoToEntity(t)
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
		params.UserID = uuid.NullUUID{UUID: *filter.UserID, Valid: true}
	}

	if filter.Completed != nil {
		params.Completed = sql.NullBool{Bool: *filter.Completed, Valid: true}
	}

	if filter.Priority != nil {
		params.Priority = sql.NullString{String: string(*filter.Priority), Valid: true}
	}

	if filter.FromDate != nil {
		params.FromDate = sql.NullTime{Time: *filter.FromDate, Valid: true}
	}

	if filter.ToDate != nil {
		params.ToDate = sql.NullTime{Time: *filter.ToDate, Valid: true}
	}

	todos, err := r.queries.GetTodosFiltered(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.Todo, len(todos))
	for i, t := range todos {
		result[i] = sqlcTodoToEntity(t)
	}

	return result, nil
}

func (r *todoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	params := sqlc.UpdateTodoParams{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: toNullString(todo.Description),
		Completed:   sql.NullBool{Bool: todo.Completed, Valid: true},
		Priority:    toNullString(string(todo.Priority)),
		DueDate:     toNullTimePtr(todo.DueDate),
		UpdatedAt:   toNullTime(todo.UpdatedAt),
		CompletedAt: toNullTimePtr(todo.CompletedAt),
	}
	return r.queries.UpdateTodo(ctx, params)
}

func (r *todoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.SoftDeleteTodo(ctx, id)
}

func (r *todoRepository) Count(ctx context.Context, filter repository.TodoFilter) (int64, error) {
	params := sqlc.CountTodosFilteredParams{}

	// Convert pointer types to nullable types
	if filter.UserID != nil {
		params.UserID = uuid.NullUUID{UUID: *filter.UserID, Valid: true}
	}

	if filter.Completed != nil {
		params.Completed = sql.NullBool{Bool: *filter.Completed, Valid: true}
	}

	if filter.Priority != nil {
		params.Priority = sql.NullString{String: string(*filter.Priority), Valid: true}
	}

	if filter.FromDate != nil {
		params.FromDate = sql.NullTime{Time: *filter.FromDate, Valid: true}
	}

	if filter.ToDate != nil {
		params.ToDate = sql.NullTime{Time: *filter.ToDate, Valid: true}
	}

	return r.queries.CountTodosFiltered(ctx, params)
}

func (r *todoRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.CountTodosByUser(ctx, userID)
}
