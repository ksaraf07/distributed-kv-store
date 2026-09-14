package store

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Store struct {
	mu      sync.RWMutex
	data    map[string]string
	logFile *os.File
}

// New creates a Store backed by the log file at logPath.
// If the file already has entries, they're replayed to rebuild state.
func New(logPath string) (*Store, error) {
	s := &Store{data: make(map[string]string)}

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	if err := s.replay(f); err != nil {
		return nil, err
	}

	s.logFile = f
	return s, nil
}

// replay reads every previously logged command and re-applies it
// to rebuild the in-memory map after a restart.
func (s *Store) replay(f *os.File) error {
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		parts := strings.Fields(scanner.Text())
		if len(parts) == 0 {
			continue
		}
		switch parts[0] {
		case "SET":
			if len(parts) == 3 {
				s.data[parts[1]] = parts[2]
			}
		case "DEL":
			if len(parts) == 2 {
				delete(s.data, parts[1])
			}
		}
	}
	return scanner.Err()
}

// Set adds or overwrites a key's value, logging it to disk first.
func (s *Store) Set(key, value string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := fmt.Fprintf(s.logFile, "SET %s %s\n", key, value); err != nil {
		return err
	}
	s.data[key] = value
	return nil
}

// Get retrieves a key's value. ok is false if the key doesn't exist.
func (s *Store) Get(key string) (value string, ok bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	value, ok = s.data[key]
	return value, ok
}

// Delete removes a key, logging it to disk first.
func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, err := fmt.Fprintf(s.logFile, "DEL %s\n", key); err != nil {
		return err
	}
	delete(s.data, key)
	return nil
}
