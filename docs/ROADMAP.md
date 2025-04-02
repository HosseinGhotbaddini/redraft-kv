# ROADMAP.md

## Redraft: Roadmap and Future Enhancements

This document outlines potential future improvements and architectural extensions, assuming additional time or team development. Prioritization is based on technical value, user experience, and system completeness.

---

## Core Improvements

### **State Persistence**
- Swap in-memory store with disk-backed store (e.g., BoltDB, Pebble)
- Enable safe restarts and crash recovery
- Separate persistence logic behind a clean interface

### **Snapshotting**
- Implement `Snapshot()` and `Restore()` methods in FSM
- Avoid unbounded log growth
- Reduce replay time during recovery

### **Leader Discovery**
- Add a `LEADER` Redis command to return the current leader's ID and address
- Improve UX by guiding writes to the correct node

---

## Dev & Ops Enhancements

### **Dynamic Cluster Membership**
- Support joining and removing nodes at runtime via Raft APIs
- Replace static config model with join tokens or discovery

### **Observability**
- Expose metrics (Prometheus-compatible) for:
  - Raft state (leader/follower)
  - Command throughput
  - Replication lag
- Add basic structured logging and trace identifiers

---

## Interface Features

### **Additional Redis Commands**
- `SETMANY`, `GETMANY` for multi-key batch operations
- `EXISTS`, `TTL`, or key metadata extensions
- Improve command dispatching to allow feature routing

### **Alternative Interfaces**
- Optional REST or gRPC gateway for non-Redis clients
- JSON-based API to mirror Redis functionality

---

## Testing & Dev Experience

### **End-to-End Cluster Tests**
- Add integration test harness for 3-node cluster
- Verify consistency and log application across replicas

### **Client Compatibility Testing**
- Test with popular Redis clients (go-redis, Jedis, redis-py)
- Ensure protocol compliance

---

## Growth Areas

### **Redraft as a Library**
Package Redraft as a reusable Go module to embed in other distributed systems or dev tools. This makes it easier to adopt Redraft’s key-value engine in broader infrastructure workflows.

### **Snapshotting & Time-Travel**
Implement Raft FSM snapshots with support for tagged versions and rollback. Combine with a replay tool to inspect or restore historical state transitions for debugging or observability purposes.

### **Verifiable State Proofs**
Use Merkle trees to let clients verify that a key exists in a given state. Tag snapshots by their Merkle root to support auditability, and trustless validation.

### **Kubernetes Deployment Support**
Add Helm chart and/or Kubernetes Operator to enable Redraft to run in modern infrastructure stacks. Automate node bootstrapping, recovery, and Raft peer discovery.


---

## Summary

This roadmap reflects Redraft’s potential evolution from a minimal demo into a robust, distributed, and extensible key-value system. Each milestone is intended to deepen real-world capability while preserving the architectural simplicity and modularity established in the initial implementation.

