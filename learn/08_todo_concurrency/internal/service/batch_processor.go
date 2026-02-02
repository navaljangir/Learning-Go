package service

import (
	"context"
	"fmt"
	"sync"
	"time"
	"todo_concurrency/internal/dto"
)

// BatchProcessor handles concurrent batch operations
//
// KEY LEARNING - GOROUTINES & CHANNELS:
// This demonstrates the worker pool pattern using goroutines and channels.
// Multiple workers process jobs concurrently, communicating via channels.
type BatchProcessor struct {
	todoService *TodoService
	workerCount int
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor(todoService *TodoService, workerCount int) *BatchProcessor {
	return &BatchProcessor{
		todoService: todoService,
		workerCount: workerCount,
	}
}

// ProcessBatch processes multiple todos concurrently
//
// LEARNING CONCEPTS:
// 1. Buffered channels - jobs and results channels can hold multiple items
// 2. Worker pool - fixed number of goroutines process jobs from a queue
// 3. Fan-out/Fan-in - distribute work (fan-out), collect results (fan-in)
// 4. WaitGroup - wait for all workers to finish
func (bp *BatchProcessor) ProcessBatch(ctx context.Context, requests []dto.CreateTodoRequest) dto.BatchCreateResponse {
	startTime := time.Now()

	// Create buffered channels
	// CHANNEL LEARNING:
	// Buffered channels can hold N items before blocking.
	// This allows senders and receivers to work at different speeds.
	jobs := make(chan dto.CreateTodoRequest, len(requests))
	results := make(chan dto.BatchResult, len(requests))

	// Start worker goroutines
	// GOROUTINE LEARNING:
	// Each worker runs concurrently in its own goroutine.
	// The 'go' keyword launches a lightweight thread.
	var wg sync.WaitGroup
	for w := 0; w < bp.workerCount; w++ {
		wg.Add(1)
		go bp.worker(ctx, w+1, jobs, results, &wg)
	}

	// Send jobs to workers (producer)
	// CHANNEL LEARNING:
	// Sending to a channel: channel <- value
	// This distributes work among workers
	for i, req := range requests {
		// Add index to track which request this is
		reqWithIndex := req
		jobs <- reqWithIndex
		fmt.Printf("📤 Sent job %d to queue\n", i+1)
	}
	close(jobs) // Close channel to signal no more jobs

	// Wait for all workers to finish
	go func() {
		wg.Wait()
		close(results) // Close results after all workers done
	}()

	// Collect results (consumer)
	// CHANNEL LEARNING:
	// Receiving from a channel: value <- channel
	// Range over channel until it's closed
	var allResults []dto.BatchResult
	successCount := 0
	failureCount := 0

	for result := range results {
		allResults = append(allResults, result)
		if result.Success {
			successCount++
		} else {
			failureCount++
		}
	}

	elapsed := time.Since(startTime)

	return dto.BatchCreateResponse{
		SuccessCount: successCount,
		FailureCount: failureCount,
		Results:      allResults,
		TimeElapsed:  elapsed.String(),
	}
}

// worker processes jobs from the jobs channel
// GOROUTINE LEARNING:
// Each worker runs in its own goroutine and processes jobs independently.
// Multiple workers can run simultaneously on different CPU cores.
func (bp *BatchProcessor) worker(ctx context.Context, id int, jobs <-chan dto.CreateTodoRequest, results chan<- dto.BatchResult, wg *sync.WaitGroup) {
	defer wg.Done() // Decrement counter when worker exits

	// CHANNEL LEARNING:
	// '<-chan' is a receive-only channel (can't send to it)
	// 'chan<-' is a send-only channel (can't receive from it)
	// This provides compile-time safety

	for req := range jobs {
		fmt.Printf("🔨 Worker %d processing job: %s\n", id, req.Title)

		// Simulate some processing time
		time.Sleep(100 * time.Millisecond)

		// Process the todo
		todo, err := bp.todoService.Create(ctx, req)

		result := dto.BatchResult{}
		if err != nil {
			result.Success = false
			result.Error = err.Error()
			fmt.Printf("❌ Worker %d failed: %s\n", id, err)
		} else {
			result.Success = true
			result.TodoID = todo.ID
			fmt.Printf("✅ Worker %d completed: %s (ID: %s)\n", id, req.Title, todo.ID)
		}

		results <- result
	}

	fmt.Printf("👋 Worker %d finished\n", id)
}

// ProcessBatchV2 demonstrates alternative pattern with error channel
//
// ADVANCED PATTERN:
// Uses separate channels for results and errors
// This pattern is useful when you want to handle errors differently
func (bp *BatchProcessor) ProcessBatchV2(ctx context.Context, requests []dto.CreateTodoRequest) dto.BatchCreateResponse {
	startTime := time.Now()

	// Two separate channels for different types of results
	successChan := make(chan dto.BatchResult, len(requests))
	errorChan := make(chan dto.BatchResult, len(requests))

	// Semaphore pattern - limit concurrent goroutines
	// CONCURRENCY CONTROL:
	// Buffered channel acts as a semaphore to limit concurrent operations
	sem := make(chan struct{}, bp.workerCount)

	var wg sync.WaitGroup

	// Launch one goroutine per request (not a worker pool)
	for i, req := range requests {
		wg.Add(1)
		go func(index int, request dto.CreateTodoRequest) {
			defer wg.Done()

			// Acquire semaphore (blocks if full)
			sem <- struct{}{}
			defer func() { <-sem }() // Release semaphore

			todo, err := bp.todoService.Create(ctx, request)

			result := dto.BatchResult{Index: index}
			if err != nil {
				result.Success = false
				result.Error = err.Error()
				errorChan <- result
			} else {
				result.Success = true
				result.TodoID = todo.ID
				successChan <- result
			}
		}(i, req)
	}

	// Close channels after all goroutines finish
	go func() {
		wg.Wait()
		close(successChan)
		close(errorChan)
	}()

	// Collect from both channels
	var allResults []dto.BatchResult
	successCount := 0
	failureCount := 0

	// Use select to receive from multiple channels
	// SELECT LEARNING:
	// Select statement waits for channel operations
	// Whichever channel is ready first, that case executes
	for successChan != nil || errorChan != nil {
		select {
		case result, ok := <-successChan:
			if !ok {
				successChan = nil // Channel closed
				continue
			}
			allResults = append(allResults, result)
			successCount++

		case result, ok := <-errorChan:
			if !ok {
				errorChan = nil // Channel closed
				continue
			}
			allResults = append(allResults, result)
			failureCount++
		}
	}

	return dto.BatchCreateResponse{
		SuccessCount: successCount,
		FailureCount: failureCount,
		Results:      allResults,
		TimeElapsed:  time.Since(startTime).String(),
	}
}
