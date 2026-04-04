package store

import (
	"sync"
	"time"
)

type item struct {
	value     string
	expiresAt int64
}

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]item
}

func New() *MemoryStore {
	store := &MemoryStore{
		data: make(map[string]item),
	}

	go store.startSweeping(time.Minute)

	return store
}

func (s *MemoryStore) Get(key string) (string, bool) {
	s.mu.RLock()
	val, exists := s.data[key]
	s.mu.RUnlock()

	if !exists {
		return "", false
	}

	if val.expiresAt > 0 && time.Now().UnixNano() > val.expiresAt {
		return "", false
	}

	return val.value, true
}

func (s *MemoryStore) startSweeping(interval time.Duration) {
	ticker := time.NewTicker(interval)

	defer ticker.Stop()

	for {
		<-ticker.C
		s.sweep()
	}
}

func (s *MemoryStore) sweep() {
	//TODO
}
