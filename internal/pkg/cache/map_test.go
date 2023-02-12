package cache_test

import (
	"testing"
	"time"

	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"github.com/OutOfStack/game-library/internal/pkg/td"
	"github.com/stretchr/testify/require"
)

func TestKVMap_NotExpired_ShouldReturnValues(t *testing.T) {
	m := cache.NewKVMap[int64, int64](1 * time.Minute)
	key, value := td.Int64(), td.Int64()
	m.Set(key, value)

	gotValue, ok := m.Get(key)
	require.True(t, ok, "ok should be true")
	require.Equal(t, value, gotValue, "value should be equal")

	size := m.Size()
	require.NotEqual(t, 0, size, "size should not be 0")
}

func TestKVMap_Expired_ShouldNotReturnValues(t *testing.T) {
	m := cache.NewKVMap[int64, int64](1 * time.Millisecond)
	key, value := td.Int64(), td.Int64()
	m.Set(key, value)

	time.Sleep(2 * time.Millisecond)

	_, ok := m.Get(key)
	require.False(t, ok, "ok should be false")

	size := m.Size()
	require.Equal(t, 0, size, "size should be 0")
}
