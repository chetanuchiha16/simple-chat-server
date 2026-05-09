# Mistakes in trial2

`trial2` fixes the slice-sharing bug from `trial1`, but several important issues still remain.

## What is fixed from trial1

- Shared-slice visibility bug is fixed.
  - `handleClient` now receives `clients *[]net.Conn`, so goroutines use one shared slice state.
  - This avoids the `trial1` behavior where goroutines could miss newly appended clients due to copied slice headers.
- Basic synchronization is added with `sync.Mutex`.
  - `append`, broadcast iteration, and connected-count reads are now guarded by `mu`.
  - This is a correct direction for shared mutable state.

## Remaining mistakes

- `Accept()` error handling order is wrong.
  - The code appends `conn` and starts a goroutine before verifying `err`.
  - On accept failure, `conn` can be invalid (`nil`) and should not enter `clients`.
- Holding `mu` while writing to every client is dangerous.
  - If any client is slow or blocked, the broadcast loop stalls and no other goroutines can append or broadcast safely.
  - This is a classic "mutex held while doing I/O" anti-pattern.
- There is no cleanup for closed or dead connections.
  - If a client disconnects, the code never removes it from the slice.
  - Writes to stale connections will continue to fail or block.
- Error handling is still incomplete.
  - `scanner.ReadString('\n')` errors are printed but the loop continues instead of breaking and cleaning up.
  - `client.Write(...)` errors are ignored completely.
- The code still uses `println` and inconsistent logging. A proper `log` package would be better.

## Conceptual issue

- Using bare `[]net.Conn` + one global mutex is okay for learning, but a better design is a client manager with per-client writer goroutines.
- The current implementation is still too close to a toy example rather than a robust server.

## What to fix next

- Check `err` from `Accept()` immediately and `continue` on failure.
- On read/write failure, close the connection, remove it from `clients`, and stop that goroutine.
- Avoid lock-held writes: copy the client slice under lock, unlock, then write; or use per-client outbound channels.
- Replace `println` with structured `log` output and include remote addresses for traceability.
