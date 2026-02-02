package dto

// CreateTodoRequest represents the data needed to create a todo
type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Priority    int    `json:"priority" binding:"required,min=1,max=3"`
}

// UpdateTodoRequest represents the data needed to update a todo
type UpdateTodoRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    int    `json:"priority" binding:"min=1,max=3"`
	Completed   *bool  `json:"completed"` // Pointer allows distinguishing between false and not provided
}

// BatchCreateRequest contains multiple todos to create
type BatchCreateRequest struct {
	Todos []CreateTodoRequest `json:"todos" binding:"required,min=1,max=100"`
}

// BatchCreateResponse contains results of batch creation
type BatchCreateResponse struct {
	SuccessCount int            `json:"success_count"`
	FailureCount int            `json:"failure_count"`
	Results      []BatchResult  `json:"results"`
	TimeElapsed  string         `json:"time_elapsed"`
}

// BatchResult represents the result of processing one todo
type BatchResult struct {
	Index   int    `json:"index"`
	Success bool   `json:"success"`
	TodoID  string `json:"todo_id,omitempty"`
	Error   string `json:"error,omitempty"`
}

// NotifyRequest represents a notification request
type NotifyRequest struct {
	Message      string `json:"message" binding:"required"`
	DelaySeconds int    `json:"delay_seconds" binding:"min=0,max=300"` // Max 5 minutes
}

// StatsResponse contains system statistics
type StatsResponse struct {
	TotalTodos       int     `json:"total_todos"`
	CompletedTodos   int     `json:"completed_todos"`
	PendingTodos     int     `json:"pending_todos"`
	CompletionRate   float64 `json:"completion_rate"`
	ActiveGoroutines int     `json:"active_goroutines"`
	StorageType      string  `json:"storage_type"`
}

// SwitchStorageRequest allows changing storage backend
type SwitchStorageRequest struct {
	Backend string `json:"backend" binding:"required,oneof=memory cache"`
}
