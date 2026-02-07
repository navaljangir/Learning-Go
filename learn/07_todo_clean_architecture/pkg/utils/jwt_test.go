package utils

import (
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

// TestJWTGeneration tests creating a valid JWT token
// What we're testing: Can JWTUtil create a token successfully?
func TestJWTGeneration(t *testing.T) {
	// ARRANGE: Create JWTUtil with test configuration
	jwtUtil := NewJWTUtil("test-secret-key-12345", 24, "test-issuer")

	// ACT: Generate a token for a user
	token, expiresAt, err := jwtUtil.GenerateToken("user-123", "john_doe")

	// ASSERT: Verify the results
	assert.NoError(t, err, "token generation should not return error")
	assert.NotEmpty(t, token, "token should not be empty string")
	assert.Greater(t, expiresAt, time.Now().Unix(), "expiration should be in the future")

	// NOTE: Token format is like "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ..."
	// NOTE: It's a base64 encoded string with 3 parts separated by dots
}

// TestJWTValidation tests validating a correct token
// What we're testing: Can we validate a token we just created?
func TestJWTValidation(t *testing.T) {
	// ARRANGE: Create util and generate a token
	jwtUtil := NewJWTUtil("test-secret", 24, "test-issuer")
	token, _, err := jwtUtil.GenerateToken("user-456", "jane_doe")
	assert.NoError(t, err, "setup: token generation should succeed")

	// ACT: Validate the token
	claims, err := jwtUtil.ValidateToken(token)

	// ASSERT: Check validation succeeded and claims are correct
	assert.NoError(t, err, "validation should succeed")
	assert.NotNil(t, claims, "claims should not be nil")
	assert.Equal(t, "user-456", claims.UserID, "user ID should match")
	assert.Equal(t, "jane_doe", claims.Username, "username should match")
	assert.Equal(t, "test-issuer", claims.Issuer, "issuer should match")

	// NOTE: Claims contain the data we encoded in the token
	// NOTE: Go's JWT library automatically checks expiration during validation
}

// TestJWTValidationWithWrongSecret tests security - wrong secret should fail
// What we're testing: Tokens signed with different secrets should NOT validate
func TestJWTValidationWithWrongSecret(t *testing.T) {
	// ARRANGE: Create token with one secret
	jwtUtil1 := NewJWTUtil("secret-key-1", 24, "issuer")
	token, _, _ := jwtUtil1.GenerateToken("user-789", "attacker")

	// ACT: Try to validate with DIFFERENT secret
	jwtUtil2 := NewJWTUtil("secret-key-2", 24, "issuer")
	claims, err := jwtUtil2.ValidateToken(token)

	// ASSERT: Validation should FAIL
	assert.Error(t, err, "validation should fail with wrong secret")
	assert.Nil(t, claims, "claims should be nil for invalid token")

	// NOTE: This tests security - if someone steals a token, they can't
	// NOTE: create their own tokens without knowing the secret key
}

// TestJWTValidationWithWrongIssuer tests issuer validation
// What we're testing: Tokens from different issuers should not validate
func TestJWTValidationWithWrongIssuer(t *testing.T) {
	// ARRANGE: Create token with one issuer
	jwtUtil1 := NewJWTUtil("same-secret", 24, "issuer-app-1")
	token, _, _ := jwtUtil1.GenerateToken("user-123", "john")

	// ACT: Try to validate with DIFFERENT issuer
	jwtUtil2 := NewJWTUtil("same-secret", 24, "issuer-app-2")
	claims, err := jwtUtil2.ValidateToken(token)

	// ASSERT: Should fail
	assert.Error(t, err, "validation should fail with wrong issuer")
	assert.Nil(t, claims, "claims should be nil")

	// NOTE: Issuer check prevents tokens from one app being used in another app
}

// TestExpiredToken tests that expired tokens are rejected
// What we're testing: Old tokens should not work
func TestExpiredToken(t *testing.T) {
	// ARRANGE: Create token that expires immediately (0 hours)
	jwtUtil := NewJWTUtil("secret", 0, "issuer")
	token, _, _ := jwtUtil.GenerateToken("user-999", "olduser")

	// Wait for token to expire
	time.Sleep(time.Second * 2)

	// ACT: Try to validate expired token
	claims, err := jwtUtil.ValidateToken(token)

	// ASSERT: Should fail
	assert.Error(t, err, "validation should fail for expired token")
	assert.Nil(t, claims, "claims should be nil")
	assert.Contains(t, err.Error(), "expired", "error message should mention expiration")

	// NOTE: In real apps, tokens typically expire in 24 hours or 7 days
}

// TestTokenExpirationTime tests expiration timestamp calculation
// What we're testing: expiresAt should be approximately now + expiryHours
func TestTokenExpirationTime(t *testing.T) {
	// ARRANGE: Create util with 48 hour expiry
	jwtUtil := NewJWTUtil("secret", 48, "issuer")

	// ACT: Generate token
	beforeGeneration := time.Now()
	_, expiresAt, _ := jwtUtil.GenerateToken("user-123", "john")
	afterGeneration := time.Now()

	// ASSERT: Expiration should be ~48 hours from now
	expectedExpiryMin := beforeGeneration.Add(48 * time.Hour).Unix()
	expectedExpiryMax := afterGeneration.Add(48 * time.Hour).Unix()

	assert.GreaterOrEqual(t, expiresAt, expectedExpiryMin, "expiry should be at least 48 hours from now")
	assert.LessOrEqual(t, expiresAt, expectedExpiryMax, "expiry should be no more than 48 hours from now")

	// NOTE: We use a range because token generation takes a few milliseconds
}

// TestMultipleTokenGeneration tests that each token is unique
// What we're testing: Generating multiple tokens should create different tokens
func TestMultipleTokenGeneration(t *testing.T) {
	// ARRANGE
	jwtUtil := NewJWTUtil("secret", 24, "issuer")

	// ACT: Generate 3 tokens for same user with sufficient delay
	token1, _, _ := jwtUtil.GenerateToken("user-123", "john")
	time.Sleep(time.Second * 1) // Wait 1 second to ensure different IssuedAt timestamps
	token2, _, _ := jwtUtil.GenerateToken("user-123", "john")
	time.Sleep(time.Second * 1)
	token3, _, _ := jwtUtil.GenerateToken("user-123", "john")

	// ASSERT: All tokens should be different (because IssuedAt timestamp differs)
	assert.NotEqual(t, token1, token2, "token1 and token2 should differ")
	assert.NotEqual(t, token2, token3, "token2 and token3 should differ")
	assert.NotEqual(t, token1, token3, "token1 and token3 should differ")

	// NOTE: Even for same user, tokens differ due to IssuedAt timestamp
	// NOTE: JWT timestamps use Unix seconds (not milliseconds), so need 1 second delay
	// NOTE: Each token is valid independently
}

// TestTokenWithDifferentUsers tests tokens for different users
// What we're testing: Different users should get different tokens with correct claims
func TestTokenWithDifferentUsers(t *testing.T) {
	// ARRANGE
	jwtUtil := NewJWTUtil("secret", 24, "issuer")

	// Test cases for different users
	users := []struct {
		userID   string
		username string
	}{
		{"user-1", "alice"},
		{"user-2", "bob"},
		{"user-3", "charlie"},
	}

	for _, user := range users {
		// ACT: Generate and validate token
		token, _, err := jwtUtil.GenerateToken(user.userID, user.username)
		assert.NoError(t, err)

		claims, err := jwtUtil.ValidateToken(token)
		assert.NoError(t, err)

		// ASSERT: Claims should match input
		assert.Equal(t, user.userID, claims.UserID, "user ID should match for %s", user.username)
		assert.Equal(t, user.username, claims.Username, "username should match for %s", user.username)
	}
}

// BenchmarkJWTGeneration benchmarks token generation performance
// Run with: go test -bench=BenchmarkJWTGeneration
func BenchmarkJWTGeneration(b *testing.B) {
	jwtUtil := NewJWTUtil("secret", 24, "issuer")

	// b.N is automatically adjusted to get reliable timing
	for i := 0; i < b.N; i++ {
		jwtUtil.GenerateToken("user-123", "john")
	}
	// Example output: BenchmarkJWTGeneration-8   50000   30000 ns/op
	// Means: 50,000 iterations, each took ~30,000 nanoseconds (0.03ms)
}

// BenchmarkJWTValidation benchmarks token validation performance
func BenchmarkJWTValidation(b *testing.B) {
	jwtUtil := NewJWTUtil("secret", 24, "issuer")
	token, _, _ := jwtUtil.GenerateToken("user-123", "john")
	b.ResetTimer() // Don't count setup time
	for i := 0; i < b.N; i++ {
		jwtUtil.ValidateToken(token)
	}
}
