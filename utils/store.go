package utils

import (
	"errors"
	"sync"
)

// ErrorNoSuchKey is returned when the key is not found in the store.
var ErrorNoSuchKey = errors.New("no such key")

// Store represents a store for data of any type.
type Store struct {
	mu sync.RWMutex
	m  map[string]any
}

// NewStore creates a new store.
func NewStore() *Store {
	return &Store{
		m: make(map[string]any),
	}
}

// Get returns the value for the given key.
func (s *Store) Get(key string) (any, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	val, ok := s.m[key]

	if !ok {
		return 0, ErrorNoSuchKey
	}

	return val, nil
}

// Add adds a new key-value pair to the store.
func (s *Store) Add(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.m[key] = value
}

// Delete deletes the key-value pair from the store.
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.m, key)
}

// List returns all the values in the store.
func (s *Store) List() map[string]any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.m
}

// Keys returns all the keys in the store.
func (s *Store) Keys() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]string, len(s.m))

	i := 0
	for k := range s.m {
		keys[i] = k
		i++
	}

	return keys
}

// Values returns all the values in the store.
func (s *Store) Values() []any {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data := make([]any, len(s.m))

	i := 0
	for _, v := range s.m {
		data[i] = v
		i++
	}

	return data
}

// Clear removes all the key-value pairs from the store.
func (s *Store) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.m)
}

// Len returns the number of elements in the store.
func (s *Store) Len() int {
	return len(s.m)
}
