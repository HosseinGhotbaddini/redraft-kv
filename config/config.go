package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the node's runtime configuration.
type Config struct {
	ID           string            `yaml:"id"`            // Unique identifier for the Raft node
	RaftAddr     string            `yaml:"raft_addr"`     // Local address for Raft transport (e.g. 127.0.0.1:7001)
	RedisAddr    string            `yaml:"redis_addr"`    // Address to expose the Redis-compatible interface
	Peers        map[string]string `yaml:"peers"`         // Map of peer node IDs to their Raft addresses
	StoreBackend string            `yaml:"store_backend"` // "bolt" or "memory"
	StorePath    string            `yaml:"store_path"`    // Path to store BoltDB file (only used if backend is bolt)
}

// Load reads a YAML config file from the given path and returns a Config struct.
func Load(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
