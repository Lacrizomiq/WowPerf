package cache

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// InitCache initializes the Redis client with the provided URL or defaults to localhost:6379 if not set.
func InitCache() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		redisURL = "localhost:6379"
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
}

// Set sets a value in the cache with the given key and expiration time.
func Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	json, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return redisClient.Set(ctx, key, json, expiration).Err()
}

// Get retrieves a value from the cache and unmarshals it into the destination interface.
func Get(key string, dest interface{}) error {
	ctx := context.Background()
	val, err := redisClient.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// Delete removes a value from the cache by key.
func Delete(key string) error {
	ctx := context.Background()
	return redisClient.Del(ctx, key).Err()
}
