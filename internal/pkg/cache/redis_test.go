package cache_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/OutOfStack/game-library/internal/pkg/cache"
	cache_mock "github.com/OutOfStack/game-library/internal/pkg/cache/mocks"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	goredis "github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

type redisTestValue struct {
	ID   int64
	Name string
}

func TestRedisStore_Get_CacheHit(t *testing.T) {
	ctx := t.Context()
	client := cache_mock.NewMockRedisClient(gomock.NewController(t))
	core, logs := observer.New(zap.ErrorLevel)
	store := cache.NewRedisStore(client, zap.New(core))
	key := td.String()
	want := redisTestValue{ID: td.Int64(), Name: td.String()}
	var got redisTestValue

	client.EXPECT().GetStruct(ctx, key, &got).SetArg(2, want).Return(nil)

	err := store.Get(ctx, key, &got, func() (redisTestValue, error) {
		t.Fatal("cache hit must not call the loader")
		return redisTestValue{}, nil
	}, time.Minute)

	require.NoError(t, err)
	require.Equal(t, want, got)
	require.Zero(t, logs.Len())
}

func TestRedisStore_Get_LoadAndCache(t *testing.T) {
	readErr := errors.New("redis read failed")
	writeErr := errors.New("redis write failed")
	tests := []struct {
		name     string
		getErr   error
		setErr   error
		ttl      time.Duration
		zero     bool
		wantLogs []string
	}{
		{name: "cache miss", getErr: goredis.Nil, ttl: time.Minute},
		{name: "wrapped cache miss", getErr: fmt.Errorf("cache lookup: %w", goredis.Nil), ttl: time.Minute},
		{name: "zero TTL", getErr: goredis.Nil},
		{name: "zero value", getErr: goredis.Nil, ttl: time.Minute, zero: true},
		{
			name: "read failure", getErr: readErr, ttl: time.Minute,
			wantLogs: []string{"get item from redis cache"},
		},
		{
			name: "write failure", getErr: goredis.Nil, setErr: writeErr, ttl: time.Minute,
			wantLogs: []string{"set item to redis cache"},
		},
		{
			name: "read and write failures", getErr: readErr, setErr: writeErr, ttl: time.Minute,
			wantLogs: []string{"get item from redis cache", "set item to redis cache"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			client := cache_mock.NewMockRedisClient(gomock.NewController(t))
			core, logs := observer.New(zap.ErrorLevel)
			store := cache.NewRedisStore(client, zap.New(core))
			key := td.String()
			want := redisTestValue{ID: td.Int64(), Name: td.String()}
			if tt.zero {
				want = redisTestValue{}
			}
			got := redisTestValue{ID: td.Int64(), Name: td.String()}

			gomock.InOrder(
				client.EXPECT().GetStruct(ctx, key, &got).Return(tt.getErr),
				client.EXPECT().SetStruct(ctx, key, want, tt.ttl).Return(tt.setErr),
			)

			calls := 0
			err := store.Get(ctx, key, &got, func() (redisTestValue, error) {
				calls++
				return want, nil
			}, tt.ttl)

			require.NoError(t, err)
			require.Equal(t, want, got)
			require.Equal(t, 1, calls)
			entries := logs.All()
			require.Len(t, entries, len(tt.wantLogs))
			for i, message := range tt.wantLogs {
				require.Equal(t, message, entries[i].Message)
				require.Equal(t, key, entries[i].ContextMap()["key"])
			}
		})
	}
}

func TestRedisStore_Get_LoaderError(t *testing.T) {
	tests := []struct {
		name   string
		getErr error
	}{
		{name: "cache miss", getErr: goredis.Nil},
		{name: "read failure", getErr: errors.New("redis read failed")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			client := cache_mock.NewMockRedisClient(gomock.NewController(t))
			store := cache.NewRedisStore(client, zap.NewNop())
			key := td.String()
			initial := redisTestValue{ID: td.Int64(), Name: td.String()}
			got := initial
			loadErr := errors.New("load failed")

			client.EXPECT().GetStruct(ctx, key, &got).Return(tt.getErr)

			calls := 0
			err := store.Get(ctx, key, &got, func() (redisTestValue, error) {
				calls++
				return redisTestValue{}, loadErr
			}, time.Minute)

			require.ErrorIs(t, err, loadErr)
			require.Equal(t, initial, got)
			require.Equal(t, 1, calls)
		})
	}
}

func TestRedisStore_Delete(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "success"},
		{name: "failure", err: errors.New("redis delete failed")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			client := cache_mock.NewMockRedisClient(gomock.NewController(t))
			store := cache.NewRedisStore(client, zap.NewNop())
			key := td.String()

			client.EXPECT().Delete(ctx, key).Return(tt.err)

			err := store.Delete(ctx, key)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRedisStore_DeleteByStartsWith(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{name: "success"},
		{name: "failure", err: errors.New("redis delete failed")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := t.Context()
			client := cache_mock.NewMockRedisClient(gomock.NewController(t))
			store := cache.NewRedisStore(client, zap.NewNop())
			prefix := td.String() + ":"

			client.EXPECT().DeleteByMatch(ctx, prefix+"*").Return(tt.err)

			err := store.DeleteByStartsWith(ctx, prefix)
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
