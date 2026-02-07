package utils

import "errors"

// Sentinel errors - use with errors.Is() for error checking
var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrForbidden is returned when user doesn't have permission to access resource
	ErrForbidden = errors.New("forbidden")

	// ErrDuplicateKey is returned when a unique constraint is violated
	ErrDuplicateKey = errors.New("duplicate key violation")

	// ErrInvalidCredentials is returned when authentication fails
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrBadRequest is returned when request is malformed
	ErrBadRequest = errors.New("bad request")
)

// AppError represents an application-level error with HTTP status code
// Services return this to specify exactly what HTTP code should be used
//
// Usage in services:
//   return &utils.AppError{Err: utils.ErrNotFound, Message: "todo not found", StatusCode: 404}
//
// Usage in middleware:
//   var appErr *utils.AppError
//   if errors.As(err, &appErr) {
//       c.JSON(appErr.StatusCode, gin.H{"error": appErr.Message})
//   }
type AppError struct {
	Err        error  // The underlying sentinel error (ErrNotFound, ErrForbidden, etc.)
	Message    string // User-friendly error message
	StatusCode int    // HTTP status code (404, 403, 500, etc.)
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "unknown error"
}

// Unwrap allows errors.Is and errors.As to work with wrapped errors
func (e *AppError) Unwrap() error {
	return e.Err
}
