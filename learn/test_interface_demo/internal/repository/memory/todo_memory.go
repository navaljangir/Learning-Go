package memory

import (
	"context"
	"errors"
	"sync"
	"time"

	"demo/domain/repository"
)

// InMemoryTodoRepository stores todos in memory
type InMemoryTodoRepository struct {
	mu          sync.RWMutex
	todos       map[string]*repository.Todo
	accessCount int
	lastAccess  time.Time
}

// NewInMemoryTodoRepository creates a new in-memory repository
func NewInMemoryTodoRepository() repository.TodoRepository {
	return &InMemoryTodoRepository{
		todos: make(map[string]*repository.Todo),
	}
}

// ============================================================================
// REQUIRED INTERFACE IMPLEMENTATION - TodoRepository
// ============================================================================

func (r *InMemoryTodoRepository) Create(ctx context.Context, todo *repository.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	todo.CreatedAt = time.Now()
	r.todos[todo.ID] = todo

	r.accessCount++
	r.lastAccess = time.Now()

	return nil
}

func (r *InMemoryTodoRepository) FindByID(ctx context.Context, id string) (*repository.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	r.accessCount++
	r.lastAccess = time.Now()

	todo, exists := r.todos[id]
	if !exists {
		return nil, errors.New("not found")
	}

	return todo, nil
}

func (r *InMemoryTodoRepository) FindAll(ctx context.Context) ([]*repository.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	r.accessCount++
	r.lastAccess = time.Now()

	todos := make([]*repository.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *InMemoryTodoRepository) Update(ctx context.Context, todo *repository.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[todo.ID]; !exists {
		return errors.New("not found")
	}

	r.todos[todo.ID] = todo

	r.accessCount++
	r.lastAccess = time.Now()

	return nil
}

func (r *InMemoryTodoRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[id]; !exists {
		return errors.New("not found")
	}

	delete(r.todos, id)

	r.accessCount++
	r.lastAccess = time.Now()

	return nil
}

// ============================================================================
// OPTIONAL INTERFACE IMPLEMENTATION - StorageInfo
// ============================================================================

// GetStorageType implements StorageInfo interface
func (r *InMemoryTodoRepository) GetStorageType() string {
	return "in-memory"
}

// GetStats implements StorageInfo interface
func (r *InMemoryTodoRepository) GetStats() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return map[string]interface{}{
		"storage_type": "in-memory",
		"total_todos":  len(r.todos),
		"access_count": r.accessCount,
		"last_access":  r.lastAccess.Format(time.RFC3339),
	}
}

// ============================================================================
// COMPILE-TIME CHECKS
// ============================================================================

// Verify this type implements the required interface
var _ repository.TodoRepository = (*InMemoryTodoRepository)(nil)

// Verify this type implements the OPTIONAL StorageInfo interface
var _ repository.StorageInfo = (*InMemoryTodoRepository)(nil)

// NOTE: We DON'T implement BatchCapable or CacheCapable
// So we DON'T have compile-time checks for those
