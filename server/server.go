package server

import (
	"log"

	"github.com/HosseinGhotbaddini/redraft-kv/raft"
	"github.com/HosseinGhotbaddini/redraft-kv/store"

	"github.com/tidwall/redcon"
)

func Start(nodeID string, r raft.Node, kv *store.Store) error {
	addr := ":9001" // static for now
	log.Printf("Starting Redis server on %s...", addr)
	return redcon.ListenAndServe(addr,
		func(conn redcon.Conn, cmd redcon.Command) {
			response := HandleCommand(nodeID, r, kv, cmd)
			conn.WriteString(response)
		},
		nil, // accept
		nil, // close
	)
}
