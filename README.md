# Func

The `task` package provides an abstraction for handling computations that might fail. It builds on the functional programming principles using lazy evaluations and error handling, offering an expressive way to compose functions and manage asynchronous tasks in Go.

## Installation

To install, run:

```bash
go get github.com/venil7/func
```

## Overview

The `task` library is based on a `Task` type, which represents a computation that can succeed or fail. By composing functions, users can handle both synchronous and asynchronous computations, and build complex workflows while preserving error handling.

### Types

- **Task[A any]**: The core type, representing a lazy computation that either returns a result of type `A` or an error.

### Constructors

- **Of**: Wraps a value in a `Task`.
  ```go
  func Of[A any](a A) Task[A]
  ```

- **Fail**: Creates a `Task` that fails with the given error.
  ```go
  func Fail[A any](err error) Task[A]
  ```

- **From**: Converts a `LazyErr` function to a `Task`.
  ```go
  func From[A any](f function.LazyErr[A]) Task[A]
  ```

- **From1**: Converts a single-argument `MapLazyErr` function to a function returning a `Task`.
  ```go
  func From1[A, B any](f function.MapLazyErr[A, B]) function.Mapping[A, Task[B]]
  ```

- **From2**: Converts a two-argument function to one that returns a `Task`.
  ```go
  func From2[A, B, C any](f func(a A, b B) (C, error)) func(a A, b B) Task[C]
  ```

### Operations on Tasks

- **Map**: Applies a mapping function to a `Task`'s result, returning a new `Task`.
  ```go
  func Map[A any, B any](t Task[A], f function.Mapping[A, B]) Task[B]
  ```

- **FlatMap**: Chains `Task`s by applying a mapping function that itself returns a `Task`.
  ```go
  func FlatMap[A any, B any](t Task[A], f function.Mapping[A, Task[B]]) Task[B]
  ```

- **Tap**: Executes a `Task` for its side effects, discarding the result but preserving the original `Task`.
  ```go
  func Tap[A, B any](t Task[A], f function.Mapping[A, Task[B]]) Task[A]
  ```

- **Then**: Combines a `Task` with a `MapLazyErr` function, creating a sequence of dependent tasks.
  ```go
  func Then[A any, B any](t Task[A], f function.MapLazyErr[A, B]) Task[B]
  ```

### Aggregation Functions

- **Sequence**: Takes multiple `Task`s and returns a `Task` containing a slice of results, or an error if any `Task` fails.
  ```go
  func Sequence[A any](ts ...Task[A]) Task[[]A]
  ```

- **Traverse**: Applies a function to each element in a list, transforming them into `Task`s, and collects the results in a single `Task`.
  ```go
  func Traverse[A, B any](ts []A, f function.Mapping[A, Task[B]]) Task[[]B]
  ```

### Utility Functions

- **ToResult**: Converts a `Task` to a `Result`, representing either success or failure.
  ```go
  func (t Task[A]) ToResult() result.Result[A]
  ```

## Example Usage

```go
package main

import (
	"fmt"
	"github.com/venil7/func/task"
)

func main() {
	// Creating a successful task
	t1 := task.Of(42)

	// Creating a failing task
	t2 := task.Fail[int](fmt.Errorf("some error"))

	// Mapping over a task
	t3 := task.Map(t1, func(i int) string { return fmt.Sprintf("Result: %d", i) })

	// Running a sequence of tasks
	result, err := t3()
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Result:", result)
	}
}
```

## License

This package is licensed under the MIT License. See `LICENSE` for details.