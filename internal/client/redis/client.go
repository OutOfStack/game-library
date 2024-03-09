package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/OutOfStack/game-library/internal/appconf"
	"github.com/redis/go-redis/v9"
)

var (
	// ErrInvalidType represents error when invalid type provided
	ErrInvalidType = errors.New("invalid type provided")
)

// Client represents redis client
type Client struct {
	rdb *redis.Client
	ttl time.Duration
}

// New creates new redis client instance
func New(cfg appconf.Redis) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       0, // use default DB
	})

	ttl, err := time.ParseDuration(cfg.TTL)
	if err != nil {
		return nil, err
	}

	return &Client{
		rdb: rdb,
		ttl: ttl,
	}, nil
}

// Get returns value by key. If key does not exist, returns empty string without error
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	value, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return value, nil
}

// Set sets value for provided key with ttl. If ttl == 0, default ttl is used
func (c *Client) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	err := c.rdb.Set(ctx, key, value, c.getTTL(ttl)).Err()
	if err != nil {
		return err
	}
	return nil
}

// GetStruct gets value by key from cache and sets value. If key does not exist, returns nil without error.
// If value cannot be unmarshalled, InvalidTypeError returned
func (c *Client) GetStruct(ctx context.Context, key string, value interface{}) error {
	bytes, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, &value)
}

// SetStruct sets struct value for provided key with ttl. If ttl == 0, default ttl is used.
// If value cannot be marshalled, InvalidTypeError returned
func (c *Client) SetStruct(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return ErrInvalidType
	}
	err = c.rdb.Set(ctx, key, bytes, c.getTTL(ttl)).Err()
	if err != nil {
		return err
	}
	return nil
}

// Delete removes value by key
func (c *Client) Delete(ctx context.Context, key string) error {
	return c.rdb.Del(ctx, key).Err()
}

// DeleteByMatch removes value by pattern
func (c *Client) DeleteByMatch(ctx context.Context, pattern string) error {
	iterator := c.rdb.Scan(ctx, 0, pattern, 0).Iterator()
	var errGroup error
	for iterator.Next(ctx) {
		err := c.rdb.Del(ctx, iterator.Val()).Err()
		if err != nil {
			errGroup = errors.Join(errGroup, err)
		}
	}
	return errGroup
}

func (c *Client) getTTL(ttl time.Duration) time.Duration {
	if ttl == 0 {
		return c.ttl
	}
	return ttl
}
