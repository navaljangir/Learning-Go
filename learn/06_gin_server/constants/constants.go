package constants

import "time"

// Server configuration
const (
	ServerPort         = ":8080"
	ServerReadTimeout  = 10 * time.Second
	ServerWriteTimeout = 10 * time.Second
)

// JWT configuration
const (
	JWTSecretKey     = "your-super-secret-key-change-in-production" // Change this in production!
	JWTExpiryHours   = 24
	JWTIssuer        = "gin_server"
)

// HTTP Status messages
const (
	MsgSuccess           = "success"
	MsgInternalError     = "internal server error"
	MsgUnauthorized      = "unauthorized"
	MsgBadRequest        = "bad request"
	MsgNotFound          = "not found"
	MsgUserExists        = "user already exists"
	MsgInvalidCredentials = "invalid credentials"
	MsgTokenRequired     = "authorization token required"
	MsgTokenInvalid      = "invalid or expired token"
)

// Context keys
const (
	ContextUserID   = "userID"
	ContextUsername = "username"
)
