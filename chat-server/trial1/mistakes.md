Now, there *are* still real issues.

And this is where programming becomes engineering.

---

# Big issue right now

Each goroutine receives its OWN COPY of the slice header:

```go id="jlwmw1"
go handleClient(conn, clients)
```

Meaning:

* old goroutines may not see newly connected clients.

This is subtle and VERY real concurrency/shared-state behavior.

You’ve officially reached the interesting systems part 😭

---

For example:

* Client A connects
  → slice contains `[A]`

* Goroutine A gets that slice snapshot.

Later:

* Client B connects
  → main slice becomes `[A, B]`

But goroutine A may still only know about `[A]`.

That’s because slices are descriptors:

```text id="jlwmw2"
pointer + len + cap
```

and you copied the descriptor into the goroutine call.

That’s a REAL systems concept.

---

# Another issue

Disconnected clients remain forever.

Eventually:

```go id="jlwmw3"
client.Write(...)
```

will fail.

You’ll need cleanup logic later.

Again:
real engineering problem.

---

# Another issue

Multiple goroutines are accessing:

```go id="jlwmw4"
clients
```

without synchronization.

You have now naturally arrived at:

* race conditions,
* shared mutable state,
* synchronization needs.

THIS is where channels/mutexes become meaningful instead of theoretical.

---