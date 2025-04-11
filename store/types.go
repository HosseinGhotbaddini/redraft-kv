package store

// Command represents a Raft-replicated operation applied via the FSM.
type Command struct {
	Op    string // "set", "delete"
	Key   string
	Value []byte
}
