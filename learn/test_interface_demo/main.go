package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"demo/domain/repository"
	"demo/internal/repository/file"
	"demo/internal/repository/memory"
	"demo/internal/repository/redis"
)

func main() {
	ctx := context.Background()

	// ========================================================================
	// CREATE DIFFERENT IMPLEMENTATIONS
	// ========================================================================

	// Implementation 1: In-Memory (has StorageInfo)
	memoryRepo := memory.NewInMemoryTodoRepository()

	// Implementation 2: PostgreSQL (has StorageInfo + BatchCapable)
	// postgresRepo := postgres.NewPostgresTodRepository(db)  // Commented out (needs real DB)

	// Implementation 3: File (NO optional interfaces)
	fileRepo := file.NewFileTodoRepository("/tmp/todos.json")

	// Implementation 4: Redis (has StorageInfo + CacheCapable)
	redisClient := redis.NewRedisClient()
	redisRepo := redis.NewRedisTodoRepository(redisClient)

	// ========================================================================
	// DEMONSTRATE SWITCHING BETWEEN IMPLEMENTATIONS
	// ========================================================================

	fmt.Println("=== Testing Different Implementations ===")

	// Test each implementation
	testRepository("Memory", memoryRepo, ctx)
	fmt.Println()

	testRepository("File", fileRepo, ctx)
	fmt.Println()

	testRepository("Redis", redisRepo, ctx)
	fmt.Println()

	// ========================================================================
	// DEMONSTRATE BULK IMPORT (checks for BatchCapable)
	// ========================================================================

	fmt.Println("=== Bulk Import Demo ===")
	fmt.Println()

	todos := []*repository.Todo{
		{ID: "bulk1", Title: "Todo 1", Completed: false, CreatedAt: time.Now()},
		{ID: "bulk2", Title: "Todo 2", Completed: false, CreatedAt: time.Now()},
		{ID: "bulk3", Title: "Todo 3", Completed: false, CreatedAt: time.Now()},
	}

	bulkImport(ctx, memoryRepo, todos) // Memory: NO batch support
	bulkImport(ctx, fileRepo, todos)   // File: NO batch support
	// bulkImport(ctx, postgresRepo, todos) // Postgres: HAS batch support (commented out)

	fmt.Println()

	// ========================================================================
	// DEMONSTRATE CACHE MANAGEMENT (checks for CacheCapable)
	// ========================================================================

	fmt.Println("=== Cache Management Demo ===")
	fmt.Println()

	manageCache(memoryRepo) // Memory: NO cache support
	manageCache(redisRepo)  // Redis: HAS cache support
}

// testRepository demonstrates using any TodoRepository implementation
func testRepository(name string, repo repository.TodoRepository, ctx context.Context) {
	fmt.Printf("--- Testing %s Repository ---\n", name)

	// Create a todo (works with ALL implementations)
	todo := &repository.Todo{
		ID:        "test-" + name,
		Title:     "Test Todo from " + name,
		Completed: false,
		CreatedAt: time.Now(),
	}

	if err := repo.Create(ctx, todo); err != nil {
		log.Printf("Error creating todo: %v\n", err)
	} else {
		fmt.Println("✓ Created todo")
	}

	// Find the todo (works with ALL implementations)
	found, err := repo.FindByID(ctx, todo.ID)
	if err != nil {
		log.Printf("Error finding todo: %v\n", err)
	} else {
		fmt.Printf("✓ Found todo: %s\n", found.Title)
	}

	// Check for OPTIONAL StorageInfo interface
	if info, ok := repo.(repository.StorageInfo); ok {
		fmt.Printf("✓ Storage type: %s\n", info.GetStorageType())
		fmt.Printf("✓ Stats: %v\n", info.GetStats())
	} else {
		fmt.Println("✗ StorageInfo not supported")
	}
}

// bulkImport demonstrates checking for BatchCapable interface
func bulkImport(ctx context.Context, repo repository.TodoRepository, todos []*repository.Todo) {
	// Check if repository supports batch operations (OPTIONAL feature)
	if batchRepo, ok := repo.(repository.BatchCapable); ok {
		// Use efficient batch insert (Postgres)
		fmt.Println("✓ Using BatchCreate (efficient)")
		if err := batchRepo.BatchCreate(ctx, todos); err != nil {
			log.Printf("Batch create failed: %v\n", err)
		}
	} else {
		// Fallback: Insert one by one (Memory, File)
		fmt.Println("✗ BatchCapable not supported, using loop (slower)")
		for _, todo := range todos {
			if err := repo.Create(ctx, todo); err != nil {
				log.Printf("Create failed: %v\n", err)
			}
		}
	}
}

// manageCache demonstrates checking for CacheCapable interface
func manageCache(repo repository.TodoRepository) {
	// Check if repository supports cache operations (OPTIONAL feature)
	if cacheRepo, ok := repo.(repository.CacheCapable); ok {
		fmt.Println("✓ Cache operations available")
		fmt.Printf("  - Hit rate: %.2f%%\n", cacheRepo.GetCacheHitRate()*100)

		// Clear cache
		if err := cacheRepo.ClearCache(); err != nil {
			log.Printf("Clear cache failed: %v\n", err)
		} else {
			fmt.Println("  - Cache cleared")
		}

		// Check TTL
		ttl, err := cacheRepo.GetTTL("test-Redis")
		if err != nil {
			fmt.Printf("  - TTL check: %v\n", err)
		} else {
			fmt.Printf("  - TTL: %v\n", ttl)
		}
	} else {
		fmt.Println("✗ CacheCapable not supported")
	}

	fmt.Println()
}

// ============================================================================
// REAL-WORLD EXAMPLE: Configuration-based switching
// ============================================================================

// GetRepository creates the appropriate repository based on config
func GetRepository(storageType string) repository.TodoRepository {
	switch storageType {
	case "memory":
		return memory.NewInMemoryTodoRepository()

	case "file":
		return file.NewFileTodoRepository("/tmp/todos.json")

	// case "postgres":
	// 	db := connectToPostgres()
	// 	return postgres.NewPostgresTodRepository(db)

	// case "redis":
	// 	client := connectToRedis()
	// 	return redis.NewRedisTodoRepository(client)

	default:
		return memory.NewInMemoryTodoRepository() // Default to memory
	}
}

// PrintCapabilities shows which optional features a repository has
func PrintCapabilities(repo repository.TodoRepository) {
	fmt.Println("Repository Capabilities:")
	fmt.Println("  ✓ TodoRepository (required)")

	if _, ok := repo.(repository.StorageInfo); ok {
		fmt.Println("  ✓ StorageInfo (optional)")
	} else {
		fmt.Println("  ✗ StorageInfo (optional)")
	}

	if _, ok := repo.(repository.BatchCapable); ok {
		fmt.Println("  ✓ BatchCapable (optional)")
	} else {
		fmt.Println("  ✗ BatchCapable (optional)")
	}

	if _, ok := repo.(repository.CacheCapable); ok {
		fmt.Println("  ✓ CacheCapable (optional)")
	} else {
		fmt.Println("  ✗ CacheCapable (optional)")
	}
}
