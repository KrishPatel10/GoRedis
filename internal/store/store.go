package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/KrishPatel10/GoRedis.git/internal/aof"
	"github.com/KrishPatel10/GoRedis.git/internal/resp"
)

type item struct {
	value     string
	expiresAt int64
}

type MemoryStore struct {
	mu   sync.RWMutex
	data map[string]item
	aof  *aof.AOF
}

func New(aof *aof.AOF) *MemoryStore {
	store := &MemoryStore{
		data: make(map[string]item),
		aof:  aof,
	}

	go store.startSweeping(time.Minute)

	go store.Recover(aof)

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
		s.mu.Lock() // Acquire Write Lock
		delete(s.data, key)
		s.mu.Unlock() // Release Write Lock
		return "", false
	}

	return val.value, true
}

func (s *MemoryStore) Set(key string, val string, expiryDuration time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	item := item{
		value:     val,
		expiresAt: time.Now().Add(expiryDuration).Unix(),
	}

	s.data[key] = item
}

func (s *MemoryStore) SetWithoutExpiry(key string, val string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[key] = item{
		value:     val,
		expiresAt: 0,
	}

	// NEW: Write the command to the disk!
	if s.aof != nil {
		// Manually formatting a RESP array for the AOF file
		cmd := fmt.Sprintf("*3\r\n$3\r\nSET\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n",
			len(key), key, len(val), val)
		s.aof.Write([]byte(cmd))
	}
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
	s.mu.Lock()
	defer s.mu.Unlock()

	// Only process 100 keys at a time so we don't freeze the database
	processed := 0
	now := time.Now().UnixNano()

	for key, val := range s.data {
		if processed >= 100 {
			break
		}
		if val.expiresAt > 0 && now > val.expiresAt {
			delete(s.data, key)
		}
		processed++
	}
}

func (s *MemoryStore) Restore(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = item{value: value, expiresAt: 0}
}

func (s *MemoryStore) Recover(aof *aof.AOF) error {
	fmt.Print("Initializing the recovery phase")

	err := aof.Read(func(value resp.Value) {
		if value.Typ == "array" && len(value.Array) >= 3 {
			command := value.Array[0].Value
			key := value.Array[1].Value
			value := value.Array[2].Value

			if command == "SET" {
				s.Restore(key, value)
			}
		}
	})

	if err != nil {
		fmt.Printf("Warning during the recovery phase %s", err)
		return err
	}

	return nil
}
