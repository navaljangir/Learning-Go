package postgres

import (
	"context"
	"database/sql"
	"errors"

	"demo/domain/repository"
)

// PostgresTodoRepository stores todos in PostgreSQL
type PostgresTodoRepository struct {
	db *sql.DB
}

// NewPostgresTodRepository creates a new postgres repository
func NewPostgresTodRepository(db *sql.DB) repository.TodoRepository {
	return &PostgresTodoRepository{db: db}
}

// ============================================================================
// REQUIRED INTERFACE IMPLEMENTATION - TodoRepository
// ============================================================================

func (r *PostgresTodoRepository) Create(ctx context.Context, todo *repository.Todo) error {
	query := `INSERT INTO todos (id, title, completed, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, todo.ID, todo.Title, todo.Completed, todo.CreatedAt)
	return err
}

func (r *PostgresTodoRepository) FindByID(ctx context.Context, id string) (*repository.Todo, error) {
	query := `SELECT id, title, completed, created_at FROM todos WHERE id = $1`

	var todo repository.Todo
	err := r.db.QueryRowContext(ctx, query, id).Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, errors.New("not found")
	}
	if err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *PostgresTodoRepository) FindAll(ctx context.Context) ([]*repository.Todo, error) {
	query := `SELECT id, title, completed, created_at FROM todos`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*repository.Todo
	for rows.Next() {
		var todo repository.Todo
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Completed, &todo.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	return todos, nil
}

func (r *PostgresTodoRepository) Update(ctx context.Context, todo *repository.Todo) error {
	query := `UPDATE todos SET title = $1, completed = $2 WHERE id = $3`
	result, err := r.db.ExecContext(ctx, query, todo.Title, todo.Completed, todo.ID)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("not found")
	}

	return nil
}

func (r *PostgresTodoRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM todos WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("not found")
	}

	return nil
}

// ============================================================================
// OPTIONAL INTERFACE IMPLEMENTATION - StorageInfo
// ============================================================================

// GetStorageType implements StorageInfo interface
func (r *PostgresTodoRepository) GetStorageType() string {
	return "postgresql"
}

// GetStats implements StorageInfo interface
func (r *PostgresTodoRepository) GetStats() map[string]interface{} {
	var count int
	r.db.QueryRow("SELECT COUNT(*) FROM todos").Scan(&count)

	var dbSize string
	r.db.QueryRow("SELECT pg_size_pretty(pg_database_size(current_database()))").Scan(&dbSize)

	return map[string]interface{}{
		"storage_type":  "postgresql",
		"total_todos":   count,
		"database_size": dbSize,
	}
}

// ============================================================================
// OPTIONAL INTERFACE IMPLEMENTATION - BatchCapable
// ============================================================================

// BatchCreate implements BatchCapable interface
func (r *PostgresTodoRepository) BatchCreate(ctx context.Context, todos []*repository.Todo) error {
	// Start transaction
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Prepare statement
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO todos (id, title, completed, created_at) VALUES ($1, $2, $3, $4)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute batch insert
	for _, todo := range todos {
		if _, err := stmt.ExecContext(ctx, todo.ID, todo.Title, todo.Completed, todo.CreatedAt); err != nil {
			return err
		}
	}

	// Commit transaction
	return tx.Commit()
}

// BatchDelete implements BatchCapable interface
func (r *PostgresTodoRepository) BatchDelete(ctx context.Context, ids []string) error {
	query := `DELETE FROM todos WHERE id = ANY($1)`
	_, err := r.db.ExecContext(ctx, query, ids)
	return err
}

// ============================================================================
// COMPILE-TIME CHECKS
// ============================================================================

// Verify this type implements the required interface
var _ repository.TodoRepository = (*PostgresTodoRepository)(nil)

// Verify this type implements OPTIONAL interfaces
var _ repository.StorageInfo = (*PostgresTodoRepository)(nil)
var _ repository.BatchCapable = (*PostgresTodoRepository)(nil)

// NOTE: We DON'T implement CacheCapable
// So we DON'T have a compile-time check for that
