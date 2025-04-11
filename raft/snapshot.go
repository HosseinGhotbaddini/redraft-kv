package raft

import (
	"encoding/json"

	hraft "github.com/hashicorp/raft"
)

type storeSnapshot struct {
	data map[string][]byte
}

// Persist writes the current state to the Raft snapshot sink.
func (s *storeSnapshot) Persist(sink hraft.SnapshotSink) error {
	defer sink.Close()

	encoder := json.NewEncoder(sink)
	if err := encoder.Encode(s.data); err != nil {
		sink.Cancel()
		return err
	}
	return nil
}

// Release is a no-op. Called when Raft is done with the snapshot.
func (s *storeSnapshot) Release() {}
