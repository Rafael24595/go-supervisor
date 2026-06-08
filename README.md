# go-supervisor

A lightweight supervisor for Go applications inspired by fault-tolerant systems.

`go-supervisor` executes tasks under supervision, automatically recovers from panics, captures tack traces, and restarts failed tasks according to configurable restart policies.

The library is designed for long-running processes such as:

- Event loops
- Background workers
- Daemons
- UI runtimes
- Terminal applications
- Service coordinators

---

## Why?

In Go, a panic inside a goroutine usually terminates that goroutine permanently.

While this behavior is often desirable, some components are expected to be resilient and recover automatically from unexpected failures.

`go-supervisor` provides a simple supervision model:

- Recover panics safely
- Capture useful debugging information
- Decide whether a task should be restarted
- Limit restart attempts
- Apply restart delays
- Customize stack trace formatting

---

## Features

- Panic recovery
- Automatic task restart
- Configurable restart policies
- Restart limits
- Restart delays
- Context-aware shutdown
- Stack trace capture
- Custom stack formatting
- Small dependency surface
- Zero goroutine management assumptions

---

## Installation

```sh
go get github.com/Rafael24595/go-supervisor
```

---

## Quick Start

```go
package main

import (
	"errors"

	"github.com/Rafael24595/go-supervisor/supervisor"
)

func main() {
	sup := supervisor.New("worker")

	err := sup.Run(func() error {
		return errors.New("unexpected failure")
	})

	if err != nil {
		panic(err)
	}
}
```

---

## Example: Supervising an Event Loop

```go
policy := policy.New(
    500*time.Millisecond,
    5,
    runtime.DefaultRestartIf,
)

stack := compact.Stack(
    compact.GroupedFormatter,
    compact.WithProcessor(
        compact.StandardProcessor(),
    ),
)

sup := supervisor.New(
    "engine",
    supervisor.WithContext(ctx),
    supervisor.WithPolicy(policy),
    supervisor.WithStackProvider(stack),
)

sup.Run(func() error {
    initialize()

    for {
        processEvents()
    }
})
```

If the event loop panics unexpectedly:

1. The panic is recovered.
2. The stack trace is captured.
3. The restart policy is evaluated.
4. The task is restarted if allowed.

---

## Restart Policies

Policies control when and how a task is restarted.

```go
policy := policy.New(
	time.Second,
	5,
	func(res result.Result) bool {
		return true
	},
)
```

| Field | Description |
|:--|:--|
| Delay | Time to wait before restarting |
| MaxRestarts | Maximum number of restart attempts. A value of `0` means unlimited restarts. |
| RestartIf | Determines whether the task should be restarted |

---

## Panic Recovery

Panics are automatically recovered.

```go
panic("something went wrong")
```

becomes:

```go
result.PanicError
```

with:

- panic value
- stack trace
- original error support through errors.Is and errors.As

```go
var panicErr result.PanicError

if errors.As(err, &panicErr) {
	fmt.Println(panicErr.Panic)
}
```

---

## Custom Stack Formatting

The stack package allows customizing panic output.

Built-in formatters:

### Grouped

```log
ORIGIN
	main.main [main.go:12]

STACK
	service.Run [service.go:42]
	handler.Execute [handler.go:17]
```

### Inline

```log
-> main.main [main.go:12]
   service.Run [service.go:42]
   handler.Execute [handler.go:17]
```

Example:

```go
provider := compact.Stack(
    compact.GroupedFormatter,
)

sup := supervisor.New(
    "worker",
    supervisor.WithStackProvider(provider),
)
```

---

## Custom Logging

Supervisor events are written to an `io.Writer`.

By default, output is written to a decorated `os.Stderr` writer that ensures every log entry is terminated with a newline. This keeps log output readable, even when multiple supervisor events are emitted in quick succession.

A custom writer can be provided using `WithWriter`.

```go
sup := supervisor.New(
    "worker",
    supervisor.WithWriter(
        log.WriterFromCategory("SUPERVISOR"),
    ),
)
```

The writer package provides utilities for decorating output streams.

### Decorator

```go
w := &writer.Decorator{
    writer: os.Stdout,
    prefix: []byte("[engine] "),
    suffix: []byte("\n"),
}
```

The decorator is safe for concurrent use and guarantees that each write operation is emitted as a single formatted message.

---

## Context-Aware Shutdown

A supervisor can be bound to a context.

When the context is cancelled, the supervisor stops restarting the task and exits gracefully.

```go
ctx, cancel := context.WithCancel(context.Background())

sup := supervisor.New(
    "worker",
    supervisor.WithContext(ctx),
)

go func() {
    time.Sleep(time.Second)
    cancel()
}()

err := sup.Run(func() error {
    return errors.New("temporary failure")
})
```

The supervisor will continue applying the configured restart policy until the context is cancelled.

---

## Philosophy

`go-supervisor` focuses on a simple goal:

| Execute tasks safely and decide whether they should be restarted.

The library embraces Go's lightweight concurrency model and aims to provide fault tolerance without introducing unnecessary complexity.

Current features focus on:

- Panic recovery
- Restart policies
- Context-aware shutdown
- Stack trace capture and formatting

The project intentionally starts with a small and focused API. More advanced supervision patterns may be explored in the future as real-world use cases emerge.

### Future Direction

Potential future areas of exploration include:

- Supervision groups
- Child task coordination
- Restart strategies
- Backoff policies
- Hierarchical supervision

Any additions will follow the same design principle:

| Keep the API small, predictable, and easy to integrate into existing Go applications.

---

## License

This project is licensed under the MIT License.

---