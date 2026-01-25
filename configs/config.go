// Package configs handles application configuration.
// It provides structures and functions for loading and accessing configuration values.
package configs

// Config holds the application configuration.
type Config struct {
	// Server configuration
	ServerPort string

	// Database configuration (add as needed)
	// DatabaseURL string

	// Environment (development, staging, production)
	Environment string
}

// Load returns the application configuration.
// In a production app, this would load from environment variables or config files.
func Load() *Config {
	return &Config{
		ServerPort:  ":8080",
		Environment: "development",
	}
}
