package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	redisClient *redis.Client
	ctx         = context.Background()
)

// Initialize sets up the Redis client
func Initialize() error {
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Test connection
	_, err := redisClient.Ping(ctx).Result()
	return err
}

// Set stores a value in Redis with expiration
func Set(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}

	return redisClient.Set(ctx, key, jsonData, expiration).Err()
}

// Get retrieves a value from Redis and unmarshals it into the provided interface
func Get(key string, dest interface{}) error {
	val, err := redisClient.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil // Key does not exist
	}
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(val), dest)
}

// Delete removes a key from Redis
func Delete(key string) error {
	return redisClient.Del(ctx, key).Err()
}

// DeletePattern removes all keys matching the pattern
func DeletePattern(pattern string) error {
	keys, err := redisClient.Keys(ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return redisClient.Del(ctx, keys...).Err()
	}
	return nil
}

// GetOrSet attempts to get a value from cache, if not found calls the provider function
func GetOrSet(key string, dest interface{}, expiration time.Duration, provider func() (interface{}, error)) error {
	// Try to get from cache first
	err := Get(key, dest)
	if err != nil {
		return err
	}

	// If we got a value, return
	if dest != nil {
		return nil
	}

	// Get fresh data from provider
	data, err := provider()
	if err != nil {
		return err
	}

	// Store in cache
	if err := Set(key, data, expiration); err != nil {
		return err
	}

	// Update destination with fresh data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonData, dest)
}
