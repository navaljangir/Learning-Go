package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}	

type ServerConfig struct {
	Port string
	Environment string
}

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

func getEnv(key, defaultValue string) string {
	if value, exists :=  os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func Load() *Config{
	if err := godotenv.Load(); err !=nil {
		log.Printf("No .env file found, using system enivronment varibale")
	}

	cfg := &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "todo_app"),
		},
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "3000"),
			Environment: getEnv("ENVIRONMENT", "development"),
		},
	}
	
	return cfg
}