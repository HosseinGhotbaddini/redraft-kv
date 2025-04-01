package raft

import (
	"encoding/json"
	"io"

	"github.com/HosseinGhotbaddini/redraft-kv/store"
	hraft "github.com/hashicorp/raft"
)

// FSMImpl implements hashicorp/raft.FSM.
// It defines how replicated log entries are applied to the key-value store.
type FSMImpl struct {
	Store *store.Store
}

// Apply is called once a log entry is committed by the Raft cluster.
// This is the only place where the store is mutated to ensure deterministic state.
func (f *FSMImpl) Apply(logEntry *hraft.Log) interface{} {
	var cmd store.Command
	if err := json.Unmarshal(logEntry.Data, &cmd); err != nil {
		return err
	}

	// Dispatch based on the operation type
	switch cmd.Op {
	case "set":
		f.Store.Set(cmd.Key, cmd.Value)
	case "delete":
		f.Store.Delete(cmd.Key)
	}
	return nil
}

// Snapshot returns a snapshot implementation. This stub doesn't persist anything.
func (f *FSMImpl) Snapshot() (hraft.FSMSnapshot, error) {
	return &noopSnapshot{}, nil
}

// Restore is called when a snapshot is loaded (stubbed here).
func (f *FSMImpl) Restore(rc io.ReadCloser) error {
	return nil
}

// noopSnapshot is a stub that cancels the snapshot sink immediately.
type noopSnapshot struct{}

// Persist attempts to write a snapshot but immediately cancels it.
func (n *noopSnapshot) Persist(sink hraft.SnapshotSink) error {
	_ = sink.Cancel()
	return nil
}

// Release is a no-op for this implementation.
func (n *noopSnapshot) Release() {}
