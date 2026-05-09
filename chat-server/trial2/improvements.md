# Improvements in trial2

This version fixes the main shared-slice bug from `trial1` and makes state sharing explicit.

## What got better

- Switched from passing `clients []net.Conn` by value to passing `&clients` into `handleClient`.
- Added `sync.Mutex` and used it around `append(clients, conn)`.
- Added locking around iteration of `clients` during broadcast.
- Added locking around `len(*clients)` when printing connected count.

## Why this is an improvement

- In `trial1`, each goroutine received a copy of the slice header, so later `append` operations in `main` were not reliably visible to older goroutines.
- In `trial2`, all goroutines read the same shared slice through `*[]net.Conn`, so they can observe newly accepted clients.
- The mutex gives a single synchronization point for shared-slice reads/writes, which removes obvious concurrent slice corruption.

## Still not production-ready

- Accept error handling is in the wrong order: failed `Accept()` should be checked before appending/spawning.
- Dead clients are never removed from `clients`, so stale connections accumulate.
- Socket writes happen while holding `mu`, so one slow client can block all broadcasting and client-list operations.
- `ReadString`/`Write` errors are printed but loops continue forever instead of disconnecting the client cleanly.
- Logging style is inconsistent (`fmt.Println` + `println`) and misses useful context.

## Good direction

- `trial2` is a meaningful correctness step from `trial1`.
- Next milestone: lifecycle cleanup, better error exits, and avoiding lock-held network I/O.
