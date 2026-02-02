package sqlc_impl

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"todo_app/domain/entity"
	"todo_app/domain/repository"
	"todo_app/internal/repository/sqlc"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// Helper functions shared between repositories

type userRepository struct {
	db      *sql.DB
	queries *sqlc.Queries
}

// NewUserRepository creates a new sqlc-based user repository
func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

// toNullString converts string to sql.NullString
func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// toNullTime converts time.Time to sql.NullTime
func toNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: t, Valid: true}
}

// toNullTimePtrconverts *time.Time to sql.NullTime
func toNullTimePtr(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{Valid: false}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// fromNullString converts sql.NullString to string
func fromNullString(ns sql.NullString) string {
	if !ns.Valid {
		return ""
	}
	return ns.String
}

// fromNullTime converts sql.NullTime to time.Time
func fromNullTime(nt sql.NullTime) time.Time {
	if !nt.Valid {
		return time.Time{}
	}
	return nt.Time
}

// fromNullTimePtr converts sql.NullTime to *time.Time
func fromNullTimePtr(nt sql.NullTime) *time.Time {
	if !nt.Valid {
		return nil
	}
	return &nt.Time
}

// sqlcUserToEntity converts sqlc.User to entity.User
func sqlcUserToEntity(u sqlc.User) *entity.User {
	return &entity.User{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		FullName:     fromNullString(u.FullName),
		CreatedAt:    fromNullTime(u.CreatedAt),
		UpdatedAt:    fromNullTime(u.UpdatedAt),
		DeletedAt:    fromNullTimePtr(u.DeletedAt),
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	params := sqlc.CreateUserParams{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		FullName:     toNullString(user.FullName),
		CreatedAt:    toNullTime(user.CreatedAt),
		UpdatedAt:    toNullTime(user.UpdatedAt),
	}

	err := r.queries.CreateUser(ctx, params)
	if err != nil {
		// Check for unique constraint violations
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			return errors.New("username or email already exists")
		}
		return err
	}
	return nil
}

func (r *userRepository) FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := r.queries.GetUserByID(ctx, id)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return sqlcUserToEntity(user), nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	user, err := r.queries.GetUserByUsername(ctx, username)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return sqlcUserToEntity(user), nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := r.queries.GetUserByEmail(ctx, email)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}
	return sqlcUserToEntity(user), nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	params := sqlc.UpdateUserParams{
		ID:        user.ID,
		FullName:  toNullString(user.FullName),
		UpdatedAt: toNullTime(user.UpdatedAt),
	}
	return r.queries.UpdateUser(ctx, params)
}

func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.SoftDeleteUser(ctx, id)
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	params := sqlc.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	users, err := r.queries.ListUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	result := make([]*entity.User, len(users))
	for i, u := range users {
		result[i] = sqlcUserToEntity(u)
	}

	return result, nil
}

func (r *userRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return r.queries.CheckUsernameExists(ctx, username)
}

func (r *userRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	return r.queries.CheckEmailExists(ctx, email)
}
