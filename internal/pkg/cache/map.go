package cache

import (
	"sync"
	"time"
)

// KVMap represents map with arbitrary key and value types
type KVMap[K comparable, V any] struct {
	m         map[K]V
	mutex     sync.RWMutex
	ttl       time.Duration
	expiresAt time.Time
}

// NewKVMap creates new map with arbitrary key and value types
func NewKVMap[K comparable, V any](ttl time.Duration) *KVMap[K, V] {
	if ttl == 0 {
		ttl = 7 * 24 * time.Hour
	}
	return &KVMap[K, V]{
		m:         make(map[K]V),
		expiresAt: time.Now().Add(ttl),
		ttl:       ttl,
	}
}

// Get returns value by key
func (m *KVMap[K, V]) Get(key K) (V, bool) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if time.Since(m.expiresAt) > 0 {
		go m.purge()
		return *new(V), false
	}
	val, ok := m.m[key]
	return val, ok
}

// Set sets value by key
func (m *KVMap[K, V]) Set(key K, val V) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.m[key] = val
	m.expiresAt = time.Now().Add(m.ttl)
}

// Size returns number of records
func (m *KVMap[K, V]) Size() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	if time.Since(m.expiresAt) > 0 {
		go m.purge()
		return 0
	}
	return len(m.m)
}

func (m *KVMap[K, V]) purge() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.m = make(map[K]V)
}
