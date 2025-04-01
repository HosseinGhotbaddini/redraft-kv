package store

import "sync"

// Store is a simple in-memory key-value store.
// It uses a read-write mutex to allow safe concurrent access.
type Store struct {
	data map[string][]byte
	mu   sync.RWMutex
}

// Command defines the structure of a replicated operation
// applied by the Raft FSM. It is serialized and passed through Raft logs.
type Command struct {
	Op    string // "set" or "delete"
	Key   string
	Value []byte
}

// New initializes and returns a new Store instance.
func New() *Store {
	return &Store{
		data: make(map[string][]byte),
	}
}

// Called only by the FSM once a Raft log entry is committed.
func (s *Store) Set(key string, value []byte) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

// Used by the Redis server to serve read operations.
func (s *Store) Get(key string) ([]byte, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	val, ok := s.data[key]
	return val, ok
}

// Called only by the FSM after log commit.
func (s *Store) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
}
