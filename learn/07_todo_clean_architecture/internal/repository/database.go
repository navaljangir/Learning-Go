package repository

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// Database wraps the sql.DB connection
type Database struct {
	*sql.DB
}

// NewDatabase creates a new database connection (MySQL)
func NewDatabase(host, port, user, password, dbname string) (*Database, error) {
	// MySQL DSN (Data Source Name) format: user:password@tcp(host:port)/dbname?params
	//
	// Query Parameters:
	// - parseTime=true: Automatically converts MySQL DATE/DATETIME/TIMESTAMP columns
	//                   to Go's time.Time type instead of []byte or string.
	//                   Without this, you'd get strings like "2024-01-15 10:30:00"
	//                   instead of proper time.Time values.
	//
	// - charset=utf8mb4: Sets the character encoding to UTF-8 (4-byte).
	//                    utf8mb4 supports ALL Unicode characters including emojis (😀),
	//                    special symbols, and characters from all languages.
	//                    Without this, emojis and some international characters would fail.
	//                    (MySQL's "utf8" is actually 3-byte and incomplete)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		user, password, host, port, dbname)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Connection pool settings
	// These control how Go manages database connections for performance and resource usage
	db.SetMaxOpenConns(25)                  // Max 25 connections to DB at once (prevents overwhelming DB)
	db.SetMaxIdleConns(5)                   // Keep 5 idle connections ready (for quick reuse)
	db.SetConnMaxLifetime(5 * time.Minute)  // Close connections after 5 min (even if idle)
	db.SetConnMaxIdleTime(10 * time.Minute) // Close idle connections after 10 min

	// Verify connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &Database{db}, nil
}
