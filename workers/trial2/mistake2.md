These two versions behave VERY differently.

---

# Version 1

```go id="9mjlwm"
go func() {
	defer wg.Done()

	for i, out := range outs {
		for result := range out {
			fmt.Printf("%vst: %v\n", i, result)
		}
	}
}()
```

This is:

```text id="rjlwm8"
ONE goroutine reading outputs SEQUENTIALLY
```

Meaning:

```text id="jlwm2v"
read out0 completely
THEN read out1 completely
THEN read out2 completely
```

That is the important part.

---

# Why this can be bad

Suppose:

```text id="n8jlwm"
worker0 sends result → out0
worker1 sends result → out1
worker2 sends result → out2
```

Reader starts with:

```go id="6jlwmv"
for result := range out0
```

It stays there until `out0` closes.

Meanwhile:

* nobody reads `out1`
* nobody reads `out2`

So workers 1 and 2 may block on:

```go id="jlwm5t"
out <- result
```

even though reader goroutine exists.

Because it is stuck reading only one channel.

---

# Version 2

```go id="jlwm9q"
for i, out := range outs {
	wg.Add(1)

	go func(out <-chan int) {
		defer wg.Done()

		for result := range out {
			fmt.Printf("%vst: %v\n", i, result)
		}
	}(out)
}
```

This creates:

```text id="jlwm3f"
ONE reader goroutine PER output channel
```

So:

```text id="jlwmz1"
reader0 handles out0
reader1 handles out1
reader2 handles out2
```

Now all outputs are consumed concurrently.

Workers won't block waiting for their output to be read.

---

# But there is STILL a bug here ⚠️

This line:

```go id="jlwmx8"
fmt.Printf("%vst: %v\n", i, result)
```

captures outer `i`.

Classic closure bug again.

Fix:

```go id="hjlwm4"
go func(id int, out <-chan int) {
	defer wg.Done()

	for result := range out {
		fmt.Printf("%vst: %v\n", id, result)
	}
}(i, out)
```

---

# Mental Model

Version 1:

```text id="jlwm7c"
reader:
out0 → out1 → out2
```

Version 2:

```text id="jlwm6m"
reader0: out0
reader1: out1
reader2: out2
```

---

# Which is better?

Usually Version 2.

Because fan-out systems are naturally concurrent.

Version 1 accidentally serializes output consumption.
