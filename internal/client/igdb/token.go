package igdb

import (
	"sync"
	"time"
)

type tokenInfo struct {
	token     string
	expiresAt time.Time
	mu        sync.RWMutex
}

func (t *tokenInfo) get() string {
	t.mu.RLock()
	defer t.mu.RUnlock()
	if t.expiresAt.After(time.Now().Add(5 * time.Minute)) {
		return t.token
	}
	return ""
}

func (t *tokenInfo) set(token string, expiresIn int64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.token = token
	t.expiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
}
