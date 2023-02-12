package cache_test

import (
	"testing"
	"time"

	"github.com/OutOfStack/game-library/internal/pkg/cache"
	"github.com/OutOfStack/game-library/internal/pkg/td"
)

func TestKVMap_NotExpired_ShouldReturnValues(t *testing.T) {
	t.Parallel()

	m := cache.NewKVMap[int64, int64](1 * time.Minute)
	key, value := td.Int64(), td.Int64()
	m.Set(key, value)

	gotValue, ok := m.Get(key)
	if !ok {
		t.Errorf("Expected value %d to exist", value)
	} else if value != gotValue {
		t.Errorf("Expected to get value %d, got %d", value, gotValue)
	}

	size := m.Size()
	if size == 0 {
		t.Errorf("Expected size not to be 0")
	}
}

func TestKVMap_Expired_ShouldNotReturnValues(t *testing.T) {
	t.Parallel()

	m := cache.NewKVMap[int64, int64](1 * time.Millisecond)
	key, value := td.Int64(), td.Int64()
	m.Set(key, value)

	time.Sleep(2 * time.Millisecond)

	_, ok := m.Get(key)
	if ok {
		t.Errorf("Expected value %d to not exist", value)
	}

	size := m.Size()
	if size != 0 {
		t.Errorf("Expected size to be 0")
	}
}
