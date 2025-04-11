package store

import "sync"

type MemoryStore struct {
	data map[string][]byte
	mu   sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string][]byte),
	}
}

func (m *MemoryStore) Set(key string, value []byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
	return nil
}

func (m *MemoryStore) Get(key string) ([]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	val, ok := m.data[key]
	if !ok {
		return nil, nil
	}
	return val, nil
}

func (m *MemoryStore) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	return nil
}

func (m *MemoryStore) Close() error {
	return nil
}

func (m *MemoryStore) Dump() (map[string][]byte, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	clone := make(map[string][]byte, len(m.data))
	for k, v := range m.data {
		clone[k] = append([]byte(nil), v...) // deep copy
	}
	return clone, nil
}

func (m *MemoryStore) Load(snapshot map[string][]byte) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = make(map[string][]byte, len(snapshot))
	for k, v := range snapshot {
		m.data[k] = append([]byte(nil), v...) // deep copy
	}
	return nil
}
