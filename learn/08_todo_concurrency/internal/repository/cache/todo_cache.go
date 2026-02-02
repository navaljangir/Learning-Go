package cache

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
	"todo_concurrency/domain/entity"
	"todo_concurrency/domain/repository"
)

// CachedTodoRepository is an alternative implementation with caching statistics
//
// KEY LEARNING - INTERFACES:
// This is a DIFFERENT implementation of TodoRepository interface.
// The service layer doesn't care which implementation it uses!
//
// This also demonstrates more complex mutex usage with hit/miss tracking
type CachedTodoRepository struct {
	mu     sync.RWMutex
	todos  map[string]*entity.Todo
	nextID int

	// Cache statistics
	hits      int
	misses    int
	evictions int
	maxSize   int
}

// NewCachedTodoRepository creates a cached repository
func NewCachedTodoRepository(maxSize int) repository.TodoRepository {
	return &CachedTodoRepository{
		todos:   make(map[string]*entity.Todo),
		nextID:  1,
		maxSize: maxSize,
	}
}

func (r *CachedTodoRepository) Create(ctx context.Context, todo *entity.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if cache is full
	if len(r.todos) >= r.maxSize {
		// Evict oldest (simple strategy)
		r.evictOldest()
		r.evictions++
	}

	todo.ID = fmt.Sprintf("cache-%d", r.nextID)
	r.nextID++

	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	r.todos[todo.ID] = todo
	return nil
}

func (r *CachedTodoRepository) FindByID(ctx context.Context, id string) (*entity.Todo, error) {
	r.mu.Lock() // Note: Using Lock (not RLock) because we update stats
	defer r.mu.Unlock()

	todo, exists := r.todos[id]
	if !exists {
		r.misses++ // Track cache miss
		return nil, nil
	}

	r.hits++ // Track cache hit

	todoCopy := *todo
	return &todoCopy, nil
}

func (r *CachedTodoRepository) FindAll(ctx context.Context) ([]*entity.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos := make([]*entity.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todoCopy := *todo
		todos = append(todos, &todoCopy)
	}

	return todos, nil
}

func (r *CachedTodoRepository) Update(ctx context.Context, todo *entity.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[todo.ID]; !exists {
		r.misses++
		return errors.New("todo not found")
	}

	r.hits++
	todo.UpdatedAt = time.Now()
	r.todos[todo.ID] = todo

	return nil
}

func (r *CachedTodoRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[id]; !exists {
		r.misses++
		return errors.New("todo not found")
	}

	r.hits++
	delete(r.todos, id)

	return nil
}

func (r *CachedTodoRepository) Count(ctx context.Context) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.todos), nil
}

func (r *CachedTodoRepository) CountCompleted(ctx context.Context) (int, error) {
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

// evictOldest removes the oldest todo (simple LRU simulation)
// MUST be called with lock held
func (r *CachedTodoRepository) evictOldest() {
	var oldestID string
	var oldestTime time.Time

	first := true
	for id, todo := range r.todos {
		if first || todo.UpdatedAt.Before(oldestTime) {
			oldestID = id
			oldestTime = todo.UpdatedAt
			first = false
		}
	}

	if oldestID != "" {
		delete(r.todos, oldestID)
	}
}

// GetStorageType implements StorageInfo
func (r *CachedTodoRepository) GetStorageType() string {
	return "cached"
}

// GetStats implements StorageInfo with cache-specific metrics
func (r *CachedTodoRepository) GetStats() map[string]interface{} {
	r.mu.RLock()
	defer r.mu.RUnlock()

	hitRate := 0.0
	total := r.hits + r.misses
	if total > 0 {
		hitRate = float64(r.hits) / float64(total) * 100
	}

	return map[string]interface{}{
		"storage_type": "cached",
		"total_todos":  len(r.todos),
		"max_size":     r.maxSize,
		"hits":         r.hits,
		"misses":       r.misses,
		"evictions":    r.evictions,
		"hit_rate":     fmt.Sprintf("%.2f%%", hitRate),
	}
}

// Verify interface implementation
var _ repository.TodoRepository = (*CachedTodoRepository)(nil)
var _ repository.StorageInfo = (*CachedTodoRepository)(nil)
