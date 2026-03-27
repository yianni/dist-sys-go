# TLA+ Specs

These specs model the challenge semantics for selected Gossip Glomers workloads.

They are intentionally abstract:
- they model logical state and invariants
- they do not model Go package structure
- they do not model Maelstrom transport details beyond what is needed for correctness

## Included Specs

- `gcounter/GCounter.tla`
- `broadcast/Broadcast.tla`
- `kafka/Kafka.tla`
- `txn/TxnLWW.tla`

## Modeling Conventions

- finite node, key, and value sets
- safety properties first
- liveness and fairness only where needed
- reads and polls modeled as derived views when they do not affect state
- message passing modeled abstractly using inflight message state

## Non-Goals

- proving implementation equivalence with the Go code
- modeling JSON or RPC encoding
- modeling goroutines or Maelstrom internals
- modeling production durability or fault tolerance beyond the challenge semantics

## Running

Use the TLA+ Toolbox or TLC CLI with the corresponding `.cfg` file for each module.

Start with small model sizes.
