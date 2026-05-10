THAT goroutine is the reader for `out`.

The problem is:

it starts **too late** because main is already blocked before reaching a working state.

Let’s walk the timeline carefully.

---

# Your code order

```go id="4m8u6v"
tasks := make(chan int)

outs := fanOut(tasks, 3)

for i := 0; i < 10; i++ {
    tasks <- i
}
```

STOP here.

Main is now sending tasks.

---

Meanwhile workers do:

```go id="c4tr4e"
for task := range tasks {
    result := task * 2
    out <- result
}
```

Suppose:

* worker receives task 0
* computes result
* tries:

```go id="zwzjlwm"
out <- result
```

But who is reading `out` right now?

Not this yet:

```go id="fjlwmr"
go func() {
    for _, out := range outs {
```

because main has NOT REACHED THAT CODE YET.

Main is still stuck in:

```go id="o6syqp"
tasks <- i
```

---

# Why main gets stuck

Worker blocks on:

```go id="bj6m5m"
out <- result
```

So worker stops reading from `tasks`.

Then eventually ALL workers block on output sends.

Now nobody receives from `tasks`.

Then main blocks sending to `tasks`.

So execution never reaches:

```go id="kzjlwm"
go func() {
```

That reader goroutine never starts.

---

# Visual timeline

```text id="thj1uq"
main:
tasks <- 0
tasks <- 1
tasks <- 2

workers:
receive task
try out <- result
BLOCK

main:
tasks <- 3
BLOCK

reader goroutine:
never started yet
```

---

# Why putting sender in goroutine fixes it

When you do:

```go id="zjlwmn"
go func() {
    for i := 0; i < 10; i++ {
        tasks <- i
    }
}()
```

main is no longer blocked sending tasks.

So it continues downward and starts:

```go id="a1jlwm"
go func() {
    for _, out := range outs {
```

NOW receivers exist for `out`.

So workers can send results.

So workers continue reading tasks.

So sender continues sending tasks.

Pipeline flows correctly.

---

This exact thing is one of the hardest beginner concurrency concepts:

```text id="c24wwh"
"code order" != "execution order"
```

Goroutines run independently.

Blocking changes the entire execution graph.
