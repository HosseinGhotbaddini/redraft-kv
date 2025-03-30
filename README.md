# Redraft

**Redraft** is a minimal, Redis-compatible distributed key-value store written in Go. It uses `redcon` for Redis protocol handling and `hashicorp/raft` for replicated state via Raft consensus.

This project is structured to demonstrate clear modular architecture, deterministic state transitions, and consensus correctness in a simplified cluster setting.

---

## Project Status

Implementation in progress.

This README will be updated post-implementation with:
- Run and usage instructions
- Test details
- Final feature list

---

## Architecture

Redraft is organized into three core modules:

- `server/` — Redis protocol server and command router
- `raft/` — Raft node management and FSM
- `store/` — In-memory key-value store

For detailed architecture, design decisions, and rationale, see [docs/DESIGN.md](docs/DESIGN.md).

---

## Planned Features

- `SET`, `GET`, `DELETE` support via Redis client
- Raft-based replication across a 3-node cluster
- Leader-only write path
- In-memory state storage
- Minimal test coverage

