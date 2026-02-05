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

type todoListRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewTodoListRepository creates a new sqlc-based todo list repository
func NewTodoListRepository(db *sql.DB) repository.TodoListRepository {
	return &todoListRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

// sqlcTodoListToEntity converts sqlc.TodoList to entity.TodoList
func sqlcTodoListToEntity(t sqlc.TodoList) *entity.TodoList {
	id, _ := uuid.Parse(t.ID)
	userID, _ := uuid.Parse(t.UserID)

	return &entity.TodoList{
		ID:        id,
		UserID:    userID,
		Name:      t.Name,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
		DeletedAt: fromNullTimePtr(t.DeletedAt),
	}
}

func (r *todoListRepository) Create(ctx context.Context, list *entity.TodoList) error {
	params := sqlc.CreateTodoListParams{
		ID:        list.ID.String(),
		UserID:    list.UserID.String(),
		Name:      list.Name,
		CreatedAt: list.CreatedAt,
		UpdatedAt: list.UpdatedAt,
	}
	return r.queries.CreateTodoList(ctx, params)
}

func (r *todoListRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.TodoList, error) {
	list, err := r.queries.GetTodoListByID(ctx, id.String())
	if err == sql.ErrNoRows {
		return nil, errors.New("list not found")
	}
	if err != nil {
		return nil, err
	}
	return sqlcTodoListToEntity(list), nil
}

func (r *todoListRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.TodoList, error) {
	lists, err := r.queries.GetTodoListsByUserID(ctx, userID.String())
	if err != nil {
		return nil, err
	}

	result := make([]*entity.TodoList, len(lists))
	for i, l := range lists {
		result[i] = sqlcTodoListToEntity(l)
	}

	return result, nil
}

func (r *todoListRepository) Update(ctx context.Context, list *entity.TodoList) error {
	params := sqlc.UpdateTodoListParams{
		ID:        list.ID.String(),
		Name:      list.Name,
		UpdatedAt: list.UpdatedAt,
	}
	return r.queries.UpdateTodoList(ctx, params)
}

func (r *todoListRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.SoftDeleteTodoList(ctx, id.String())
}

func (r *todoListRepository) CountByUser(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.CountTodoListsByUser(ctx, userID.String())
}

func (r *todoListRepository) Duplicate(ctx context.Context, sourceListID uuid.UUID, newName string) (*entity.TodoList, error) {
	// First, get the source list to verify it exists and get user_id
	sourceList, err := r.FindByID(ctx, sourceListID)
	if err != nil {
		return nil, err
	}

	// Create new list with the given name
	newList := entity.NewTodoList(sourceList.UserID, newName)

	// Save new list
	if err := r.Create(ctx, newList); err != nil {
		return nil, err
	}

	return newList, nil
}
