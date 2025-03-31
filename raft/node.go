package raft

type Node interface {
	Apply([]byte) error
}

func NewRaftNode(id string, fsm FSM) (Node, error) {
	return &mockRaft{fsm: fsm}, nil // placeholder mock
}

type mockRaft struct {
	fsm FSM
}

func (r *mockRaft) Apply(data []byte) error {
	return r.fsm.Apply(data)
}
