# Go Quests 🏹

A hands-on, iterative exploration of **Go concurrency patterns** — built from scratch, broken deliberately, and fixed incrementally. Each project progresses through multiple "trials" that introduce real concurrency bugs and then resolve them, documenting every mistake along the way.

## Projects

### 🗨️ Chat Server — `chat-server/`

A TCP chat server that accepts multiple clients and broadcasts messages between them.
Built to learn goroutines, shared mutable state, and synchronization primitives.

| Trial | Focus | Key Concepts |
|-------|-------|--------------|
| **trial1** | Naive implementation — get it working first | `net.Listen`, `goroutines`, `bufio.Reader`, slice-based client tracking |
| **trial2** | Fix concurrency bugs from trial1 | `sync.Mutex`, passing pointers vs. values, guarded shared state |

#### Trial Progression

**trial1** — A bare-bones TCP server where each client connection spawns a goroutine. Messages are broadcast to all connected clients. This version has intentional concurrency issues:

- Each goroutine receives a **copy of the slice header**, so older goroutines can't see newly connected clients.
- No synchronization — multiple goroutines read/write the `clients` slice concurrently (race condition).
- Disconnected clients are never cleaned up.

**trial2** — Fixes the shared-slice visibility bug by passing `*[]net.Conn` (pointer to slice) and introduces a `sync.Mutex` to guard all access. Remaining issues (documented in `mistakes.md`) include lock-held I/O, missing error handling, and no client cleanup.

#### Running

```bash
# Start the server
cd chat-server/trial2
go run main.go

# Connect clients (in separate terminals)
nc localhost 8080
```

---

### ⚙️ Worker Pool — `workers/`

An implementation of the **fan-out / fan-in** concurrency pattern — distributing tasks across multiple worker goroutines and collecting results.

| Trial | Focus | Key Concepts |
|-------|-------|--------------|
| **trial1** | Basic fan-out with WaitGroup | Spawning N workers, channel distribution, `sync.WaitGroup` |
| **trial2** | Fan-out with per-worker output channels | Output channels, goroutine scheduling, deadlock debugging |
| **trial3** | Adding fan-in to merge results | `fanIn` function, channel merging, pipeline composition |

#### Trial Progression

**trial1** — Workers consume tasks from a shared channel and print results. Uses `sync.WaitGroup` for synchronization. No output channels yet — results are printed directly.

**trial2** — Each worker gets its own output channel. Demonstrates a **classic deadlock**: if the task sender runs on the main goroutine, it blocks before the output reader goroutine can start. Fixed by moving task-sending into a goroutine. Also documents the difference between sequential vs. concurrent output consumption (one reader goroutine vs. one-per-channel).

**trial3** — Introduces a `fanIn` function that merges all worker output channels into a single result channel. This completes the full fan-out/fan-in pipeline pattern.

#### Running

```bash
cd workers/trial3
go run main.go
```

---

## Repository Structure

```
go-quests/
├── chat-server/
│   ├── trial1/
│   │   ├── main.go          # Naive TCP chat server
│   │   ├── mistakes.md      # Concurrency bugs documented
│   │   └── go.mod
│   └── trial2/
│       ├── main.go          # Mutex-protected version
│       ├── mistakes.md      # Remaining issues
│       └── improvements.md  # What changed and why
├── workers/
│   ├── trial1/
│   │   └── main.go          # Basic fan-out
│   ├── trial2/
│   │   ├── main.go          # Fan-out with output channels
│   │   ├── mistake.md       # Deadlock walkthrough
│   │   └── mistake2.md      # Sequential vs. concurrent readers
│   └── trial3/
│       └── main.go          # Fan-out + fan-in
└── README.md
```

## Concepts Covered

- **Goroutines** — lightweight concurrent functions
- **Channels** — typed communication between goroutines
- **`sync.Mutex`** — protecting shared mutable state
- **`sync.WaitGroup`** — waiting for goroutine completion
- **Fan-out** — distributing work across multiple workers
- **Fan-in** — merging multiple channels into one
- **Deadlocks** — understanding channel blocking and execution order
- **Slice internals** — how Go slice headers (pointer + len + cap) affect concurrency
- **Closure bugs** — capturing loop variables in goroutines

## Prerequisites

- **Go 1.22+** (uses `range` over integers: `for i := range N`)

## License

This is a personal learning repository. Feel free to reference it for your own concurrency journey.
