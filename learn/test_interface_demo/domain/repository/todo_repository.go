package repository

import (
	"context"
	"time"
)

// Todo entity (simplified)
type Todo struct {
	ID        string
	Title     string
	Completed bool
	CreatedAt time.Time
}

// ============================================================================
// REQUIRED INTERFACE - All implementations MUST have these
// ============================================================================

// TodoRepository is the REQUIRED interface
// Every implementation (memory, postgres, file) MUST implement ALL these methods
type TodoRepository interface {
	Create(ctx context.Context, todo *Todo) error
	FindByID(ctx context.Context, id string) (*Todo, error)
	FindAll(ctx context.Context) ([]*Todo, error)
	Update(ctx context.Context, todo *Todo) error
	Delete(ctx context.Context, id string) error
}

// ============================================================================
// OPTIONAL INTERFACES - Only some implementations will have these
// ============================================================================

// StorageInfo is an OPTIONAL interface
// Only implementations that can provide statistics should implement this
// Examples: Memory (has stats), Postgres (has stats), File (NO stats)
type StorageInfo interface {
	GetStorageType() string
	GetStats() map[string]interface{}
}

// BatchCapable is an OPTIONAL interface
// Only implementations that support batch operations should implement this
// Examples: Postgres (YES), Memory (NO), File (NO)
type BatchCapable interface {
	BatchCreate(ctx context.Context, todos []*Todo) error
	BatchDelete(ctx context.Context, ids []string) error
}

// CacheCapable is an OPTIONAL interface
// Only cache implementations should have this
// Examples: Redis (YES), Memory (NO), Postgres (NO)
type CacheCapable interface {
	ClearCache() error
	GetCacheHitRate() float64
	GetTTL(key string) (time.Duration, error)
}
