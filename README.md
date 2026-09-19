# distributed-kv-store

A distributed, replicated key-value store built from scratch in Go, as a
learning project to understand the fundamentals behind systems like Redis
and etcd: networking, persistence, replication, and fault tolerance.

## Status

🚧 In progress — building incrementally, stage by stage.

- [x] Stage 1.1 — project skeleton
- [x] Stage 1.2 — in-memory store
- [x] Stage 1.3 — TCP server
- [x] Stage 1.4 — CLI client
- [x] Stage 2 — persistence (write-ahead log)
- [x] Stage 3 — replication

- [x] Stage 4 — leader election (stretch)

## Known limitations

- Only the leader is *conventionally* expected to accept writes — nothing in
  the code currently enforces this, so a client could technically SET
  directly on a follower.
- With only 2 nodes, "election" is really just "the one other node
  promotes itself" — no voting/quorum logic yet.
- Clients aren't automatically redirected to a newly promoted leader;
  they'd need to know to reconnect to the new address.


- [x] Stage 5 — evaluation / load testing

## Benchmark Results

Tested with 20 concurrent clients, 100 requests each (2,000 total), on a single machine.

| Configuration | Throughput | Total time |
|---|---|---|
| No replication | ~64,450 req/sec | 31.0ms |
| 1 follower (replicated) | ~40,900 req/sec | 48.9ms |

Replication introduces roughly 37% throughput overhead — the cost of the
leader synchronously forwarding each write to its follower before
returning a response, in exchange for data redundancy.

## Why Go

Go's concurrency primitives (goroutines, channels) map directly onto the
core problems a networked, concurrent store has to solve, and it's the
language behind real systems in this space (etcd, Docker, Kubernetes).

## Running

\`\`\`bash
go run main.go
\`\`\`