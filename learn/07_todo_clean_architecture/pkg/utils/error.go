package utils

import "errors"

// Custom error types for the application
var (
	// ErrNotFound is returned when a resource is not found
	ErrNotFound = errors.New("resource not found")

	// ErrDuplicateKey is returned when a unique constraint is violated
	ErrDuplicateKey = errors.New("duplicate key violation")

	// ErrInvalidCredentials is returned when authentication fails
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrUnauthorized is returned when user is not authorized
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden is returned when user doesn't have permission
	ErrForbidden = errors.New("forbidden")

	// ErrBadRequest is returned when request is malformed
	ErrBadRequest = errors.New("bad request")

	// ErrInternalServer is returned for unexpected errors
	ErrInternalServer = errors.New("internal server error")

	// ErrValidation is returned when validation fails
	ErrValidation = errors.New("validation error")
)

// AppError represents an application-level error with additional context
type AppError struct {
	Err     error
	Message string
	Code    int
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

// NewAppError creates a new application error
func NewAppError(err error, message string, code int) *AppError {
	return &AppError{
		Err:     err,
		Message: message,
		Code:    code,
	}
}
