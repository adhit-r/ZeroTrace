package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache provides unified caching interface for Valkey (Redis-compatible)
// Uses github.com/redis/go-redis/v9 which works with Valkey
type Cache struct {
	client *redis.Client // Redis client (compatible with Valkey)
	prefix string
}

// NewCache creates a new unified cache instance
func NewCache(client *redis.Client, prefix string) *Cache {
	return &Cache{
		client: client,
		prefix: prefix,
	}
}

// Get retrieves a value from cache
func (c *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	fullKey := c.prefix + key
	val, err := c.client.Get(ctx, fullKey).Bytes()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, fmt.Errorf("cache get failed: %w", err)
	}
	return val, nil
}

// GetString retrieves a string value from cache
func (c *Cache) GetString(ctx context.Context, key string) (string, error) {
	fullKey := c.prefix + key
	val, err := c.client.Get(ctx, fullKey).Result()
	if err == redis.Nil {
		return "", ErrCacheMiss
	}
	if err != nil {
		return "", fmt.Errorf("cache get failed: %w", err)
	}
	return val, nil
}

// Set stores a value in cache with TTL
func (c *Cache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	fullKey := c.prefix + key
	
	var val string
	switch v := value.(type) {
	case string:
		val = v
	case []byte:
		val = string(v)
	default:
		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal cache value: %w", err)
		}
		val = string(data)
	}
	
	return c.client.Set(ctx, fullKey, val, ttl).Err()
}

// Delete removes a key from cache
func (c *Cache) Delete(ctx context.Context, key string) error {
	fullKey := c.prefix + key
	return c.client.Del(ctx, fullKey).Err()
}

// DeletePattern removes all keys matching a pattern
func (c *Cache) DeletePattern(ctx context.Context, pattern string) error {
	iter := c.client.Scan(ctx, 0, c.prefix+pattern, 0).Iterator()
	for iter.Next(ctx) {
		if err := c.client.Del(ctx, iter.Val()).Err(); err != nil {
			return fmt.Errorf("failed to delete key %s: %w", iter.Val(), err)
		}
	}
	return iter.Err()
}

// Exists checks if a key exists in cache
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	fullKey := c.prefix + key
	count, err := c.client.Exists(ctx, fullKey).Result()
	if err != nil {
		return false, fmt.Errorf("cache exists check failed: %w", err)
	}
	return count > 0, nil
}

// WarmCache preloads frequently accessed data into cache
func (c *Cache) WarmCache(ctx context.Context, items map[string]interface{}, ttl time.Duration) error {
	pipe := c.client.Pipeline()
	
	for key, value := range items {
		fullKey := c.prefix + key
		var val string
		switch v := value.(type) {
		case string:
			val = v
		case []byte:
			val = string(v)
		default:
			data, err := json.Marshal(value)
			if err != nil {
				continue
			}
			val = string(data)
		}
		pipe.Set(ctx, fullKey, val, ttl)
	}
	
	_, err := pipe.Exec(ctx)
	return err
}

// InvalidatePattern invalidates all keys matching a pattern
func (c *Cache) InvalidatePattern(ctx context.Context, pattern string) error {
	return c.DeletePattern(ctx, pattern)
}

// PublishInvalidation publishes cache invalidation message via pub/sub
func (c *Cache) PublishInvalidation(ctx context.Context, channel string, keys []string) error {
	data, err := json.Marshal(map[string]interface{}{
		"keys": keys,
		"timestamp": time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("failed to marshal invalidation message: %w", err)
	}
	
	return c.client.Publish(ctx, c.prefix+channel, data).Err()
}

// SubscribeInvalidation subscribes to cache invalidation messages
func (c *Cache) SubscribeInvalidation(ctx context.Context, channel string) (*redis.PubSub, error) {
	pubsub := c.client.Subscribe(ctx, c.prefix+channel)
	return pubsub, nil
}

// GetStats returns cache statistics
func (c *Cache) GetStats(ctx context.Context) (map[string]interface{}, error) {
	info, err := c.client.Info(ctx, "stats").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache stats: %w", err)
	}
	
	stats := make(map[string]interface{})
	// Parse info string (simplified - in production use proper parser)
	stats["info"] = info
	
	dbSize, err := c.client.DBSize(ctx).Result()
	if err == nil {
		stats["db_size"] = dbSize
	}
	
	return stats, nil
}

var ErrCacheMiss = fmt.Errorf("cache miss")

