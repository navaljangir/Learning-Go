package dto

type CreateTodoRequest struct {
	Title       string  `json:"title" binding:"required,min=1,max=150"`
	Description string  `json:"description" binding:"max=300"`
	Priority    string  `json:"priority" binding:"oneof=low medium high"`
	DueDate     *string `json:"due_date" binding:"omitempty,datetime=2026-02-20T15:04:05Z07:00"`
	ListID      *string `json:"list_id" binding:"omitempty,uuid4"`
}

type UpdateTodoRequest struct {
	Title       *string `json:"title" binding:"omitempty,min=1,max=150"`
	Description *string `json:"description" binding:"omitempty,max=300"`
	Priority    *string `json:"priority" binding:"omitempty,oneof=low medium high"`
	DueDate     *string `json:"due_date" binding:"omitempty,datetime=2026-02-20T15:04:05Z07:00"`
	ListID      *string `json:"list_id" binding:"omitempty,uuid4"`
	Completed   *bool   `json:"completed,omitempty"`
}

type TodoResponse struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	ListID      *string `json:"list_id,omitempty"`
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Completed   bool    `json:"completed"`
	Priority    string  `json:"priority,omitempty"`
	DueDate     *string `json:"due_date,omitempty"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
	CompletedAt *string `json:"completed_at,omitempty"`
}
