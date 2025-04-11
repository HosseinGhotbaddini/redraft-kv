package raft

import (
	"encoding/json"
	"io"

	"github.com/HosseinGhotbaddini/redraft-kv/store"
	hraft "github.com/hashicorp/raft"
)

type FSMImpl struct {
	Store store.Store
}

// Apply is called once a log entry is committed by the Raft cluster.
func (f *FSMImpl) Apply(logEntry *hraft.Log) interface{} {
	var cmd store.Command
	if err := json.Unmarshal(logEntry.Data, &cmd); err != nil {
		return err
	}

	switch cmd.Op {
	case "set":
		f.Store.Set(cmd.Key, cmd.Value)
	case "delete":
		f.Store.Delete(cmd.Key)
	}
	return nil
}

// Snapshot generates a point-in-time snapshot of the store.
func (f *FSMImpl) Snapshot() (hraft.FSMSnapshot, error) {
	// Dump entire store into a map
	state, err := f.Store.Dump()
	if err != nil {
		return nil, err
	}
	return &storeSnapshot{data: state}, nil
}

// Restore loads a snapshot from a reader into the store.
func (f *FSMImpl) Restore(rc io.ReadCloser) error {
	defer rc.Close()

	var snapshot map[string][]byte
	if err := json.NewDecoder(rc).Decode(&snapshot); err != nil {
		return err
	}
	return f.Store.Load(snapshot)
}
