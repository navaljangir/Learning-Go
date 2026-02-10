package service

import (
	"context"
	"errors"
	"testing"
	"todo_app/domain/entity"
	"todo_app/internal/dto"
	"todo_app/pkg/utils"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// =============================================================================
// Mock UserRepository
// =============================================================================

type mockUserRepo struct {
	users map[uuid.UUID]*entity.User

	createErr              error
	findByIDErr            error
	findByUsernameErr      error
	findByEmailErr         error
	updateErr              error
	existsByUsernameErr    error
	existsByEmailErr       error
	existsByUsernameResult bool
	existsByEmailResult    bool
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{
		users: make(map[uuid.UUID]*entity.User),
	}
}

func (m *mockUserRepo) Create(_ context.Context, user *entity.User) error {
	if m.createErr != nil {
		return m.createErr
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) FindByID(_ context.Context, id uuid.UUID) (*entity.User, error) {
	if m.findByIDErr != nil {
		return nil, m.findByIDErr
	}
	user, ok := m.users[id]
	if !ok {
		return nil, utils.ErrNotFound
	}
	return user, nil
}

func (m *mockUserRepo) FindByUsername(_ context.Context, username string) (*entity.User, error) {
	if m.findByUsernameErr != nil {
		return nil, m.findByUsernameErr
	}
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, utils.ErrNotFound
}

func (m *mockUserRepo) FindByEmail(_ context.Context, email string) (*entity.User, error) {
	if m.findByEmailErr != nil {
		return nil, m.findByEmailErr
	}
	for _, user := range m.users {
		if user.Email == email {
			return user, nil
		}
	}
	return nil, utils.ErrNotFound
}

func (m *mockUserRepo) Update(_ context.Context, user *entity.User) error {
	if m.updateErr != nil {
		return m.updateErr
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockUserRepo) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	return nil, nil
}

func (m *mockUserRepo) ExistsByUsername(_ context.Context, username string) (bool, error) {
	if m.existsByUsernameErr != nil {
		return false, m.existsByUsernameErr
	}
	if m.existsByUsernameResult {
		return true, nil
	}
	for _, user := range m.users {
		if user.Username == username {
			return true, nil
		}
	}
	return false, nil
}

func (m *mockUserRepo) ExistsByEmail(_ context.Context, email string) (bool, error) {
	if m.existsByEmailErr != nil {
		return false, m.existsByEmailErr
	}
	if m.existsByEmailResult {
		return true, nil
	}
	for _, user := range m.users {
		if user.Email == email {
			return true, nil
		}
	}
	return false, nil
}

// =============================================================================
// Test helpers
// =============================================================================

// seedUser creates a user in the mock repo and returns it
func seedUser(repo *mockUserRepo, username, email, fullName string) *entity.User {
	hashedPassword, _ := utils.HashPassword("password123")
	user := entity.NewUser(username, email, hashedPassword, fullName)
	repo.users[user.ID] = user
	return user
}

// assertAppErrorUser checks that err is an *AppError with the expected status and message
func assertAppErrorUser(t *testing.T, err error, wantStatus int, wantMsg string) {
	t.Helper()
	assert.Error(t, err)
	var appErr *utils.AppError
	assert.True(t, errors.As(err, &appErr), "error should be *utils.AppError, got %T", err)
	assert.Equal(t, wantStatus, appErr.StatusCode)
	assert.Equal(t, wantMsg, appErr.Message)
}

// =============================================================================
// Register Tests
// =============================================================================

func TestRegister(t *testing.T) {
	ctx := context.Background()
	jwtUtil := utils.NewJWTUtil("test-secret", 24, "test-issuer")

	t.Run("success: create new user", func(t *testing.T) {
		userRepo := newMockUserRepo()
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.Register(ctx, dto.RegisterRequest{
			Username: "john_doe",
			Email:    "john@example.com",
			Password: "password123",
			FullName: "John Doe",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, "john_doe", resp.User.Username)
		assert.Equal(t, "john@example.com", resp.User.Email)
		assert.Equal(t, "John Doe", resp.User.FullName)
		assert.Equal(t, 1, len(userRepo.users), "user should be saved in repo")
	})

	t.Run("fail: username already exists", func(t *testing.T) {
		userRepo := newMockUserRepo()
		seedUser(userRepo, "existing_user", "existing@example.com", "Existing User")
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.Register(ctx, dto.RegisterRequest{
			Username: "existing_user",
			Email:    "new@example.com",
			Password: "password123",
			FullName: "New User",
		})

		assert.Nil(t, resp)
		assertAppErrorUser(t, err, 400, "Username already exists")
	})

	t.Run("fail: email already exists", func(t *testing.T) {
		userRepo := newMockUserRepo()
		seedUser(userRepo, "existing_user", "existing@example.com", "Existing User")
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.Register(ctx, dto.RegisterRequest{
			Username: "new_user",
			Email:    "existing@example.com",
			Password: "password123",
			FullName: "New User",
		})

		assert.Nil(t, resp)
		assertAppErrorUser(t, err, 400, "Email already exists")
	})

	t.Run("fail: repo error on username check", func(t *testing.T) {
		userRepo := newMockUserRepo()
		userRepo.existsByUsernameErr = errors.New("db error")
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.Register(ctx, dto.RegisterRequest{
			Username: "john_doe",
			Email:    "john@example.com",
			Password: "password123",
			FullName: "John Doe",
		})

		assert.Nil(t, resp)
		assertAppErrorUser(t, err, 500, "Failed to check username")
	})

	t.Run("fail: repo error on email check", func(t *testing.T) {
		userRepo := newMockUserRepo()
		userRepo.existsByEmailErr = errors.New("db error")
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.Register(ctx, dto.RegisterRequest{
			Username: "john_doe",
			Email:    "john@example.com",
			Password: "password123",
			FullName: "John Doe",
		})

		assert.Nil(t, resp)
		assertAppErrorUser(t, err, 500, "Failed to check email")
	})

	t.Run("fail: repo error on create", func(t *testing.T) {
		userRepo := newMockUserRepo()
		userRepo.createErr = errors.New("db error")
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.Register(ctx, dto.RegisterRequest{
			Username: "john_doe",
			Email:    "john@example.com",
			Password: "password123",
			FullName: "John Doe",
		})

		assert.Nil(t, resp)
		assertAppErrorUser(t, err, 500, "Failed to create user")
	})
}

// =============================================================================
// Login Tests
// =============================================================================

func TestLogin(t *testing.T) {
	ctx := context.Background()
	jwtUtil := utils.NewJWTUtil("test-secret", 24, "test-issuer")

	t.Run("success: login with correct credentials", func(t *testing.T) {
		userRepo := newMockUserRepo()
		user := seedUser(userRepo, "john_doe", "john@example.com", "John Doe")
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.Login(ctx, dto.LoginRequest{
			Username: "john_doe",
			Password: "password123",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.NotEmpty(t, resp.Token)
		assert.Equal(t, user.ID, resp.User.ID)
		assert.Equal(t, "john_doe", resp.User.Username)
	})

	t.Run("fail: user not found", func(t *testing.T) {
		userRepo := newMockUserRepo()
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.Login(ctx, dto.LoginRequest{
			Username: "nonexistent",
			Password: "password123",
		})

		assert.Nil(t, resp)
		assertAppErrorUser(t, err, 401, "Invalid credentials")
	})

	t.Run("fail: incorrect password", func(t *testing.T) {
		userRepo := newMockUserRepo()
		seedUser(userRepo, "john_doe", "john@example.com", "John Doe")
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.Login(ctx, dto.LoginRequest{
			Username: "john_doe",
			Password: "wrongpassword",
		})

		assert.Nil(t, resp)
		assertAppErrorUser(t, err, 401, "Invalid credentials")
	})

	t.Run("fail: user is deleted", func(t *testing.T) {
		userRepo := newMockUserRepo()
		user := seedUser(userRepo, "john_doe", "john@example.com", "John Doe")
		user.MarkDeleted()
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.Login(ctx, dto.LoginRequest{
			Username: "john_doe",
			Password: "password123",
		})

		assert.Nil(t, resp)
		assertAppErrorUser(t, err, 404, "Account not found")
	})
}

// =============================================================================
// GetProfile Tests
// =============================================================================

func TestGetProfile(t *testing.T) {
	ctx := context.Background()
	jwtUtil := utils.NewJWTUtil("test-secret", 24, "test-issuer")

	t.Run("success: get user profile", func(t *testing.T) {
		userRepo := newMockUserRepo()
		user := seedUser(userRepo, "john_doe", "john@example.com", "John Doe")
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.GetProfile(ctx, user.ID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, user.ID, resp.ID)
		assert.Equal(t, "john_doe", resp.Username)
		assert.Equal(t, "john@example.com", resp.Email)
		assert.Equal(t, "John Doe", resp.FullName)
	})

	t.Run("fail: user not found", func(t *testing.T) {
		userRepo := newMockUserRepo()
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.GetProfile(ctx, uuid.New())

		assert.Nil(t, resp)
		assertAppErrorUser(t, err, 404, "User not found")
	})

	t.Run("fail: user is deleted", func(t *testing.T) {
		userRepo := newMockUserRepo()
		user := seedUser(userRepo, "john_doe", "john@example.com", "John Doe")
		user.MarkDeleted()
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.GetProfile(ctx, user.ID)

		assert.Nil(t, resp)
		assertAppErrorUser(t, err, 404, "User not found")
	})

	t.Run("fail: repo error", func(t *testing.T) {
		userRepo := newMockUserRepo()
		userRepo.findByIDErr = errors.New("db error")
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.GetProfile(ctx, uuid.New())

		assert.Nil(t, resp)
		assertAppErrorUser(t, err, 404, "User not found")
	})
}

// =============================================================================
// UpdateProfile Tests
// =============================================================================

func TestUpdateProfile(t *testing.T) {
	ctx := context.Background()
	jwtUtil := utils.NewJWTUtil("test-secret", 24, "test-issuer")

	t.Run("success: update user profile", func(t *testing.T) {
		userRepo := newMockUserRepo()
		user := seedUser(userRepo, "john_doe", "john@example.com", "John Doe")
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.UpdateProfile(ctx, user.ID, dto.UpdateUserRequest{
			FullName: "John Smith",
		})

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "John Smith", resp.FullName)
		assert.Equal(t, "john_doe", resp.Username)      // Username should not change
		assert.Equal(t, "john@example.com", resp.Email) // Email should not change
	})

	t.Run("fail: user not found", func(t *testing.T) {
		userRepo := newMockUserRepo()
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.UpdateProfile(ctx, uuid.New(), dto.UpdateUserRequest{
			FullName: "John Smith",
		})

		assert.Nil(t, resp)
		assertAppErrorUser(t, err, 404, "User not found")
	})

	t.Run("fail: user is deleted", func(t *testing.T) {
		userRepo := newMockUserRepo()
		user := seedUser(userRepo, "john_doe", "john@example.com", "John Doe")
		user.MarkDeleted()
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.UpdateProfile(ctx, user.ID, dto.UpdateUserRequest{
			FullName: "John Smith",
		})

		assert.Nil(t, resp)
		assertAppErrorUser(t, err, 404, "User not found")
	})

	t.Run("fail: repo error on update", func(t *testing.T) {
		userRepo := newMockUserRepo()
		user := seedUser(userRepo, "john_doe", "john@example.com", "John Doe")
		userRepo.updateErr = errors.New("db error")
		svc := NewUserService(userRepo, jwtUtil)

		resp, err := svc.UpdateProfile(ctx, user.ID, dto.UpdateUserRequest{
			FullName: "John Smith",
		})

		assert.Nil(t, resp)
		assertAppErrorUser(t, err, 500, "Failed to update user")
	})
}
