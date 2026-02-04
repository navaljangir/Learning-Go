package file

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"demo/domain/repository"
)

// FileTodoRepository stores todos in a JSON file
type FileTodoRepository struct {
	mu       sync.RWMutex
	filePath string
	todos    map[string]*repository.Todo
}

// NewFileTodoRepository creates a new file-based repository
func NewFileTodoRepository(filePath string) repository.TodoRepository {
	repo := &FileTodoRepository{
		filePath: filePath,
		todos:    make(map[string]*repository.Todo),
	}

	// Load existing todos from file
	repo.load()

	return repo
}

// load reads todos from file
func (r *FileTodoRepository) load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, that's OK
		}
		return err
	}

	return json.Unmarshal(data, &r.todos)
}

// save writes todos to file
func (r *FileTodoRepository) save() error {
	data, err := json.MarshalIndent(r.todos, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}

// ============================================================================
// REQUIRED INTERFACE IMPLEMENTATION - TodoRepository
// ============================================================================

func (r *FileTodoRepository) Create(ctx context.Context, todo *repository.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.todos[todo.ID] = todo
	return r.save()
}

func (r *FileTodoRepository) FindByID(ctx context.Context, id string) (*repository.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todo, exists := r.todos[id]
	if !exists {
		return nil, errors.New("not found")
	}

	return todo, nil
}

func (r *FileTodoRepository) FindAll(ctx context.Context) ([]*repository.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	todos := make([]*repository.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		todos = append(todos, todo)
	}

	return todos, nil
}

func (r *FileTodoRepository) Update(ctx context.Context, todo *repository.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[todo.ID]; !exists {
		return errors.New("not found")
	}

	r.todos[todo.ID] = todo
	return r.save()
}

func (r *FileTodoRepository) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.todos[id]; !exists {
		return errors.New("not found")
	}

	delete(r.todos, id)
	return r.save()
}

// ============================================================================
// NO OPTIONAL INTERFACE IMPLEMENTATIONS
// ============================================================================

// We DON'T implement StorageInfo, BatchCapable, or CacheCapable
// because file-based storage doesn't have these capabilities

// ============================================================================
// COMPILE-TIME CHECKS
// ============================================================================

// Verify this type implements the required interface
var _ repository.TodoRepository = (*FileTodoRepository)(nil)

// NO checks for optional interfaces because we don't implement them!
// var _ repository.StorageInfo = (*FileTodoRepository)(nil)    // ← Would FAIL if uncommented
// var _ repository.BatchCapable = (*FileTodoRepository)(nil)   // ← Would FAIL if uncommented
