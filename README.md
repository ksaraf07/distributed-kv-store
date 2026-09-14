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
- [ ] Stage 3 — replication
- [ ] Stage 4 — leader election (stretch)
- [ ] Stage 5 — evaluation / load testing

## Why Go

Go's concurrency primitives (goroutines, channels) map directly onto the
core problems a networked, concurrent store has to solve, and it's the
language behind real systems in this space (etcd, Docker, Kubernetes).

## Running

\`\`\`bash
go run main.go
\`\`\`