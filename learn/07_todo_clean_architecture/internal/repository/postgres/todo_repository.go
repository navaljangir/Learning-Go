package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"todo_app/domain/entity"
	"todo_app/domain/repository"

	"github.com/google/uuid"
)

type todoRepository struct {
	db *sql.DB
}

// NewTodoRepository creates a new PostgreSQL todo repository
func NewTodoRepository(db *sql.DB) repository.TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	query := `
		INSERT INTO todos (id, user_id, title, description, completed, priority, due_date, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		todo.ID, todo.UserID, todo.Title, todo.Description,
		todo.Completed, todo.Priority, todo.DueDate,
		todo.CreatedAt, todo.UpdatedAt,
	)
	return err
}

func (r *todoRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.Todo, error) {
	query := `
		SELECT id, user_id, title, description, completed, priority, due_date,
		       created_at, updated_at, completed_at, deleted_at
		FROM todos WHERE id = $1 AND deleted_at IS NULL
	`
	todo := &entity.Todo{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
		&todo.Completed, &todo.Priority, &todo.DueDate,
		&todo.CreatedAt, &todo.UpdatedAt, &todo.CompletedAt, &todo.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("todo not found")
	}
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (r *todoRepository) FindByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Todo, error) {
	query := `
		SELECT id, user_id, title, description, completed, priority, due_date,
		       created_at, updated_at, completed_at, deleted_at
		FROM todos
		WHERE user_id = $1 AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTodos(rows)
}

func (r *todoRepository) FindWithFilters(ctx context.Context, filter repository.TodoFilter, limit, offset int) ([]*entity.Todo, error) {
	query := `
		SELECT id, user_id, title, description, completed, priority, due_date,
		       created_at, updated_at, completed_at, deleted_at
		FROM todos
		WHERE deleted_at IS NULL
	`

	var conditions []string
	var args []interface{}
	argCount := 1

	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argCount))
		args = append(args, *filter.UserID)
		argCount++
	}

	if filter.Completed != nil {
		conditions = append(conditions, fmt.Sprintf("completed = $%d", argCount))
		args = append(args, *filter.Completed)
		argCount++
	}

	if filter.Priority != nil {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argCount))
		args = append(args, *filter.Priority)
		argCount++
	}

	if filter.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("due_date >= $%d", argCount))
		args = append(args, *filter.FromDate)
		argCount++
	}

	if filter.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("due_date <= $%d", argCount))
		args = append(args, *filter.ToDate)
		argCount++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanTodos(rows)
}

func (r *todoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	query := `
		UPDATE todos
		SET title = $1, description = $2, completed = $3, priority = $4,
		    due_date = $5, updated_at = $6, completed_at = $7
		WHERE id = $8 AND deleted_at IS NULL
	`
	result, err := r.db.ExecContext(ctx, query,
		todo.Title, todo.Description, todo.Completed, todo.Priority,
		todo.DueDate, todo.UpdatedAt, todo.CompletedAt, todo.ID,
	)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("todo not found")
	}

	return nil
}

func (r *todoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE todos SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("todo not found")
	}

	return nil
}

func (r *todoRepository) Count(ctx context.Context, filter repository.TodoFilter) (int64, error) {
	query := `SELECT COUNT(*) FROM todos WHERE deleted_at IS NULL`

	var conditions []string
	var args []interface{}
	argCount := 1

	if filter.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argCount))
		args = append(args, *filter.UserID)
		argCount++
	}

	if filter.Completed != nil {
		conditions = append(conditions, fmt.Sprintf("completed = $%d", argCount))
		args = append(args, *filter.Completed)
		argCount++
	}

	if filter.Priority != nil {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", argCount))
		args = append(args, *filter.Priority)
		argCount++
	}

	if filter.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("due_date >= $%d", argCount))
		args = append(args, *filter.FromDate)
		argCount++
	}

	if filter.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("due_date <= $%d", argCount))
		args = append(args, *filter.ToDate)
		argCount++
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	var count int64
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *todoRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	query := `SELECT COUNT(*) FROM todos WHERE user_id = $1 AND deleted_at IS NULL`
	var count int64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	return count, err
}

// scanTodos is a helper function to scan multiple todos from rows
func (r *todoRepository) scanTodos(rows *sql.Rows) ([]*entity.Todo, error) {
	var todos []*entity.Todo
	for rows.Next() {
		todo := &entity.Todo{}
		err := rows.Scan(
			&todo.ID, &todo.UserID, &todo.Title, &todo.Description,
			&todo.Completed, &todo.Priority, &todo.DueDate,
			&todo.CreatedAt, &todo.UpdatedAt, &todo.CompletedAt, &todo.DeletedAt,
		)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}
