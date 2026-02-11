package tests

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"
	"todo_app/domain/entity"
	"todo_app/domain/repository"

	"github.com/google/uuid"
)

// ---------------------------------------------------------------------------
// InMemoryUserRepo
// ---------------------------------------------------------------------------

// InMemoryUserRepo implements repository.UserRepository using an in-memory map.
// Thread-safe via sync.RWMutex.
type InMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[uuid.UUID]*entity.User
}

// Compile-time check
var _ repository.UserRepository = (*InMemoryUserRepo)(nil)

func NewInMemoryUserRepo() *InMemoryUserRepo {
	return &InMemoryUserRepo{users: make(map[uuid.UUID]*entity.User)}
}

func (r *InMemoryUserRepo) Create(_ context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok || u.IsDeleted() {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (r *InMemoryUserRepo) FindByUsername(_ context.Context, username string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Username == username && !u.IsDeleted() {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *InMemoryUserRepo) FindByEmail(_ context.Context, email string) (*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == email && !u.IsDeleted() {
			return u, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *InMemoryUserRepo) Update(_ context.Context, user *entity.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[user.ID]; !ok {
		return fmt.Errorf("user not found")
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepo) Delete(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[id]
	if !ok {
		return fmt.Errorf("user not found")
	}
	u.MarkDeleted()
	return nil
}

func (r *InMemoryUserRepo) List(_ context.Context, limit, offset int) ([]*entity.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var all []*entity.User
	for _, u := range r.users {
		if !u.IsDeleted() {
			all = append(all, u)
		}
	}
	// Sort by CreatedAt for deterministic ordering
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt.Before(all[j].CreatedAt)
	})

	if offset >= len(all) {
		return nil, nil
	}
	end := min(offset+limit, len(all))
	return all[offset:end], nil
}

func (r *InMemoryUserRepo) ExistsByUsername(_ context.Context, username string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Username == username && !u.IsDeleted() {
			return true, nil
		}
	}
	return false, nil
}

func (r *InMemoryUserRepo) ExistsByEmail(_ context.Context, email string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == email && !u.IsDeleted() {
			return true, nil
		}
	}
	return false, nil
}

// ---------------------------------------------------------------------------
// InMemoryTodoRepo
// ---------------------------------------------------------------------------

// InMemoryTodoRepo implements repository.TodoRepository using an in-memory map.
type InMemoryTodoRepo struct {
	mu    sync.RWMutex
	todos map[uuid.UUID]*entity.Todo
}

var _ repository.TodoRepository = (*InMemoryTodoRepo)(nil)

func NewInMemoryTodoRepo() *InMemoryTodoRepo {
	return &InMemoryTodoRepo{todos: make(map[uuid.UUID]*entity.Todo)}
}

func (r *InMemoryTodoRepo) Create(_ context.Context, todo *entity.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.todos[todo.ID] = todo
	return nil
}

func (r *InMemoryTodoRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.todos[id]
	if !ok || t.IsDeleted() {
		return nil, fmt.Errorf("todo not found")
	}
	return t, nil
}

func (r *InMemoryTodoRepo) FindByUserID(_ context.Context, userID uuid.UUID, limit, offset int) ([]*entity.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Todo
	for _, t := range r.todos {
		if t.UserID == userID && !t.IsDeleted() {
			result = append(result, t)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	if offset >= len(result) {
		return nil, nil
	}
	end := min(offset+limit, len(result))
	return result[offset:end], nil
}

func (r *InMemoryTodoRepo) FindByListID(_ context.Context, listID uuid.UUID) ([]*entity.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Todo
	for _, t := range r.todos {
		if t.ListID != nil && *t.ListID == listID && !t.IsDeleted() {
			result = append(result, t)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result, nil
}

func (r *InMemoryTodoRepo) FindWithFilters(_ context.Context, filter repository.TodoFilter, limit, offset int) ([]*entity.Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.Todo
	for _, t := range r.todos {
		if t.IsDeleted() {
			continue
		}
		if filter.UserID != nil && t.UserID != *filter.UserID {
			continue
		}
		if filter.Completed != nil && t.Completed != *filter.Completed {
			continue
		}
		if filter.Priority != nil && t.Priority != *filter.Priority {
			continue
		}
		if filter.FromDate != nil && t.CreatedAt.Before(*filter.FromDate) {
			continue
		}
		if filter.ToDate != nil && t.CreatedAt.After(*filter.ToDate) {
			continue
		}
		result = append(result, t)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

	if offset >= len(result) {
		return nil, nil
	}
	end := min(offset+limit, len(result))
	return result[offset:end], nil
}

func (r *InMemoryTodoRepo) Update(_ context.Context, todo *entity.Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.todos[todo.ID]; !ok {
		return fmt.Errorf("todo not found")
	}
	r.todos[todo.ID] = todo
	return nil
}

func (r *InMemoryTodoRepo) Delete(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	t, ok := r.todos[id]
	if !ok {
		return fmt.Errorf("todo not found")
	}
	t.MarkDeleted()
	return nil
}

func (r *InMemoryTodoRepo) Count(_ context.Context, filter repository.TodoFilter) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var count int64
	for _, t := range r.todos {
		if t.IsDeleted() {
			continue
		}
		if filter.UserID != nil && t.UserID != *filter.UserID {
			continue
		}
		if filter.Completed != nil && t.Completed != *filter.Completed {
			continue
		}
		if filter.Priority != nil && t.Priority != *filter.Priority {
			continue
		}
		count++
	}
	return count, nil
}

func (r *InMemoryTodoRepo) CountByUser(_ context.Context, userID uuid.UUID) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var count int64
	for _, t := range r.todos {
		if t.UserID == userID && !t.IsDeleted() {
			count++
		}
	}
	return count, nil
}

func (r *InMemoryTodoRepo) UpdateListID(_ context.Context, todoIDs []uuid.UUID, listID *uuid.UUID, userID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, id := range todoIDs {
		t, ok := r.todos[id]
		if !ok || t.IsDeleted() {
			return fmt.Errorf("todo not found")
		}
		if t.UserID != userID {
			return fmt.Errorf("forbidden")
		}
		t.ListID = listID
		t.UpdatedAt = time.Now()
	}
	return nil
}

// DeleteByListID removes all todos belonging to a list (simulates CASCADE).
// Called by InMemoryTodoListRepo.Delete.
func (r *InMemoryTodoRepo) DeleteByListID(listID uuid.UUID) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, t := range r.todos {
		if t.ListID != nil && *t.ListID == listID {
			t.MarkDeleted()
		}
	}
}

// ---------------------------------------------------------------------------
// InMemoryTodoListRepo
// ---------------------------------------------------------------------------

// InMemoryTodoListRepo implements repository.TodoListRepository using an in-memory map.
type InMemoryTodoListRepo struct {
	mu       sync.RWMutex
	lists    map[uuid.UUID]*entity.TodoList
	todoRepo *InMemoryTodoRepo // reference for CASCADE delete
}

var _ repository.TodoListRepository = (*InMemoryTodoListRepo)(nil)

func NewInMemoryTodoListRepo(todoRepo *InMemoryTodoRepo) *InMemoryTodoListRepo {
	return &InMemoryTodoListRepo{
		lists:    make(map[uuid.UUID]*entity.TodoList),
		todoRepo: todoRepo,
	}
}

func (r *InMemoryTodoListRepo) Create(_ context.Context, list *entity.TodoList) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lists[list.ID] = list
	return nil
}

func (r *InMemoryTodoListRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.TodoList, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	l, ok := r.lists[id]
	if !ok || l.IsDeleted() {
		return nil, fmt.Errorf("list not found")
	}
	return l, nil
}

func (r *InMemoryTodoListRepo) FindByUserID(_ context.Context, userID uuid.UUID) ([]*entity.TodoList, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*entity.TodoList
	for _, l := range r.lists {
		if l.UserID == userID && !l.IsDeleted() {
			result = append(result, l)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})
	return result, nil
}

func (r *InMemoryTodoListRepo) Update(_ context.Context, list *entity.TodoList) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.lists[list.ID]; !ok {
		return fmt.Errorf("list not found")
	}
	r.lists[list.ID] = list
	return nil
}

func (r *InMemoryTodoListRepo) Delete(_ context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	l, ok := r.lists[id]
	if !ok {
		return fmt.Errorf("list not found")
	}
	l.MarkDeleted()
	// CASCADE: delete all todos in this list
	r.todoRepo.DeleteByListID(id)
	return nil
}

func (r *InMemoryTodoListRepo) CountByUser(_ context.Context, userID uuid.UUID) (int64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var count int64
	for _, l := range r.lists {
		if l.UserID == userID && !l.IsDeleted() {
			count++
		}
	}
	return count, nil
}

func (r *InMemoryTodoListRepo) Duplicate(_ context.Context, sourceListID uuid.UUID, newName string) (*entity.TodoList, error) {
	r.mu.RLock()
	source, ok := r.lists[sourceListID]
	r.mu.RUnlock()
	if !ok || source.IsDeleted() {
		return nil, fmt.Errorf("list not found")
	}

	newList := entity.NewTodoList(source.UserID, newName)
	r.mu.Lock()
	r.lists[newList.ID] = newList
	r.mu.Unlock()
	return newList, nil
}
