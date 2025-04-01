package server

import (
	"encoding/json"
	"strings"

	"github.com/HosseinGhotbaddini/redraft-kv/raft"
	"github.com/HosseinGhotbaddini/redraft-kv/store"
	"github.com/tidwall/redcon"
)

// HandleCommand interprets a Redis command and routes it to the appropriate backend logic.
// - GET is handled locally from the store (eventually consistent).
// - SET and DELETE are replicated through the Raft log for consistency.
func HandleCommand(nodeID string, r *raft.Node, kv *store.Store, cmd redcon.Command) string {
	args := cmd.Args
	if len(args) == 0 {
		return "ERR empty command"
	}

	switch strings.ToUpper(string(args[0])) {

	case "SET":
		if len(args) != 3 {
			return "ERR wrong number of arguments for SET"
		}

		// Create a command struct and marshal it for Raft replication
		c := store.Command{Op: "set", Key: string(args[1]), Value: args[2]}
		b, _ := json.Marshal(c)

		// Apply through Raft to ensure the operation is replicated and ordered
		if err := r.Apply(b); err != nil {
			return "ERR Raft apply failed"
		}
		return "OK"

	case "GET":
		if len(args) != 2 {
			return "ERR wrong number of arguments for GET"
		}

		// Read from the local store (not replicated, eventual consistency)
		val, ok := kv.Get(string(args[1]))
		if !ok {
			return "(nil)"
		}
		return string(val)

	case "DELETE":
		if len(args) != 2 {
			return "ERR wrong number of arguments for DELETE"
		}

		// Create a delete command and send it through Raft
		c := store.Command{Op: "delete", Key: string(args[1])}
		b, _ := json.Marshal(c)

		if err := r.Apply(b); err != nil {
			return "ERR Raft apply failed"
		}
		return "OK"

	default:
		return "ERR unknown command"
	}
}
