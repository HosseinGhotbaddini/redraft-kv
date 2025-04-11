package store

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"go.etcd.io/bbolt"
)

const defaultBucket = "kv"

type BoltStore struct {
	db *bbolt.DB
	mu sync.RWMutex
}

func NewBoltStore(path string) *BoltStore {
	// Ensure the parent directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatalf("failed to create BoltDB directory: %v", err)
	}

	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		log.Fatalf("failed to open BoltDB: %v", err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(defaultBucket))
		return err
	})
	if err != nil {
		log.Fatalf("failed to create bucket: %v", err)
	}

	return &BoltStore{db: db}
}

func (s *BoltStore) Set(key string, value []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucket))
		return b.Put([]byte(key), value)
	})
}

func (s *BoltStore) Get(key string) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var val []byte
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucket))
		v := b.Get([]byte(key))
		if v != nil {
			val = append([]byte{}, v...)
		}
		return nil
	})
	return val, err
}

func (s *BoltStore) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucket))
		return b.Delete([]byte(key))
	})
}

func (s *BoltStore) Close() error {
	return s.db.Close()
}

func (s *BoltStore) Dump() (map[string][]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	state := make(map[string][]byte)
	err := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucket))
		return b.ForEach(func(k, v []byte) error {
			state[string(k)] = append([]byte(nil), v...)
			return nil
		})
	})
	return state, err
}

func (s *BoltStore) Load(snapshot map[string][]byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucket))
		if err := b.ForEach(func(k, _ []byte) error {
			return b.Delete(k)
		}); err != nil {
			return err
		}
		for k, v := range snapshot {
			if err := b.Put([]byte(k), v); err != nil {
				return err
			}
		}
		return nil
	})
}
