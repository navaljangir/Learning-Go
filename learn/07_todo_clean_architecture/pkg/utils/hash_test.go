package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

// TestHashPassword tests successful password hashing
// What we're testing: Can we hash a password and get a valid bcrypt hash?
func TestHashPassword(t *testing.T) {
	// ARRANGE: Prepare test password
	password := "MySecurePassword123!"

	// ACT: Hash the password
	hash, err := HashPassword(password)

	// ASSERT: Verify successful hashing
	assert.NoError(t, err, "hashing should succeed")
	assert.NotEmpty(t, hash, "hash should not be empty")
	assert.NotEqual(t, password, hash, "hash should differ from password")

	// Verify it's a valid bcrypt hash (starts with $2a$, $2b$, or $2y$)
	assert.True(t, strings.HasPrefix(hash, "$2"), "should be valid bcrypt hash")
}

// TestHashPasswordTooLong tests that passwords over 72 chars are rejected
// What we're testing: bcrypt has a 72 character limit
func TestHashPasswordTooLong(t *testing.T) {
	// ARRANGE: Create a password with 73 characters (exceeds bcrypt limit)
	password := strings.Repeat("a", 73)

	// ACT: Try to hash
	hash, err := HashPassword(password)

	// ASSERT: Should fail with appropriate error
	assert.Error(t, err, "should return error for password > 72 chars")
	assert.Empty(t, hash, "hash should be empty on error")
	assert.Contains(t, err.Error(), "password too long", "error should mention length issue")
}

// TestHashPasswordEmptyString tests hashing empty password
// What we're testing: Empty passwords should still hash (library allows it)
func TestHashPasswordEmptyString(t *testing.T) {
	// ARRANGE
	password := ""

	// ACT
	hash, err := HashPassword(password)

	// ASSERT: bcrypt allows empty passwords
	assert.NoError(t, err, "empty password should hash")
	assert.NotEmpty(t, hash, "should produce a hash")
}

// TestHashPasswordUniqueness tests that same password produces different hashes
// What we're testing: bcrypt uses salt, so same password = different hashes
func TestHashPasswordUniqueness(t *testing.T) {
	// ARRANGE
	password := "SamePassword123!"

	// ACT: Hash the same password twice
	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	// ASSERT: Both succeed but produce different hashes (due to random salt)
	assert.NoError(t, err1)
	assert.NoError(t, err2)
	assert.NotEqual(t, hash1, hash2, "same password should produce different hashes")

	// NOTE: bcrypt automatically generates random salt for each hash
	// NOTE: This is a security feature - prevents rainbow table attacks
}

// TestCheckPasswordCorrect tests validating correct password
// What we're testing: Correct password should pass validation
func TestCheckPasswordCorrect(t *testing.T) {
	// ARRANGE: Hash a password
	password := "CorrectPassword123!"
	hash, _ := HashPassword(password)

	// ACT: Check the correct password
	result := CheckPassword(password, hash)

	// ASSERT: Should return true
	assert.True(t, result, "correct password should validate")
}

// TestCheckPasswordIncorrect tests validating wrong password
// What we're testing: Wrong password should fail validation
func TestCheckPasswordIncorrect(t *testing.T) {
	// ARRANGE: Hash a password
	correctPassword := "CorrectPassword123!"
	wrongPassword := "WrongPassword123!"
	hash, _ := HashPassword(correctPassword)

	// ACT: Check wrong password
	result := CheckPassword(wrongPassword, hash)

	// ASSERT: Should return false
	assert.False(t, result, "wrong password should not validate")
}

// TestCheckPasswordEmptyInputs tests edge cases with empty strings
func TestCheckPasswordEmptyInputs(t *testing.T) {
	// Test cases for empty inputs
	tests := []struct {
		name     string
		password string
		hash     string
		expected bool
	}{
		{
			name:     "empty password with valid hash",
			password: "",
			hash:     "$2a$10$abc123...", // some hash
			expected: false,
		},
		{
			name:     "valid password with empty hash",
			password: "password123",
			hash:     "",
			expected: false,
		},
		{
			name:     "both empty",
			password: "",
			hash:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPassword(tt.password, tt.hash)
			assert.Equal(t, tt.expected, result, tt.name)
		})
	}
}

// TestCheckPasswordInvalidHash tests checking against invalid hash format
func TestCheckPasswordInvalidHash(t *testing.T) {
	// ARRANGE
	password := "password123"
	invalidHash := "not-a-valid-bcrypt-hash"

	// ACT
	result := CheckPassword(password, invalidHash)

	// ASSERT: Should return false (bcrypt validation fails)
	assert.False(t, result, "invalid hash should return false")
}

// TestHashPasswordDifferentLengths tests various password lengths
func TestHashPasswordDifferentLengths(t *testing.T) {
	// Test passwords of different lengths (all within bcrypt limit)
	passwordLengths := []int{1, 8, 16, 32, 64, 72}

	for _, length := range passwordLengths {
		t.Run(string(rune(length))+" characters", func(t *testing.T) {
			// ARRANGE: Create password of specific length
			password := strings.Repeat("a", length)

			// ACT: Hash it
			hash, err := HashPassword(password)

			// ASSERT: Should succeed
			assert.NoError(t, err, "password of %d chars should hash", length)
			assert.NotEmpty(t, hash)

			// Verify it validates correctly
			assert.True(t, CheckPassword(password, hash), "should validate correctly")
		})
	}
}

// TestCheckPasswordWithActualHash tests with real bcrypt hash
func TestCheckPasswordWithActualHash(t *testing.T) {
	// ARRANGE: Use a pre-generated bcrypt hash
	// This is the hash for "testpassword" with bcrypt cost 10
	password := "testpassword"

	// Generate a fresh hash
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	// ACT & ASSERT: Test correct password
	assert.True(t, CheckPassword(password, string(hash)), "correct password should match")

	// ACT & ASSERT: Test wrong password
	assert.False(t, CheckPassword("wrongpassword", string(hash)), "wrong password should not match")
}

// TestHashPasswordSpecialCharacters tests passwords with special chars
func TestHashPasswordSpecialCharacters(t *testing.T) {
	// Test passwords with various special characters
	passwords := []string{
		"Pass@word123!",
		"P@$$w0rd#2024",
		"Tëst™Pàss©wørd",
		"密码测试123",   // Chinese characters
		"パスワード123",  // Japanese characters
		"пароль123", // Cyrillic
	}

	for _, password := range passwords {
		t.Run(password, func(t *testing.T) {
			// ACT: Hash and verify
			hash, err := HashPassword(password)

			// ASSERT
			assert.NoError(t, err, "should hash password with special chars")
			assert.True(t, CheckPassword(password, hash), "should validate correctly")
		})
	}
}

// BenchmarkHashPassword benchmarks password hashing performance
// Run with: go test -bench=BenchmarkHashPassword
func BenchmarkHashPassword(b *testing.B) {
	password := "BenchmarkPassword123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashPassword(password)
	}
	// NOTE: bcrypt is intentionally slow (CPU-intensive) for security
	// NOTE: Typical result: ~50-100ms per hash with default cost
}

// BenchmarkCheckPassword benchmarks password validation performance
func BenchmarkCheckPassword(b *testing.B) {
	password := "BenchmarkPassword123!"
	hash, _ := HashPassword(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		CheckPassword(password, hash)
	}
}
