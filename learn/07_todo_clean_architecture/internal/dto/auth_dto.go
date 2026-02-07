package dto

// RegisterRequest represents a user registration request
// Validation rules:
// - Username: 3-30 chars, starts with letter, only letters/numbers/underscores, no spaces
// - Email: valid email format, max 255 chars
// - Password: 8-72 chars, must be strong (upper, lower, number, special char)
// - FullName: required, 2-100 chars, allows spaces
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=30,alphanumunder,nospaces"`
	Email    string `json:"email" binding:"required,email,max=255"`
	Password string `json:"password" binding:"required,min=8,max=72,strongpassword"`
	FullName string `json:"full_name" binding:"required,min=2,max=100"`
}

// LoginRequest represents a user login request
// For login, we only require the fields exist (validation was done at registration)
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the response after successful login
type LoginResponse struct {
	Token     string       `json:"token"`
	User      UserResponse `json:"user"`
	ExpiresAt int64        `json:"expires_at"`
}
