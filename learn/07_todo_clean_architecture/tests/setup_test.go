package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"todo_app/api/handler"
	"todo_app/api/router"
	serviceImpl "todo_app/internal/service"
	"todo_app/pkg/utils"
	customValidator "todo_app/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// TestMain runs once before all tests in this package.
// Registers custom validators and sets gin to test mode.
func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := customValidator.RegisterCustomValidators(v); err != nil {
			panic("failed to register custom validators: " + err.Error())
		}
	}

	os.Exit(m.Run())
}

// TestApp wires the full application stack with in-memory repositories.
type TestApp struct {
	Router   *gin.Engine
	JWTUtil  *utils.JWTUtil
	UserRepo *InMemoryUserRepo
	TodoRepo *InMemoryTodoRepo
	ListRepo *InMemoryTodoListRepo
}

// NewTestApp creates a fully wired application with in-memory repos.
func NewTestApp() *TestApp {
	// Repos
	userRepo := NewInMemoryUserRepo()
	todoRepo := NewInMemoryTodoRepo()
	listRepo := NewInMemoryTodoListRepo(todoRepo)

	// JWT
	jwtUtil := utils.NewJWTUtil("test-secret-key-for-integration", 24, "test-issuer")

	// Services (real implementations, mock repos)
	userService := serviceImpl.NewUserService(userRepo, jwtUtil)
	todoService := serviceImpl.NewTodoService(todoRepo, listRepo)
	listService := serviceImpl.NewTodoListService(listRepo, todoRepo, userRepo, "test-share-secret")

	// Handlers
	authHandler := handler.NewAuthHandler(userService)
	userHandler := handler.NewUserHandler(userService)
	todoHandler := handler.NewTodoHandler(todoService)
	listHandler := handler.NewTodoListHandler(listService)

	// Router (includes all middleware: recovery, request ID, logger, CORS, error handler)
	r := router.SetupRouter(authHandler, userHandler, todoHandler, listHandler, jwtUtil)

	return &TestApp{
		Router:   r,
		JWTUtil:  jwtUtil,
		UserRepo: userRepo,
		TodoRepo: todoRepo,
		ListRepo: listRepo,
	}
}

// RegisterUser registers a user and returns the auth token.
func (app *TestApp) RegisterUser(t *testing.T, username, email, password, fullName string) string {
	t.Helper()
	body := map[string]string{
		"username":  username,
		"email":     email,
		"password":  password,
		"full_name": fullName,
	}
	rec := app.DoRequest("POST", "/api/v1/auth/register", "", body)
	if rec.Code != http.StatusCreated {
		t.Fatalf("RegisterUser: expected 201, got %d — body: %s", rec.Code, rec.Body.String())
	}
	return extractToken(t, rec)
}

// LoginUser logs in and returns the auth token.
func (app *TestApp) LoginUser(t *testing.T, username, password string) string {
	t.Helper()
	body := map[string]string{
		"username": username,
		"password": password,
	}
	rec := app.DoRequest("POST", "/api/v1/auth/login", "", body)
	if rec.Code != http.StatusOK {
		t.Fatalf("LoginUser: expected 200, got %d — body: %s", rec.Code, rec.Body.String())
	}
	return extractToken(t, rec)
}

// DoRequest performs an HTTP request against the test router.
// If token is non-empty, it's set as a Bearer token.
// If body is non-nil, it's JSON-encoded.
func (app *TestApp) DoRequest(method, path, token string, body any) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		b, _ := json.Marshal(body)
		reqBody = bytes.NewBuffer(b)
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req := httptest.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	rec := httptest.NewRecorder()
	app.Router.ServeHTTP(rec, req)
	return rec
}

// ParseResponse unmarshals a response body into the standard envelope + data.
func ParseResponse(t *testing.T, rec *httptest.ResponseRecorder) map[string]any {
	t.Helper()
	var result map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to parse response: %v — body: %s", err, rec.Body.String())
	}
	return result
}

// extractToken pulls the token from the standard {success, data: {token: ...}} envelope.
func extractToken(t *testing.T, rec *httptest.ResponseRecorder) string {
	t.Helper()
	resp := ParseResponse(t, rec)
	data, ok := resp["data"].(map[string]any)
	if !ok {
		t.Fatalf("extractToken: missing 'data' field — body: %s", rec.Body.String())
	}
	token, ok := data["token"].(string)
	if !ok || token == "" {
		t.Fatalf("extractToken: missing 'token' field — body: %s", rec.Body.String())
	}
	return token
}
