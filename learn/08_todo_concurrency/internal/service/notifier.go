package service

import (
	"context"
	"fmt"
	"sync"
	"time"
	"todo_concurrency/domain/entity"
)

// Notification represents a notification to be sent
type Notification struct {
	TodoID  string
	Message string
	Delay   time.Duration
}

// NotificationResult represents the result of sending a notification
type NotificationResult struct {
	TodoID    string
	Success   bool
	Error     error
	Timestamp time.Time
}

// Notifier handles asynchronous notifications
//
// KEY LEARNING - GOROUTINES & CHANNELS:
// Demonstrates non-blocking operations and background processing
type Notifier struct {
	todoService *TodoService

	// Channel for sending notifications
	// BUFFERED CHANNEL: Can hold up to 100 notifications without blocking
	notificationQueue chan Notification

	// Channel for receiving results
	resultQueue chan NotificationResult

	// Protect statistics
	mu              sync.Mutex
	totalSent       int
	totalFailed     int
	isRunning       bool
	stopChan        chan struct{} // Signal to stop the worker
}

// NewNotifier creates a new notifier and starts the background worker
func NewNotifier(todoService *TodoService) *Notifier {
	n := &Notifier{
		todoService:       todoService,
		notificationQueue: make(chan Notification, 100), // Buffered channel
		resultQueue:       make(chan NotificationResult, 100),
		stopChan:          make(chan struct{}),
	}

	// Start background worker
	// GOROUTINE LEARNING:
	// This goroutine runs for the lifetime of the application
	// It continuously processes notifications from the queue
	go n.worker()

	return n
}

// SendAsync sends a notification asynchronously (non-blocking)
//
// GOROUTINE LEARNING:
// This method returns immediately. The actual notification
// is processed by the background worker goroutine.
func (n *Notifier) SendAsync(ctx context.Context, todoID string, message string, delaySeconds int) error {
	// First, verify the todo exists
	_, err := n.todoService.GetByID(ctx, todoID)
	if err != nil {
		return err
	}

	notification := Notification{
		TodoID:  todoID,
		Message: message,
		Delay:   time.Duration(delaySeconds) * time.Second,
	}

	// CHANNEL LEARNING - NON-BLOCKING SEND:
	// Use select with default to avoid blocking if queue is full
	select {
	case n.notificationQueue <- notification:
		fmt.Printf("📬 Notification queued for todo %s (will send in %ds)\n", todoID, delaySeconds)
		return nil
	default:
		return fmt.Errorf("notification queue is full")
	}
}

// worker processes notifications from the queue
// GOROUTINE LEARNING:
// This runs continuously in the background
func (n *Notifier) worker() {
	n.mu.Lock()
	n.isRunning = true
	n.mu.Unlock()

	fmt.Println("🚀 Notification worker started")

	for {
		select {
		case notification := <-n.notificationQueue:
			// Process this notification
			n.processNotification(notification)

		case <-n.stopChan:
			// Received stop signal
			fmt.Println("🛑 Notification worker stopped")
			n.mu.Lock()
			n.isRunning = false
			n.mu.Unlock()
			return
		}
	}
}

// processNotification handles a single notification
func (n *Notifier) processNotification(notification Notification) {
	// Wait for the delay
	if notification.Delay > 0 {
		fmt.Printf("⏰ Waiting %v before sending notification for todo %s...\n",
			notification.Delay, notification.TodoID)
		time.Sleep(notification.Delay)
	}

	// Simulate sending notification (in real app, this would send email/SMS/push)
	fmt.Printf("📧 Sending notification for todo %s: %s\n",
		notification.TodoID, notification.Message)

	// Simulate network delay
	time.Sleep(500 * time.Millisecond)

	// Record result
	result := NotificationResult{
		TodoID:    notification.TodoID,
		Success:   true,
		Timestamp: time.Now(),
	}

	// Update statistics
	n.mu.Lock()
	n.totalSent++
	n.mu.Unlock()

	// Send result (non-blocking)
	select {
	case n.resultQueue <- result:
	default:
		// Result queue full, drop the result
	}

	fmt.Printf("✅ Notification sent successfully for todo %s\n", notification.TodoID)
}

// GetStats returns notification statistics
func (n *Notifier) GetStats() map[string]interface{} {
	n.mu.Lock()
	defer n.mu.Unlock()

	return map[string]interface{}{
		"total_sent":       n.totalSent,
		"total_failed":     n.totalFailed,
		"queue_length":     len(n.notificationQueue),
		"results_pending":  len(n.resultQueue),
		"worker_running":   n.isRunning,
	}
}

// Stop gracefully stops the notification worker
func (n *Notifier) Stop() {
	close(n.stopChan)
}

// SendBatchAsync sends multiple notifications concurrently
//
// ADVANCED PATTERN:
// Demonstrates spawning multiple goroutines for parallel processing
func (n *Notifier) SendBatchAsync(ctx context.Context, todos []*entity.Todo, message string) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(todos))

	for _, todo := range todos {
		wg.Add(1)

		// Launch a goroutine for each notification
		// GOROUTINE LEARNING:
		// Each todo gets its own goroutine for parallel processing
		go func(t *entity.Todo) {
			defer wg.Done()

			if err := n.SendAsync(ctx, t.ID, message, 0); err != nil {
				errChan <- err
			}
		}(todo)
	}

	// Wait for all to complete
	wg.Wait()
	close(errChan)

	// Check for errors
	if len(errChan) > 0 {
		return <-errChan
	}

	return nil
}

// Example of using Context for cancellation
//
// CONTEXT LEARNING:
// Context carries deadlines, cancellation signals across API boundaries
func (n *Notifier) SendWithTimeout(ctx context.Context, todoID string, message string, timeout time.Duration) error {
	// Create context with timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resultChan := make(chan error, 1)

	go func() {
		notification := Notification{
			TodoID:  todoID,
			Message: message,
			Delay:   0,
		}

		select {
		case n.notificationQueue <- notification:
			resultChan <- nil
		case <-ctx.Done():
			resultChan <- ctx.Err()
		}
	}()

	// Wait for result or timeout
	select {
	case err := <-resultChan:
		return err
	case <-ctx.Done():
		return fmt.Errorf("notification timed out: %w", ctx.Err())
	}
}
