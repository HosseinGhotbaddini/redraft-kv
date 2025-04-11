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

// NewRaftNode sets up and starts a Raft node without static cluster bootstrapping.
func NewRaftNode(id, raftAddr string, kv store.Store) (*Node, error) {
	dataDir := filepath.Join("data", id)
	if err := os.MkdirAll(dataDir, 0700); err != nil {
		return nil, err
	}

	config := hraft.DefaultConfig()
	config.LocalID = hraft.ServerID(id)

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

	addr, err := net.ResolveTCPAddr("tcp", raftAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid raft address: %w", err)
	}

	transport, err := hraft.NewTCPTransport(raftAddr, addr, 3, 10*time.Second, os.Stderr)
	if err != nil {
		return nil, err
	}

	fsm := &FSMImpl{Store: kv}

	r, err := hraft.NewRaft(config, fsm, logStore, stableStore, snapshots, transport)
	if err != nil {
		return nil, err
	}

	return &Node{raft: r}, nil
}

// Apply sends a command to Raft to be replicated and applied to the FSM.
func (n *Node) Apply(cmd []byte) error {
	f := n.raft.Apply(cmd, 5*time.Second)
	return f.Error()
}

// BootstrapSelf initializes the cluster with this node as the first voter.
func (n *Node) BootstrapSelf(id, addr string) error {
	cfg := hraft.Configuration{
		Servers: []hraft.Server{
			{ID: hraft.ServerID(id), Address: hraft.ServerAddress(addr)},
		},
	}
	future := n.raft.BootstrapCluster(cfg)
	return future.Error()
}

// JoinNode adds a new voter to the existing Raft cluster.
func (n *Node) JoinNode(id, addr string) error {
	future := n.raft.AddVoter(hraft.ServerID(id), hraft.ServerAddress(addr), 0, 0)
	return future.Error()
}

// GetLeader returns the current leader's ID and address.
func (n *Node) GetLeader() (string, string) {
	leaderAddr := n.raft.Leader()
	if leaderAddr == "" {
		return "", ""
	}

	f := n.raft.GetConfiguration()
	if err := f.Error(); err != nil {
		return "", ""
	}

	for _, srv := range f.Configuration().Servers {
		if srv.Address == leaderAddr {
			return string(srv.ID), string(srv.Address)
		}
	}
	return "", ""
}

// ListPeers returns the current configuration of the Raft cluster.
func (n *Node) ListPeers() []hraft.Server {
	f := n.raft.GetConfiguration()
	if err := f.Error(); err != nil {
		return nil
	}
	return f.Configuration().Servers
}
