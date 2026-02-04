package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"todo_concurrency/domain/entity"
	"todo_concurrency/domain/repository"
)

// InMemoryTodoRepository stores todos in memory using a map
//
// KEY LEARNING - MUTEX:
// Since multiple goroutines (HTTP handlers) can access this map concurrently,
// we need a mutex (mutual exclusion lock) to prevent race conditions.
//
// Without mutex: Two goroutines could read/write simultaneously, causing data corruption
// With mutex: Only one goroutine can access the map at a time
type InMemoryTodoRepository struct {
	mu     sync.RWMutex            // RWMutex allows multiple readers OR one writer
	todos  map[string]*entity.Todo // The actual storage
	nextID int                     // Auto-incrementing ID

	// Statistics (protected by the same mutex)
	accessCount int
	lastAccess  time.Time
}

// NewInMemoryTodoRepository creates a new in-memory repository
func NewInMemoryTodoRepository() repository.TodoRepository {
	return &InMemoryTodoRepository{
		todos:  make(map[string]*entity.Todo),
		nextID: 1,
	}
}

// Create adds a new todo
// MUTEX LEARNING: We use Lock() for write operations
func (r *InMemoryTodoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	r.mu.Lock()         // LOCK - No other goroutine can access now
	defer r.mu.Unlock() // UNLOCK - Always unlock when function returns

	// Generate ID
	todo.ID = fmt.Sprintf("%d", r.nextID)
	r.nextID++

	// Set timestamps
	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	// Store
	r.todos[todo.ID] = todo

	// Update stats
	r.accessCount++
	r.lastAccess = time.Now()

	return nil
}

// FindByID retrieves a todo by ID
// MUTEX LEARNING: We use RLock() for read operations
// RLock allows multiple readers simultaneously (better performance)
func (r *InMemoryTodoRepository) FindByID(ctx context.Context, id string) (*entity.Todo, error) {
	r.mu.RLock()         // READ LOCK - Multiple goroutines can hold this
	defer r.mu.RUnlock() // UNLOCK

	r.accessCount++
	r.lastAccess = time.Now()

	todo, exists := r.todos[id]
	if !exists {
		return nil, nil // Not found (not an error)
	}

	// Return a copy to prevent external modification
	todoCopy := *todo
	return &todoCopy, nil
}

// FindAll retrieves all todos
func (r *InMemoryTodoRepository) FindAll(ctx context.Context) ([]*entity.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	r.accessCount++
	r.lastAccess = time.Now()

	// Convert map to slice
	todos := make([]*entity.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todoCopy := *todo
		todos = append(todos, &todoCopy)
	}

	return todos, nil
}

// Update modifies an existing todo
func (r *InMemoryTodoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[todo.ID]; !exists {
		return errors.New("todo not found")
	}

	todo.UpdatedAt = time.Now()
	r.todos[todo.ID] = todo

	r.accessCount++
	r.lastAccess = time.Now()

	return nil
}

// Delete removes a todo
func (r *InMemoryTodoRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[id]; !exists {
		return errors.New("todo not found")
	}

	delete(r.todos, id)

	r.accessCount++
	r.lastAccess = time.Now()

	return nil
}

// Count returns total todos
func (r *InMemoryTodoRepository) Count(ctx context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.todos), nil
}

// CountCompleted returns completed todos count
func (r *InMemoryTodoRepository) CountCompleted(ctx context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, todo := range r.todos {
		if todo.Completed {
			count++
		}
	}

	return count, nil
}

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

// Verify interface implementation at compile time
var _ repository.TodoRepository = (*InMemoryTodoRepository)(nil)
var _ repository.StorageInfo = (*InMemoryTodoRepository)(nil)
