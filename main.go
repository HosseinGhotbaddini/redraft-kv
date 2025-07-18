package main

import (
	"log"
	"os"

	"github.com/HosseinGhotbaddini/redraft-kv/config"
	"github.com/HosseinGhotbaddini/redraft-kv/raft"
	"github.com/HosseinGhotbaddini/redraft-kv/server"
	"github.com/HosseinGhotbaddini/redraft-kv/store"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: go run main.go <config-file>")
	}
	configPath := os.Args[1]

	// Load config file
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize store backend (memory or bolt)
	var kvStore store.Store
	switch cfg.StoreBackend {
	case "bolt":
		kvStore = store.NewBoltStore(cfg.StorePath)
	default:
		kvStore = store.NewMemoryStore()
	}
	defer kvStore.Close()

	log.Printf("Initialized store backend: %T", kvStore)

	// Initialize Raft node (no static peers)
	raftNode, err := raft.NewRaftNode(cfg.ID, cfg.RaftAddr, kvStore)
	if err != nil {
		log.Fatalf("Failed to start Raft node: %v", err)
	}

	// Start Redis-compatible server
	if err := server.Start(cfg.RedisAddr, cfg.ID, raftNode, kvStore); err != nil {
		log.Fatalf("Redis server error: %v", err)
	}
}
