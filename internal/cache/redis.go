package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// RedisClient wraps Redis operations
type RedisClient struct {
	client *redis.Client
	prefix string
}

// Config contains Redis configuration
type Config struct {
	Addr     string
	Password string
	DB       int
	Prefix   string
}

// NewRedisClient creates a new Redis client
func NewRedisClient(config Config) (*RedisClient, error) {
	if config.Addr == "" {
		config.Addr = "localhost:6379"
	}
	if config.Prefix == "" {
		config.Prefix = "gateway:"
	}

	client := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis connection failed: %w", err)
	}

	log.Info().Str("addr", config.Addr).Msg("Redis connected")

	return &RedisClient{
		client: client,
		prefix: config.Prefix,
	}, nil
}

// Get retrieves a value from cache
func (r *RedisClient) Get(ctx context.Context, key string, dest interface{}) error {
	fullKey := r.prefix + key

	data, err := r.client.Get(ctx, fullKey).Bytes()
	if err == redis.Nil {
		return nil // Cache miss
	}
	if err != nil {
		return fmt.Errorf("redis get failed: %w", err)
	}

	if err := json.Unmarshal(data, dest); err != nil {
		return fmt.Errorf("unmarshal failed: %w", err)
	}

	log.Debug().Str("key", fullKey).Msg("Cache hit")
	return nil
}

// Set stores a value in cache
func (r *RedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	fullKey := r.prefix + key

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("marshal failed: %w", err)
	}

	if err := r.client.Set(ctx, fullKey, data, ttl).Err(); err != nil {
		return fmt.Errorf("redis set failed: %w", err)
	}

	log.Debug().Str("key", fullKey).Dur("ttl", ttl).Msg("Cache set")
	return nil
}

// Delete removes a value from cache
func (r *RedisClient) Delete(ctx context.Context, key string) error {
	fullKey := r.prefix + key

	if err := r.client.Del(ctx, fullKey).Err(); err != nil {
		return fmt.Errorf("redis delete failed: %w", err)
	}

	log.Debug().Str("key", fullKey).Msg("Cache deleted")
	return nil
}

// DeletePattern removes all keys matching a pattern
func (r *RedisClient) DeletePattern(ctx context.Context, pattern string) error {
	fullPattern := r.prefix + pattern

	var cursor uint64
	for {
		var keys []string
		var err error

		keys, cursor, err = r.client.Scan(ctx, cursor, fullPattern, 100).Result()
		if err != nil {
			return fmt.Errorf("redis scan failed: %w", err)
		}

		if len(keys) > 0 {
			if err := r.client.Del(ctx, keys...).Err(); err != nil {
				return fmt.Errorf("redis delete failed: %w", err)
			}
		}

		if cursor == 0 {
			break
		}
	}

	log.Debug().Str("pattern", fullPattern).Msg("Cache pattern deleted")
	return nil
}

// HealthCheck checks Redis health
func (r *RedisClient) HealthCheck(ctx context.Context) error {
	if err := r.client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis ping failed: %w", err)
	}
	return nil
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	return r.client.Close()
}
