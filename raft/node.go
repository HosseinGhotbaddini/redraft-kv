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

func NewRaftNode(id string, kv *store.Store) (*Node, error) {
	dataDir := filepath.Join("data", id)
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, err
	}

	// Create Raft config
	config := hraft.DefaultConfig()
	config.LocalID = hraft.ServerID(id)

	// Log and snapshot stores
	logStore, err := raftboltdb.NewBoltStore(filepath.Join(dataDir, "raft-log.db"))
	if err != nil {
		return nil, fmt.Errorf("log store: %w", err)
	}
	stableStore, err := raftboltdb.NewBoltStore(filepath.Join(dataDir, "raft-stable.db"))
	if err != nil {
		return nil, fmt.Errorf("stable store: %w", err)
	}
	snapshots, err := hraft.NewFileSnapshotStore(dataDir, 1, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("snapshot store: %w", err)
	}

	// In-memory TCP transport
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:0") // random port for now
	if err != nil {
		return nil, err
	}
	transport, err := hraft.NewTCPTransport(addr.String(), addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	// Initialize FSM
	fsm := &FSMImpl{Store: kv}

	// Construct Raft system
	r, err := hraft.NewRaft(config, fsm, logStore, stableStore, snapshots, transport)
	if err != nil {
		return nil, err
	}

	// Bootstrap single-node cluster (for now)
	cfg := hraft.Configuration{
		Servers: []hraft.Server{
			{
				ID:      config.LocalID,
				Address: transport.LocalAddr(),
			},
		},
	}
	r.BootstrapCluster(cfg)

	return &Node{raft: r}, nil
}

func (n *Node) Apply(cmd []byte) error {
	f := n.raft.Apply(cmd, 5*time.Second)
	return f.Error()
}
