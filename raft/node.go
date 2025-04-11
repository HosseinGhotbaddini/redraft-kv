package raft

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	hraft "github.com/hashicorp/raft"
	raftboltdb "github.com/hashicorp/raft-boltdb"

	"github.com/HosseinGhotbaddini/redraft-kv/store"
)

type Node struct {
	raft *hraft.Raft
}

// NewRaftNode sets up and starts a real Raft node using hashicorp/raft.
// - id:       local node ID
// - raftAddr: local Raft TCP bind address (e.g., 127.0.0.1:7001)
// - peers:    map of other node IDs to their raft addresses
// - kv:       store backend (memory or bolt)
func NewRaftNode(id, raftAddr string, peers map[string]string, kv store.Store) (*Node, error) {
	// Create a data directory per node (for logs, snapshots, etc.)
	dataDir := filepath.Join("data", id)
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, err
	}

	// Configure Raft
	config := hraft.DefaultConfig()
	config.LocalID = hraft.ServerID(id)

	// Set up log storage, stable storage, and snapshot store
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(dataDir, "raft-log.db"))
	if err != nil {
		return nil, err
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(dataDir, "raft-stable.db"))
	if err != nil {
		return nil, err
	}
	snapshots, err := hraft.NewFileSnapshotStore(dataDir, 1, os.Stderr)
	if err != nil {
		return nil, err
	}

	// Resolve the Raft bind address
	addr, err := net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid raft address: %w", err)
	}

	// Create a TCP transport for Raft communication
	transport, err := hraft.NewTCPTransport(raftAddr, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	// Initialize the FSM which applies commands to the store
	fsm := &FSMImpl{Store: kv}

	// Create the full Raft system
	r, err := hraft.NewRaft(config, fsm, logStore, stableStore, snapshots, transport)
	if err != nil {
		return nil, err
	}

	// Build static cluster configuration (including this node and peers)
	servers := []hraft.Server{
		{ID: hraft.ServerID(id), Address: transport.LocalAddr()},
	}
	for peerID, peerAddr := range peers {
		if peerID == id {
			continue // skip self
		}
		servers = append(servers, hraft.Server{
			ID:      hraft.ServerID(peerID),
			Address: hraft.ServerAddress(peerAddr),
		})
	}

	// Bootstrap the Raft cluster with the static configuration
	cfg := hraft.Configuration{Servers: servers}
	future := r.BootstrapCluster(cfg)
	if err := future.Error(); err != nil {
		return nil, fmt.Errorf("bootstrap error: %w", err)
	}

	return &Node{raft: r}, nil
}

// Apply sends a command to Raft to be replicated and applied to the FSM
func (n *Node) Apply(cmd []byte) error {
	f := n.raft.Apply(cmd, 5*time.Second)
	return f.Error()
}
