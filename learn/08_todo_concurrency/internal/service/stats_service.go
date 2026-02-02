package service

import (
	"context"
	"runtime"
	"sync"
	"time"
	"todo_concurrency/domain/repository"
	"todo_concurrency/internal/dto"
)

// StatsService tracks application statistics in a thread-safe manner
//
// KEY LEARNING - MUTEX:
// Multiple HTTP handlers (goroutines) can call these methods concurrently.
// Mutex ensures only one goroutine modifies the counters at a time.
type StatsService struct {
	// RWMutex allows multiple readers OR one writer
	// Use RLock for reads, Lock for writes
	mu sync.RWMutex

	// Statistics (protected by mutex)
	requestCount     int64
	lastRequestTime  time.Time
	totalResponseTime time.Duration

	// Request type counters
	createRequests int64
	readRequests   int64
	updateRequests int64
	deleteRequests int64

	// Reference to repository for todo counts
	repo repository.TodoRepository
}

// NewStatsService creates a new stats service
func NewStatsService(repo repository.TodoRepository) *StatsService {
	return &StatsService{
		repo:            repo,
		lastRequestTime: time.Now(),
	}
}

// RecordRequest records a new request
//
// MUTEX LEARNING - WRITE LOCK:
// Lock() acquires exclusive access. No other goroutine can read or write
// until Unlock() is called. Always use defer to ensure unlock happens.
func (s *StatsService) RecordRequest(requestType string, duration time.Duration) {
	s.mu.Lock()                // EXCLUSIVE LOCK - blocks all other access
	defer s.mu.Unlock()        // ALWAYS unlock when done

	s.requestCount++
	s.totalResponseTime += duration
	s.lastRequestTime = time.Now()

	// Update type-specific counters
	switch requestType {
	case "create":
		s.createRequests++
	case "read":
		s.readRequests++
	case "update":
		s.updateRequests++
	case "delete":
		s.deleteRequests++
	}
}

// GetStats returns current statistics
//
// MUTEX LEARNING - READ LOCK:
// RLock() allows multiple goroutines to read simultaneously.
// This is more efficient than Lock() for read-heavy workloads.
func (s *StatsService) GetStats(ctx context.Context) (dto.StatsResponse, error) {
	// Get todo counts from repository
	totalTodos, err := s.repo.Count(ctx)
	if err != nil {
		return dto.StatsResponse{}, err
	}

	completedTodos, err := s.repo.CountCompleted(ctx)
	if err != nil {
		return dto.StatsResponse{}, err
	}

	// Calculate stats with read lock
	s.mu.RLock()              // SHARED LOCK - allows concurrent reads
	defer s.mu.RUnlock()

	pendingTodos := totalTodos - completedTodos
	completionRate := 0.0
	if totalTodos > 0 {
		completionRate = float64(completedTodos) / float64(totalTodos) * 100
	}

	// Get storage type if available
	storageType := "unknown"
	if si, ok := s.repo.(repository.StorageInfo); ok {
		storageType = si.GetStorageType()
	}

	return dto.StatsResponse{
		TotalTodos:       totalTodos,
		CompletedTodos:   completedTodos,
		PendingTodos:     pendingTodos,
		CompletionRate:   completionRate,
		ActiveGoroutines: runtime.NumGoroutine(),
		StorageType:      storageType,
	}, nil
}

// GetDetailedStats returns more detailed statistics
func (s *StatsService) GetDetailedStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	avgResponseTime := time.Duration(0)
	if s.requestCount > 0 {
		avgResponseTime = s.totalResponseTime / time.Duration(s.requestCount)
	}

	return map[string]interface{}{
		"request_count":       s.requestCount,
		"last_request_time":   s.lastRequestTime.Format(time.RFC3339),
		"avg_response_time":   avgResponseTime.String(),
		"create_requests":     s.createRequests,
		"read_requests":       s.readRequests,
		"update_requests":     s.updateRequests,
		"delete_requests":     s.deleteRequests,
		"active_goroutines":   runtime.NumGoroutine(),
		"memory_alloc_mb":     memoryUsageMB(),
	}
}

// GetStorageStats returns storage-specific statistics
func (s *StatsService) GetStorageStats() map[string]interface{} {
	if si, ok := s.repo.(repository.StorageInfo); ok {
		return si.GetStats()
	}
	return map[string]interface{}{
		"error": "storage does not provide statistics",
	}
}

// ResetStats resets all statistics
//
// MUTEX LEARNING:
// Even simple assignments need protection in concurrent code
func (s *StatsService) ResetStats() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.requestCount = 0
	s.totalResponseTime = 0
	s.createRequests = 0
	s.readRequests = 0
	s.updateRequests = 0
	s.deleteRequests = 0
	s.lastRequestTime = time.Now()
}

// IncrementCounter is a simple example of atomic counter
//
// MUTEX LEARNING:
// Without mutex, this would be a race condition:
// 1. Goroutine A reads counter = 5
// 2. Goroutine B reads counter = 5
// 3. Goroutine A writes counter = 6
// 4. Goroutine B writes counter = 6
// Result: Counter should be 7, but it's 6! Data race!
func (s *StatsService) IncrementCounter() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.requestCount++
	return s.requestCount
}

// GetRequestCount returns the request count (thread-safe read)
func (s *StatsService) GetRequestCount() int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.requestCount
}

// memoryUsageMB returns current memory usage in MB
func memoryUsageMB() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.Alloc) / 1024 / 1024
}

// Example of incorrect code (RACE CONDITION - DO NOT USE)
//
// type UnsafeCounter struct {
//     count int
// }
//
// func (c *UnsafeCounter) Increment() {
//     // RACE CONDITION: Multiple goroutines can read/write simultaneously
//     c.count++  // This is actually: temp = c.count; temp++; c.count = temp
// }
//
// Always use mutex to protect shared data!

// Example of using sync.Atomic for simple counters (alternative to mutex)
//
// ADVANCED:
// For simple integer counters, atomic operations are faster than mutex
// But for complex operations or multiple fields, use mutex
//
// import "sync/atomic"
//
// type AtomicCounter struct {
//     count int64
// }
//
// func (c *AtomicCounter) Increment() {
//     atomic.AddInt64(&c.count, 1)
// }
//
// func (c *AtomicCounter) Get() int64 {
//     return atomic.LoadInt64(&c.count)
// }
