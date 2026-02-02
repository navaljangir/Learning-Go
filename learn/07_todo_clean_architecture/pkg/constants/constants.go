package constants

// Context keys for storing values in request context
const (
	ContextUserID    = "user_id"
	ContextUsername  = "username"
	ContextRequestID = "request_id"
)

// HTTP response messages
const (
	MessageSuccess         = "success"
	MessageCreated         = "resource created successfully"
	MessageUpdated         = "resource updated successfully"
	MessageDeleted         = "resource deleted successfully"
	MessageInternalError   = "internal server error"
	MessageBadRequest      = "bad request"
	MessageUnauthorized    = "unauthorized"
	MessageForbidden       = "forbidden"
	MessageNotFound        = "resource not found"
	MessageValidationError = "validation error"
)

// Validation constraints
const (
	MinUsernameLength = 3
	MaxUsernameLength = 50
	MinPasswordLength = 8
	MaxPasswordLength = 72 // bcrypt limit
	MaxTitleLength    = 255
	MaxDescLength     = 2000
)

// Pagination defaults
const (
	DefaultPage     = 1
	DefaultPageSize = 10
	MaxPageSize     = 100
)

// Priority levels for todos
const (
	PriorityLow    = "low"
	PriorityMedium = "medium"
	PriorityHigh   = "high"
	PriorityUrgent = "urgent"
)
