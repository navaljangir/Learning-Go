package config

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestLoad tests loading config with default values
func TestLoad(t *testing.T) {
	clearEnvVars()
	cfg := Load()

	assert.NotNil(t, cfg)
	assert.Equal(t, ":8080", cfg.Server.Port)
	assert.Equal(t, "development", cfg.Server.Environment)
	assert.Equal(t, "localhost", cfg.Database.Host)
	assert.Equal(t, 24, cfg.JWT.ExpiryHours)
}

// TestLoadWithEnvironmentVariables tests loading with custom env vars
func TestLoadWithEnvironmentVariables(t *testing.T) {
	os.Setenv("SERVER_PORT", ":3000")
	os.Setenv("DB_HOST", "postgres.example.com")
	os.Setenv("JWT_EXPIRY_HOURS", "48")
	defer clearEnvVars()

	cfg := Load()

	assert.Equal(t, ":3000", cfg.Server.Port)
	assert.Equal(t, "postgres.example.com", cfg.Database.Host)
	assert.Equal(t, 48, cfg.JWT.ExpiryHours)
}

// TestGetEnvInt tests integer environment variable parsing
func TestGetEnvInt(t *testing.T) {
	os.Setenv("TEST_INT", "100")
	defer os.Unsetenv("TEST_INT")

	result := getEnvInt("TEST_INT", 42)
	assert.Equal(t, 100, result)

	// Test with invalid int
	os.Setenv("TEST_INT", "not-a-number")
	result = getEnvInt("TEST_INT", 42)
	assert.Equal(t, 42, result, "should fallback to default")
}

// TestGetDuration tests duration parsing from environment
func TestGetDuration(t *testing.T) {
	os.Setenv("TEST_DURATION", "60")
	defer os.Unsetenv("TEST_DURATION")

	result := getDuration("TEST_DURATION", 30)
	assert.Equal(t, 60*time.Second, result)
}

// clearEnvVars clears all config-related environment variables
func clearEnvVars() {
	envVars := []string{
		"SERVER_PORT", "SERVER_READ_TIMEOUT", "SERVER_WRITE_TIMEOUT", "ENVIRONMENT",
		"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME",
		"JWT_SECRET", "JWT_EXPIRY_HOURS", "JWT_ISSUER",
	}
	for _, key := range envVars {
		os.Unsetenv(key)
	}
}
