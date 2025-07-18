package test

import (
	"encoding/json"
	"testing"

	"github.com/HosseinGhotbaddini/redraft-kv/raft"
	"github.com/HosseinGhotbaddini/redraft-kv/store"
	hraft "github.com/hashicorp/raft"
)

// helper: marshal store.Command into a mock raft.Log entry
func makeLog(cmd store.Command) *hraft.Log {
	data, _ := json.Marshal(cmd)
	return &hraft.Log{Data: data}
}

func TestFSM_Apply_SetAndGet(t *testing.T) {
	kv := store.NewMemoryStore()
	fsm := &raft.FSMImpl{Store: kv}

	cmd := store.Command{
		Op:    "set",
		Key:   "foo",
		Value: []byte("bar"),
	}

	// Apply log entry to FSM
	result := fsm.Apply(makeLog(cmd))
	if err, ok := result.(error); ok && err != nil {
		t.Fatalf("Apply returned error: %v", err)
	}

	// Verify store state
	val, err := kv.Get("foo")
	if err != nil || val == nil || string(val) != "bar" {
		t.Errorf("expected 'bar', got '%s'", val)
	}
}

func TestFSM_Apply_Delete(t *testing.T) {
	kv := store.NewMemoryStore()
	kv.Set("delete-me", []byte("bye"))
	fsm := &raft.FSMImpl{Store: kv}

	cmd := store.Command{
		Op:  "delete",
		Key: "delete-me",
	}

	result := fsm.Apply(makeLog(cmd))
	if err, ok := result.(error); ok && err != nil {
		t.Fatalf("Apply returned error: %v", err)
	}

	val, _ := kv.Get("delete-me")
	if val != nil {
		t.Error("expected key to be deleted")
	}
}

func TestFSM_Apply_UnknownOperation(t *testing.T) {
	kv := store.NewMemoryStore()
	fsm := &raft.FSMImpl{Store: kv}

	cmd := store.Command{
		Op:    "noop", // unsupported operation
		Key:   "x",
		Value: []byte("y"),
	}

	result := fsm.Apply(makeLog(cmd))
	if err, ok := result.(error); ok && err != nil {
		t.Errorf("expected no error on unknown op, got: %v", err)
	}

	val, _ := kv.Get("x")
	if val != nil {
		t.Error("expected no key to be created for unknown op")
	}
}

func TestFSM_Apply_MalformedData(t *testing.T) {
	kv := store.NewMemoryStore()
	fsm := &raft.FSMImpl{Store: kv}

	// Malformed JSON (not a store.Command)
	badLog := &hraft.Log{
		Data: []byte("{bad json"),
	}

	result := fsm.Apply(badLog)
	if result == nil {
		t.Fatal("expected error for malformed data")
	}
}
