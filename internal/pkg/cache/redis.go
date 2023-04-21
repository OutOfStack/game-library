package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/OutOfStack/game-library/internal/client/redis"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Cache represents cache service
type Cache struct {
	redis *redis.Client
	log   *zap.Logger
}

// New creates new Cache instance
func New(redis *redis.Client, log *zap.Logger) *Cache {
	return &Cache{
		redis: redis,
		log:   log,
	}
}

// Get gets value by key from cache. If key not present runs fn and stores result in cache and returns result.
// If ttl == 0, default ttl is used.
func Get[T any](ctx context.Context, c *Cache, key string, val *T, fn func() (T, error), ttl time.Duration) error {
	err := c.redis.GetStruct(ctx, key, val)
	if err == nil {
		return nil
	}
	if !errors.Is(err, goredis.Nil) {
		c.log.Error("get item from redis cache", zap.String("key", key), zap.Error(err))
	}

	res, err := fn()
	if err != nil {
		return err
	}

	*val = res

	err = c.redis.SetStruct(ctx, key, res, ttl)
	if err != nil {
		c.log.Error("set item to redis cache", zap.String("key", key), zap.Error(err))
	}

	return nil
}

// Delete removes data by key from cache
func Delete(ctx context.Context, c *Cache, key string) error {
	return c.redis.Delete(ctx, key)
}

// DeleteByStartsWith removes data with key starting with provided key from cache
func DeleteByStartsWith(ctx context.Context, c *Cache, key string) error {
	pattern := fmt.Sprintf("%s*", key)

	return c.redis.DeleteByMatch(ctx, pattern)
}
