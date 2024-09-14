package services

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

var rdb *redis.Client
var ctx = context.Background()

// Initialize Redis client
func InitRedis(redisURL string) {
	rdb = redis.NewClient(&redis.Options{
		Addr: redisURL,
	})
}

// Cache file metadata
func CacheMetadata(key string, value string) error {
	return rdb.Set(ctx, key, value, 5*time.Minute).Err()
}

// Retrieve cached metadata
func GetCachedMetadata(key string) (string, error) {
	return rdb.Get(ctx, key).Result()
}
