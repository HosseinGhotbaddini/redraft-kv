package main

import (
	"log"
	"os"

	"github.com/HosseinGhotbaddini/redraft-kv/raft"
	"github.com/HosseinGhotbaddini/redraft-kv/server"
	"github.com/HosseinGhotbaddini/redraft-kv/store"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run main.go <nodeID>")
	}
	nodeID := os.Args[1]

	// Initialize store
	kvStore := store.New()

	// Initialize Raft
	raftNode, err := raft.NewRaftNode(nodeID, kvStore)
	if err != nil {
		log.Fatalf("Failed to start Raft node: %v", err)
	}

	// Start Redis server
	if err := server.Start(nodeID, *raftNode, kvStore); err != nil {
		log.Fatalf("Redis server error: %v", err)
	}
}
