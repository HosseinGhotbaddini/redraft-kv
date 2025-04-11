# Redraft

Redraft is a distributed key-value store that uses the Redis wire protocol for client interaction and Raft consensus for state replication. It is implemented in Go and designed with a modular, extensible architecture.

---

## Features

- Redis-compatible interface via [`redcon`](https://github.com/tidwall/redcon)
- Raft-based state replication using [`hashicorp/raft`](https://github.com/hashicorp/raft)
- Dynamic cluster membership via `BOOTSTRAP` and `JOIN` commands
- Introspection commands: `LEADER`, `NODES`
- FSM snapshot and restore support for log compaction and faster recovery
- Modular architecture: server / raft / store
- In-memory or BoltDB-based key-value store
- Configurable via YAML
- Deterministic FSM with basic operations (`SET`, `GET`, `DELETE`)
- Unit test coverage and dev-friendly Makefile

---

## Architecture

Redraft is composed of three core modules:

- **server/** — Redis protocol server and command routing
- **raft/** — Raft node lifecycle and FSM integration
- **store/** — Pluggable key-value store (in-memory or BoltDB)

See [docs/DESIGN.md](docs/DESIGN.md) for detailed architecture notes.

---

## Getting Started

### 1. Install Dependencies

```bash
go mod tidy
```

---

### 2. Prepare Config Files

Each node config YAML should define only its own local addresses and store backend.

```yaml
# config/node1.yaml
id: node1
raft_addr: 127.0.0.1:7001
redis_addr: 127.0.0.1:9001
store_backend: bolt
store_path: data/node1/store.db
```

Repeat for `node2.yaml`, `node3.yaml`, etc.

---

### 3. Run Nodes

```bash
make node1
make node2
make node3
```

Or start all three in a tmux session:

```bash
make cluster
```

---

### 4. Form the Cluster

Connect to the leader node (e.g., `redis-cli -p 9001`), then:

```bash
BOOTSTRAP node1 127.0.0.1:7001
JOIN node2 127.0.0.1:7002
JOIN node3 127.0.0.1:7003
```

---

### 5. Test Key Commands

```bash
SET foo bar
GET foo
DELETE foo
```

Only the leader accepts writes. Reads (`GET`) work on all nodes.

---

### 6. Cluster Introspection

```bash
LEADER     # returns "node1 127.0.0.1:7001"
NODES      # returns all nodes, one per line
SNAPSHOT   # triggers a manual FSM snapshot (dev tool)
```

---

## Development

### Run All Tests

```bash
make test
```

### Clean Build Artifacts

```bash
make clean
```

Removes the binary and data directories.

---

## Roadmap

See [docs/ROADMAP.md](docs/ROADMAP.md) for planned features and design goals.