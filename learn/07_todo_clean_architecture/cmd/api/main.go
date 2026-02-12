package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo_app/api/handler"
	"todo_app/api/router"
	"todo_app/config"
	"todo_app/domain/repository"
	"todo_app/domain/service"
	db_repo "todo_app/repository"
	"todo_app/repository/sqlc_impl"
	serviceImpl "todo_app/service"
	"todo_app/pkg/utils"
	customValidator "todo_app/pkg/validator"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

// init runs automatically before main()
// Used for one-time initialization that doesn't need error handling in main
func init() {
	// Load environment variables from .env file
	// In production, you'd typically use real environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}
}

func main() {
	// Load configuration
	cfg := config.Load()

	// Configure Gin mode (production = no debug logs)
	configureGinMode(cfg.Server.Environment)

	// Register custom validators
	registerCustomValidators()

	// Initialize database connection (Infrastructure layer)
	db := initDatabase(cfg)
	defer db.Close()

	// Initialize repositories (Infrastructure layer)
	userRepo, todoRepo, listRepo := initRepositories(db)

	// Initialize utilities
	jwtUtil := initJWT(cfg)

	// Initialize services (Application layer)
	userService, todoService, listService := initServices(userRepo, todoRepo, listRepo, jwtUtil, cfg.JWT.Secret)

	// Initialize handlers (Presentation layer)
	authHandler, userHandler, todoHandler, listHandler := initHandlers(userService, todoService, listService)

	// Setup router with all routes
	r := setupRouter(authHandler, userHandler, todoHandler, listHandler, jwtUtil)

	// Create and start HTTP server
	srv := createServer(cfg, r)
	startServer(srv, cfg)

	// Wait for graceful shutdown
	waitForShutdown(srv)
}

// configureGinMode sets Gin to production or debug mode
func configureGinMode(environment string) {
	if environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	log.Printf("[OK] Gin mode: %s", gin.Mode())
}

// initDatabase initializes the database connection
func initDatabase(cfg *config.Config) *db_repo.Database {
	db, err := db_repo.NewDatabase(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("[OK] Database connected successfully (MySQL)")
	return db
}

// initRepositories initializes all repositories
func initRepositories(db *db_repo.Database) (repository.UserRepository, repository.TodoRepository, repository.TodoListRepository) {
	userRepo := sqlc_impl.NewUserRepository(db.DB)
	todoRepo := sqlc_impl.NewTodoRepository(db.DB)
	listRepo := sqlc_impl.NewTodoListRepository(db.DB)
	log.Println("[OK] Repositories initialized (sqlc)")
	return userRepo, todoRepo, listRepo
}

// initJWT initializes JWT utility
func initJWT(cfg *config.Config) *utils.JWTUtil {
	jwtUtil := utils.NewJWTUtil(cfg.JWT.Secret, cfg.JWT.ExpiryHours, cfg.JWT.Issuer)
	log.Println("[OK] JWT utility initialized")
	return jwtUtil
}

// initServices initializes all services
func initServices(
	userRepo repository.UserRepository,
	todoRepo repository.TodoRepository,
	listRepo repository.TodoListRepository,
	jwtUtil *utils.JWTUtil,
	shareSecret string,
) (service.UserService, service.TodoService, service.TodoListService) {
	userService := serviceImpl.NewUserService(userRepo, jwtUtil)
	todoService := serviceImpl.NewTodoService(todoRepo, listRepo)
	listService := serviceImpl.NewTodoListService(listRepo, todoRepo, userRepo, shareSecret)
	log.Println("[OK] Services initialized")
	return userService, todoService, listService
}

// initHandlers initializes all HTTP handlers
// Returns handler interfaces for better testability and flexibility
func initHandlers(
	userService service.UserService,
	todoService service.TodoService,
	listService service.TodoListService,
) (handler.AuthHandlerInterface, handler.UserHandlerInterface, handler.TodoHandlerInterface, handler.TodoListHandlerInterface) {
	authHandler := handler.NewAuthHandler(userService)
	userHandler := handler.NewUserHandler(userService)
	todoHandler := handler.NewTodoHandler(todoService)
	listHandler := handler.NewTodoListHandler(listService)
	log.Println("[OK] Handlers initialized")
	return authHandler, userHandler, todoHandler, listHandler
}

// setupRouter configures all routes
// Accepts handler interfaces for flexibility and testability
func setupRouter(
	authHandler handler.AuthHandlerInterface,
	userHandler handler.UserHandlerInterface,
	todoHandler handler.TodoHandlerInterface,
	listHandler handler.TodoListHandlerInterface,
	jwtUtil *utils.JWTUtil,
) *gin.Engine {
	r := router.SetupRouter(authHandler, userHandler, todoHandler, listHandler, jwtUtil)
	log.Println("[OK] Router configured")
	return r
}

// createServer creates HTTP server with timeouts
func createServer(cfg *config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      handler,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}
}

// startServer starts the HTTP server in a goroutine and prints startup info
func startServer(srv *http.Server, cfg *config.Config) {
	go func() {
		printStartupBanner(cfg)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
}

// printStartupBanner prints the server startup information
func printStartupBanner(cfg *config.Config) {
	separator := repeat("=", 60)
	log.Println("\n" + separator)
	log.Printf("Server starting on http://localhost%s", cfg.Server.Port)
	log.Println(separator)
	log.Println("\nAvailable Endpoints:")
	log.Println("\n  Public Endpoints:")
	log.Println("    GET    /health                    - Health check")
	log.Println("    POST   /api/v1/auth/register      - Register new user")
	log.Println("    POST   /api/v1/auth/login         - Login and get token")
	log.Println("\n  Protected Endpoints (require Bearer token):")
	log.Println("    GET    /api/v1/users/profile      - Get user profile")
	log.Println("    PUT    /api/v1/users/profile      - Update user profile")
	log.Println("\n  Todo Endpoints:")
	log.Println("    GET    /api/v1/todos              - List todos (with pagination)")
	log.Println("    POST   /api/v1/todos              - Create new todo")
	log.Println("    GET    /api/v1/todos/:id          - Get specific todo")
	log.Println("    PUT    /api/v1/todos/:id          - Update todo")
	log.Println("    PATCH  /api/v1/todos/:id/toggle   - Toggle completion")
	log.Println("    DELETE /api/v1/todos/:id          - Delete todo")
	log.Println("    PATCH  /api/v1/todos/move         - Move todos to list/global")
	log.Println("\n  List Endpoints:")
	log.Println("    GET    /api/v1/lists              - Get all lists")
	log.Println("    POST   /api/v1/lists              - Create new list")
	log.Println("    GET    /api/v1/lists/:id          - Get list with todos")
	log.Println("    PUT    /api/v1/lists/:id          - Update list (rename)")
	log.Println("    DELETE /api/v1/lists/:id          - Delete list")
	log.Println("    POST   /api/v1/lists/:id/duplicate - Duplicate list")
	log.Println("    POST   /api/v1/lists/:id/share    - Generate share link")
	log.Println("    POST   /api/v1/lists/import/:token - Import shared list")
	log.Println("\n" + separator)
	log.Printf("Environment: %s", cfg.Server.Environment)
	log.Printf("Database: %s@%s:%s/%s", cfg.Database.User, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	log.Println("\nPress Ctrl+C to stop the server")
	log.Println(separator + "\n")
}

// waitForShutdown waits for interrupt signal and performs graceful shutdown
func waitForShutdown(srv *http.Server) {
	// Create channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until signal received
	sig := <-quit
	log.Printf("\n\n[WARN] Received signal: %v", sig)
	log.Println("[SHUTDOWN] Shutting down server gracefully...")

	// Give outstanding requests 10 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("[OK] Server stopped gracefully")
}

// Helper function to repeat strings (like Python's str * n)
func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}

// registerCustomValidators registers custom validation rules with Gin's validator
func registerCustomValidators() {
	// Get Gin's validator instance (it uses go-playground/validator under the hood)
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		if err := customValidator.RegisterCustomValidators(v); err != nil {
			log.Fatalf("Failed to register custom validators: %v", err)
		}
		log.Println("[OK] Custom validators registered (nospaces, alphanumunder, strongpassword)")
	}
}
