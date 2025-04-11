package store

type Store interface {
	Set(key string, value []byte) error
	Get(key string) ([]byte, error)
	Delete(key string) error
	Close() error

	// For Raft snapshot and restore
	Dump() (map[string][]byte, error)
	Load(snapshot map[string][]byte) error
}
