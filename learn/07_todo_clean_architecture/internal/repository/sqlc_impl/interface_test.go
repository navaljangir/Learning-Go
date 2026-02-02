package sqlc_impl

import (
	"database/sql"
	"testing"
	"todo_app/domain/repository"
)

// TestUserRepositoryImplementsInterface verifies that userRepository implements the interface
func TestUserRepositoryImplementsInterface(t *testing.T) {
	var _ repository.UserRepository = (*userRepository)(nil)
}

// TestTodoRepositoryImplementsInterface verifies that todoRepository implements the interface
func TestTodoRepositoryImplementsInterface(t *testing.T) {
	var _ repository.TodoRepository = (*todoRepository)(nil)
}

// TestNewUserRepositoryReturnsInterface verifies NewUserRepository returns the interface
func TestNewUserRepositoryReturnsInterface(t *testing.T) {
	db := &sql.DB{} // Mock DB for type checking only
	repo := NewUserRepository(db)

	// This would fail to compile if repo doesn't implement repository.UserRepository
	var _ repository.UserRepository = repo
}

// TestNewTodoRepositoryReturnsInterface verifies NewTodoRepository returns the interface
func TestNewTodoRepositoryReturnsInterface(t *testing.T) {
	db := &sql.DB{} // Mock DB for type checking only
	repo := NewTodoRepository(db)

	// This would fail to compile if repo doesn't implement repository.TodoRepository
	var _ repository.TodoRepository = repo
}
