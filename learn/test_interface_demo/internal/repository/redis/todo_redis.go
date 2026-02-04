package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"demo/domain/repository"
)

// Mock Redis client for demonstration
type RedisClient struct {
	data   map[string]string
	ttl    map[string]time.Time
	hits   int
	misses int
}

// NewRedisClient creates a new mock Redis client
func NewRedisClient() *RedisClient {
	return &RedisClient{
		data: make(map[string]string),
		ttl:  make(map[string]time.Time),
	}
}

// RedisTodoRepository stores todos in Redis cache
type RedisTodoRepository struct {
	client *RedisClient
	prefix string
}

// NewRedisTodoRepository creates a new redis repository
func NewRedisTodoRepository(client *RedisClient) repository.TodoRepository {
	return &RedisTodoRepository{
		client: client,
		prefix: "todo:",
	}
}

// ============================================================================
// REQUIRED INTERFACE IMPLEMENTATION - TodoRepository
// ============================================================================

func (r *RedisTodoRepository) Create(ctx context.Context, todo *repository.Todo) error {
	key := r.prefix + todo.ID

	data, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	r.client.data[key] = string(data)
	r.client.ttl[key] = time.Now().Add(24 * time.Hour) // 24 hour TTL

	return nil
}

func (r *RedisTodoRepository) FindByID(ctx context.Context, id string) (*repository.Todo, error) {
	key := r.prefix + id

	data, exists := r.client.data[key]
	if !exists {
		r.client.misses++
		return nil, errors.New("not found")
	}

	// Check TTL
	if time.Now().After(r.client.ttl[key]) {
		delete(r.client.data, key)
		delete(r.client.ttl, key)
		r.client.misses++
		return nil, errors.New("expired")
	}

	r.client.hits++

	var todo repository.Todo
	if err := json.Unmarshal([]byte(data), &todo); err != nil {
		return nil, err
	}

	return &todo, nil
}

func (r *RedisTodoRepository) FindAll(ctx context.Context) ([]*repository.Todo, error) {
	var todos []*repository.Todo

	for key, data := range r.client.data {
		// Check if it's a todo key
		if len(key) < len(r.prefix) || key[:len(r.prefix)] != r.prefix {
			continue
		}

		// Check TTL
		if time.Now().After(r.client.ttl[key]) {
			delete(r.client.data, key)
			delete(r.client.ttl, key)
			continue
		}

		var todo repository.Todo
		if err := json.Unmarshal([]byte(data), &todo); err != nil {
			continue
		}
		todos = append(todos, &todo)
	}

	return todos, nil
}

func (r *RedisTodoRepository) Update(ctx context.Context, todo *repository.Todo) error {
	key := r.prefix + todo.ID

	if _, exists := r.client.data[key]; !exists {
		return errors.New("not found")
	}

	data, err := json.Marshal(todo)
	if err != nil {
		return err
	}

	r.client.data[key] = string(data)
	r.client.ttl[key] = time.Now().Add(24 * time.Hour) // Reset TTL

	return nil
}

func (r *RedisTodoRepository) Delete(ctx context.Context, id string) error {
	key := r.prefix + id

	if _, exists := r.client.data[key]; !exists {
		return errors.New("not found")
	}

	delete(r.client.data, key)
	delete(r.client.ttl, key)

	return nil
}

// ============================================================================
// OPTIONAL INTERFACE IMPLEMENTATION - StorageInfo
// ============================================================================

// GetStorageType implements StorageInfo interface
func (r *RedisTodoRepository) GetStorageType() string {
	return "redis-cache"
}

// GetStats implements StorageInfo interface
func (r *RedisTodoRepository) GetStats() map[string]interface{} {
	todoCount := 0
	for key := range r.client.data {
		if len(key) >= len(r.prefix) && key[:len(r.prefix)] == r.prefix {
			todoCount++
		}
	}

	return map[string]interface{}{
		"storage_type":   "redis-cache",
		"total_todos":    todoCount,
		"cache_hits":     r.client.hits,
		"cache_misses":   r.client.misses,
		"cache_hit_rate": r.GetCacheHitRate(),
	}
}

// ============================================================================
// OPTIONAL INTERFACE IMPLEMENTATION - CacheCapable
// ============================================================================

// ClearCache implements CacheCapable interface
func (r *RedisTodoRepository) ClearCache() error {
	// Remove all todo keys
	for key := range r.client.data {
		if len(key) >= len(r.prefix) && key[:len(r.prefix)] == r.prefix {
			delete(r.client.data, key)
			delete(r.client.ttl, key)
		}
	}

	// Reset stats
	r.client.hits = 0
	r.client.misses = 0

	return nil
}

// GetCacheHitRate implements CacheCapable interface
func (r *RedisTodoRepository) GetCacheHitRate() float64 {
	total := r.client.hits + r.client.misses
	if total == 0 {
		return 0.0
	}
	return float64(r.client.hits) / float64(total)
}

// GetTTL implements CacheCapable interface
func (r *RedisTodoRepository) GetTTL(key string) (time.Duration, error) {
	fullKey := r.prefix + key

	ttl, exists := r.client.ttl[fullKey]
	if !exists {
		return 0, errors.New("key not found")
	}

	remaining := time.Until(ttl)
	if remaining < 0 {
		return 0, errors.New("key expired")
	}

	return remaining, nil
}

// ============================================================================
// COMPILE-TIME CHECKS
// ============================================================================

// Verify this type implements the required interface
var _ repository.TodoRepository = (*RedisTodoRepository)(nil)

// Verify this type implements OPTIONAL interfaces
var _ repository.StorageInfo = (*RedisTodoRepository)(nil)
var _ repository.CacheCapable = (*RedisTodoRepository)(nil)

// NOTE: We DON'T implement BatchCapable (Redis doesn't need batch operations)
// So we DON'T have a compile-time check for that
