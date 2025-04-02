# Redraft

**Redraft** is a distributed key-value store that uses the Redis wire protocol for client interaction and Raft consensus for replication. It is implemented in Go and designed with modular clarity in mind.

---

## Features

- Redis-compatible interface via [`redcon`](https://github.com/tidwall/redcon)
- Raft-based state replication using [`hashicorp/raft`](https://github.com/hashicorp/raft)
- Modular architecture
- Static 3-node cluster bootstrapped
- Deterministic FSM with basic operations
- In-memory key-value store
- Configurable via YAML
- Unit test scaffolding and automated test target
- Developer Makefile for build/run/test


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

Using the Makefile:

```bash
make node1
make node2
make node3
```

Or to start all three in a tmux session:

```bash
make cluster
```

You can also run manually with:

```bash
make run CONFIG=config/custom.yaml
```

---

### 4. Interact via redis-cli

```bash
redis-cli -p 9001
```

Or use ports 9002, 9003 to interact with other nodes.

Then test it:

```bash
SET foo bar
GET foo
DEL foo
```

Only the leader node will accept write commands. Use GET on any node.

---

## Run Tests

```bash
make test
```

This runs all unit tests in test/.

---

## Clean Build Artifacts

```bash
make clean
```

Removes the binary and local data/ directories used by Raft.

---

## Notes

For future improvements, see [docs/ROADMAP.md](docs/ROADMAP.md).