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
	"todo_app/internal/repository"
	"todo_app/internal/repository/postgres"
	"todo_app/internal/service"
	"todo_app/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg := config.Load()

	// Set Gin mode based on environment
	if cfg.Server.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database connection
	db, err := repository.NewDatabase(
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
	)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("✓ Database connected successfully")

	// Initialize repositories (Infrastructure layer)
	userRepo := postgres.NewUserRepository(db.DB)
	todoRepo := postgres.NewTodoRepository(db.DB)
	log.Println("✓ Repositories initialized")

	// Initialize utilities
	jwtUtil := utils.NewJWTUtil(cfg.JWT.Secret, cfg.JWT.ExpiryHours, cfg.JWT.Issuer)
	log.Println("✓ JWT utility initialized")

	// Initialize services (Application layer)
	userService := service.NewUserService(userRepo, jwtUtil)
	todoService := service.NewTodoService(todoRepo)
	log.Println("✓ Services initialized")

	// Initialize handlers (Presentation layer)
	authHandler := handler.NewAuthHandler(userService)
	userHandler := handler.NewUserHandler(userService)
	todoHandler := handler.NewTodoHandler(todoService)
	log.Println("✓ Handlers initialized")

	// Setup router with all dependencies
	r := router.SetupRouter(
		authHandler,
		userHandler,
		todoHandler,
		jwtUtil,
	)
	log.Println("✓ Router configured")

	// Create HTTP server with timeouts
	srv := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// Start server in a goroutine
	go func() {
		separator := repeat("=", 60)
		log.Println("\n" + separator)
		log.Printf("Server starting on http://localhost%s", cfg.Server.Port)
		log.Println(separator)
		log.Println("\nAvailable endpoints:")
		log.Println("\n  Public Endpoints:")
		log.Println("    GET    /health                    - Health check")
		log.Println("    POST   /api/v1/auth/register      - Register new user")
		log.Println("    POST   /api/v1/auth/login         - Login and get token")
		log.Println("\n  Protected Endpoints (require Bearer token):")
		log.Println("    GET    /api/v1/users/profile      - Get user profile")
		log.Println("    PUT    /api/v1/users/profile      - Update user profile")
		log.Println("    GET    /api/v1/todos              - List todos (with pagination)")
		log.Println("    POST   /api/v1/todos              - Create new todo")
		log.Println("    GET    /api/v1/todos/:id          - Get specific todo")
		log.Println("    PUT    /api/v1/todos/:id          - Update todo")
		log.Println("    PATCH  /api/v1/todos/:id/toggle   - Toggle completion")
		log.Println("    DELETE /api/v1/todos/:id          - Delete todo")
		log.Println("\n" + separator)
		log.Printf("Environment: %s\n", cfg.Server.Environment)
		log.Println("Press Ctrl+C to stop the server")
		log.Println(separator + "\n")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	log.Printf("\n\n🛑 Received signal: %v", sig)
	log.Println("🔄 Shutting down server gracefully...")

	// Give outstanding requests 10 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✓ Server stopped gracefully")
}

// Helper function to repeat strings (like Python's str * n)
func repeat(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
