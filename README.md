# Redraft

**Redraft** is a distributed key-value store that uses the Redis wire protocol for client interaction and Raft consensus for replication. It is implemented in Go and designed with modular clarity in mind.

---

## Features

- Redis-compatible interface via [`redcon`](https://github.com/tidwall/redcon)
- Raft-based state replication using [`hashicorp/raft`](https://github.com/hashicorp/raft)
- Configurable via YAML (no CLI flag jungle)
- Modular architecture: Server / Raft / Store
- Deterministic state transitions through FSM
- In-memory key-value store

---

## Architecture

Redraft is composed of three core modules:

- **server/** — Redis protocol server and command routing
- **raft/** — Raft node lifecycle and FSM integration
- **store/** — Thread-safe in-memory key-value state

For a full system overview, see [docs/DESIGN.md](docs/DESIGN.md).

---

## Getting Started

### 1. Install Dependencies

```bash
go mod tidy
```

---

### 2. Prepare Config Files

Create one config file per node under `config/`:

```yaml
# config/node1.yaml
id: node1
raft_addr: 127.0.0.1:7001
redis_addr: 127.0.0.1:9001
peers:
  node2: 127.0.0.1:7002
  node3: 127.0.0.1:7003
```

Duplicate this for node2.yaml, node3.yaml with appropriate IDs and ports.

---

### 3. Run Nodes

In three separate terminals:

```bash
go run main.go config/node1.yaml
go run main.go config/node2.yaml
go run main.go config/node3.yaml
```

---

### 4. Interact via redis-cli

```bash
redis-cli -p 9001
```

Then test it:

```bash
SET foo bar
GET foo
DEL foo
```

Only the leader node will accept write commands. Use GET on any node.

---

## Development Notes

- All node runtime state is written to `data/<nodeID>/` (ignored in Git)
- Each node is bootstrapped using a static cluster config
- Snapshotting is not yet implemented (stub provided)
- Cluster membership is static (no dynamic join)

---

## Next Steps

- Add integration tests
