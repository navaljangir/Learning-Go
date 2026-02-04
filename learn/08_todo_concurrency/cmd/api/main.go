package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo_concurrency/api/handler"
	"todo_concurrency/api/router"
	"todo_concurrency/internal/repository/memory"
	"todo_concurrency/internal/service"
)

func main() {
	printBanner()

	// Initialize repository (in-memory implementation)
	// INTERFACE LEARNING: We're using the in-memory implementation,
	// but we could easily switch to cache or database without changing the code below!
	repo := memory.NewInMemoryTodoRepository()
	log.Println("✓ Repository initialized (in-memory)")

	// Initialize services
	todoService := service.NewTodoService(repo)
	statsService := service.NewStatsService(repo)
	batchProcessor := service.NewBatchProcessor(todoService, 3) // 3 workers
	notifier := service.NewNotifier(todoService)
	log.Println("✓ Services initialized")
	log.Println("  - TodoService: Handles business logic")
	log.Println("  - StatsService: Thread-safe statistics with mutex")
	log.Println("  - BatchProcessor: Concurrent processing with worker pool")
	log.Println("  - Notifier: Async notifications with goroutines")

	// Initialize handlers
	todoHandler := handler.NewTodoHandler(todoService, statsService)
	batchHandler := handler.NewBatchHandler(batchProcessor)
	statsHandler := handler.NewStatsHandler(statsService)
	notifyHandler := handler.NewNotifyHandler(notifier)
	adminHandler := handler.NewAdminHandler(todoService, statsService)
	log.Println("✓ Handlers initialized")

	// Setup router
	r := router.SetupRouter(todoHandler, batchHandler, statsHandler, notifyHandler, adminHandler)
	log.Println("✓ Router configured")

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server in a goroutine
	// GOROUTINE LEARNING:
	// The server runs in background, main goroutine waits for shutdown signal
	go func() {
		printEndpoints()

		log.Println("\n🚀 Server starting on http://localhost:8080")
		log.Println("Press Ctrl+C to stop")

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	// CHANNEL LEARNING:
	// os.Signal is sent through a channel when user presses Ctrl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Block until signal received

	log.Println("\n🛑 Shutdown signal received...")

	// Stop notifier worker
	notifier.Stop()

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("✅ Server stopped gracefully")
}

func printBanner() {
	banner := `
╔════════════════════════════════════════════════════════════════╗
║                                                                ║
║           TODO APP - Learning Concurrency & Interfaces        ║
║                                                                ║
║  📚 This app teaches:                                         ║
║     • Interfaces - Multiple implementations of same contract  ║
║     • Goroutines - Lightweight concurrent execution           ║
║     • Channels - Communication between goroutines             ║
║     • Mutex - Thread-safe shared state                        ║
║                                                                ║
╚════════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

func printEndpoints() {
	fmt.Println("\n" + strings("=", 70))
	fmt.Println("                      AVAILABLE ENDPOINTS")
	fmt.Println(strings("=", 70))

	fmt.Println("\n📋 BASIC CRUD (Learn: Interfaces)")
	fmt.Println("  POST   /api/v1/todos              - Create todo")
	fmt.Println("  GET    /api/v1/todos              - List all todos")
	fmt.Println("  GET    /api/v1/todos/:id          - Get specific todo")
	fmt.Println("  PUT    /api/v1/todos/:id          - Update todo")
	fmt.Println("  DELETE /api/v1/todos/:id          - Delete todo")
	fmt.Println("  PATCH  /api/v1/todos/:id/toggle   - Toggle completion")

	fmt.Println("\n🔄 BATCH OPERATIONS (Learn: Goroutines + Channels)")
	fmt.Println("  POST   /api/v1/todos/batch        - Process batch (worker pool pattern)")
	fmt.Println("  POST   /api/v1/todos/batch-v2     - Process batch (semaphore pattern)")

	fmt.Println("\n📧 NOTIFICATIONS (Learn: Async Goroutines)")
	fmt.Println("  POST   /api/v1/todos/:id/notify   - Send async notification")
	fmt.Println("  GET    /api/v1/notifications/stats - Notification statistics")

	fmt.Println("\n📊 STATISTICS (Learn: Mutex)")
	fmt.Println("  GET    /api/v1/stats               - Basic statistics")
	fmt.Println("  GET    /api/v1/stats/detailed      - Detailed statistics")
	fmt.Println("  GET    /api/v1/stats/storage       - Storage-specific stats")
	fmt.Println("  GET    /api/v1/stats/goroutines    - Active goroutine count")
	fmt.Println("  POST   /api/v1/stats/reset         - Reset statistics")

	fmt.Println("\n⚙️  ADMIN (Learn: Interface Switching)")
	fmt.Println("  POST   /api/v1/admin/switch-storage - Switch storage backend")
	fmt.Println("  GET    /api/v1/admin/storage-info   - Current storage info")

	fmt.Println("\n💡 TRY THESE EXAMPLES:")
	fmt.Println("  # Create a todo")
	fmt.Println(`  curl -X POST http://localhost:8080/api/v1/todos \`)
	fmt.Println(`    -H "Content-Type: application/json" \`)
	fmt.Println(`    -d '{"title":"Learn Go","description":"Master concurrency","priority":3}'`)

	fmt.Println("\n  # Batch create (watch console for worker activity!)")
	fmt.Println(`  curl -X POST http://localhost:8080/api/v1/todos/batch \`)
	fmt.Println(`    -H "Content-Type: application/json" \`)
	fmt.Println(`    -d '{"todos":[{"title":"Task 1","priority":2},{"title":"Task 2","priority":1}]}'`)

	fmt.Println("\n  # Send async notification (returns immediately!)")
	fmt.Println(`  curl -X POST http://localhost:8080/api/v1/todos/1/notify \`)
	fmt.Println(`    -H "Content-Type: application/json" \`)
	fmt.Println(`    -d '{"message":"Don't forget!","delay_seconds":5}'`)

	fmt.Println("\n  # View statistics")
	fmt.Println(`  curl http://localhost:8080/api/v1/stats`)

	fmt.Println("\n" + strings("=", 70) + "\n")
}

func strings(s string, count int) string {
	result := ""
	for i := 0; i < count; i++ {
		result += s
	}
	return result
}
