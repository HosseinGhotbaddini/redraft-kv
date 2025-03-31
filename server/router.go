package server

import (
	"encoding/json"
	"strings"

	"github.com/HosseinGhotbaddini/redraft-kv/raft"
	"github.com/HosseinGhotbaddini/redraft-kv/store"
	"github.com/tidwall/redcon"
)

func HandleCommand(nodeID string, r raft.Node, kv *store.Store, cmd redcon.Command) string {
	args := cmd.Args
	if len(args) == 0 {
		return "ERR empty command"
	}

	switch strings.ToUpper(string(args[0])) {
	case "SET":
		if len(args) != 3 {
			return "ERR wrong number of arguments for SET"
		}
		c := store.Command{Op: "set", Key: string(args[1]), Value: args[2]}
		b, _ := json.Marshal(c)
		if err := r.Apply(b); err != nil {
			return "ERR Raft apply failed"
		}
		return "OK"

	case "GET":
		if len(args) != 2 {
			return "ERR wrong number of arguments for GET"
		}
		val, ok := kv.Get(string(args[1]))
		if !ok {
			return "(nil)"
		}
		return string(val)

	case "DELETE":
		if len(args) != 2 {
			return "ERR wrong number of arguments for DELETE"
		}
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
