package server

import (
	"log"

	"github.com/HosseinGhotbaddini/redraft-kv/raft"
	"github.com/HosseinGhotbaddini/redraft-kv/store"

	"github.com/tidwall/redcon"
)

// Start launches a Redis-compatible TCP server using redcon.
// - redisAddr: bind address for Redis (e.g., 127.0.0.1:9001)
// - nodeID: used for command tracing/debugging
// - r: the raft engine to route write commands through
// - kv: the local key-value store (implements store.Store)
func Start(redisAddr string, nodeID string, r *raft.Node, kv store.Store) error {
	log.Printf("Starting Redis server on %s...", redisAddr)
	return redcon.ListenAndServe(redisAddr,
		func(conn redcon.Conn, cmd redcon.Command) {
			// Handle each incoming Redis command
			response := HandleCommand(nodeID, r, kv, cmd)
			conn.WriteString(response)
		},
		nil, // accept hook (not used)
		nil, // close hook (not used)
	)
}
