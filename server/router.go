package server

import (
	"encoding/json"
	"strings"

	"github.com/HosseinGhotbaddini/redraft-kv/raft"
	"github.com/HosseinGhotbaddini/redraft-kv/store"
	"github.com/tidwall/redcon"
)

// HandleCommand interprets a Redis command and routes it to the appropriate backend logic.
func HandleCommand(nodeID string, r *raft.Node, kv store.Store, cmd redcon.Command) string {
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
		val, err := kv.Get(string(args[1]))
		if err != nil || val == nil {
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

	case "BOOTSTRAP":
		if len(args) != 3 {
			return "ERR usage: BOOTSTRAP <node-id> <raft-addr>"
		}
		id := string(args[1])
		addr := string(args[2])
		if err := r.BootstrapSelf(id, addr); err != nil {
			return "ERR bootstrap failed: " + err.Error()
		}
		return "OK"

	case "JOIN":
		if len(args) != 3 {
			return "ERR usage: JOIN <node-id> <raft-addr>"
		}
		id := string(args[1])
		addr := string(args[2])
		if err := r.JoinNode(id, addr); err != nil {
			return "ERR join failed: " + err.Error()
		}
		return "OK"

	case "LEADER":
		id, addr := r.GetLeader()
		if id == "" || addr == "" {
			return "(nil)"
		}
		return id + " " + addr

	case "NODES":
		var lines []string
		for _, peer := range r.ListPeers() {
			lines = append(lines, string(peer.ID)+" "+string(peer.Address))
		}
		return strings.Join(lines, "\n")

	case "SNAPSHOT":
		if err := r.Snapshot(); err != nil {
			return "ERR snapshot failed: " + err.Error()
		}
		return "OK"

	default:
		return "ERR unknown command"
	}
}
