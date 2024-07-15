package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisClient provides methods for working with redis
type RedisClient interface {
	GetStruct(ctx context.Context, key string, value interface{}) error
	SetStruct(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	DeleteByMatch(ctx context.Context, pattern string) error
	Delete(ctx context.Context, key string) error
}

// RedisStore represents redis store
type RedisStore struct {
	redisClient RedisClient
	log         *zap.Logger
}

// NewRedisStore creates new Cache instance
func NewRedisStore(redisClient RedisClient, log *zap.Logger) *RedisStore {
	return &RedisStore{
		redisClient: redisClient,
		log:         log,
	}
}

// Get gets value by key from cache. If key not present, runs fn and stores result in cache and returns result.
// If ttl == 0, default ttl is used.
func Get[T any](ctx context.Context, rs *RedisStore, key string, val *T, fn func() (T, error), ttl time.Duration) error {
	err := rs.redisClient.GetStruct(ctx, key, val)
	if err == nil {
		return nil
	}
	if !errors.Is(err, goredis.Nil) {
		rs.log.Error("get item from redis cache", zap.String("key", key), zap.Error(err))
	}

	res, err := fn()
	if err != nil {
		return err
	}

	*val = res

	err = rs.redisClient.SetStruct(ctx, key, res, ttl)
	if err != nil {
		rs.log.Error("set item to redis cache", zap.String("key", key), zap.Error(err))
	}

	return nil
}

// Delete removes data by key from cache
func Delete(ctx context.Context, c *RedisStore, key string) error {
	return c.redisClient.Delete(ctx, key)
}

// DeleteByStartsWith removes data with key starting with provided key from cache
func DeleteByStartsWith(ctx context.Context, rs *RedisStore, key string) error {
	pattern := fmt.Sprintf("%s*", key)

	return rs.redisClient.DeleteByMatch(ctx, pattern)
}
