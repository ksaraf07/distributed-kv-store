package store

import "sync"

// Store is a concurrency-safe in-memory key-value store.
type Store struct {
	mu   sync.RWMutex
	data map[string]string
}

// New creates an empty Store, ready to use.
func New() *Store {
	return &Store{data: make(map[string]string)}
}

// Set adds or overwrites a key's value.
func (s *Store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Get retrieves a key's value. ok is false if the key doesn't exist.
func (s *Store) Get(key string) (value string, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok = s.data[key]
	return value, ok
}

// Delete removes a key. It's a no-op if the key doesn't exist.
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}
