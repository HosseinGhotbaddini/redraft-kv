package raft

import (
	"encoding/json"
	"io"

	"github.com/HosseinGhotbaddini/redraft-kv/store"
	hraft "github.com/hashicorp/raft"
)

type FSMImpl struct {
	Store *store.Store
}

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

func (f *FSMImpl) Snapshot() (hraft.FSMSnapshot, error) {
	return &noopSnapshot{}, nil
}

func (f *FSMImpl) Restore(rc io.ReadCloser) error {
	return nil
}

type noopSnapshot struct{}

func (n *noopSnapshot) Persist(sink hraft.SnapshotSink) error {
	_ = sink.Cancel()
	return nil
}

func (n *noopSnapshot) Release() {}
