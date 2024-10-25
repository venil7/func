package task

import (
	"sync"

	"github.com/venil7/func/function"
	"github.com/venil7/func/result"
)

type Task[A any] function.LazyErr[A]

func Of[A any](a A) Task[A] {
	return (func() (A, error) {
		return a, nil
	})
}

func Fail[A any](err error) Task[A] {
	return (func() (A, error) {
		var a A
		return a, err
	})
}

func From[A any](f function.LazyErr[A]) Task[A] {
	return (func() (A, error) {
		return f()
	})
}

func From1[A, B any](f function.MapLazyErr[A, B]) function.Mapping[A, Task[B]] {
	return (func(a A) Task[B] {
		return (func() (B, error) {
			return f(a)
		})
	})
}

func From2[A, B, C any](f func(a A, b B) (C, error)) func(a A, b B) Task[C] {
	return (func(a A, b B) Task[C] {
		return (func() (C, error) {
			return f(a, b)
		})
	})
}

func (t Task[A]) ToResult() result.Result[A] {
	return result.From[A](func() (A, error) {
		return t()
	})
}

func Map[A any, B any](t Task[A], f function.Mapping[A, B]) Task[B] {
	return (func() (B, error) {
		a, err := t()
		if err != nil {
			return *new(B), err
		}
		return f(a), nil
	})
}

func FlatMap[A any, B any](t Task[A], f function.Mapping[A, Task[B]]) Task[B] {
	return (func() (B, error) {
		a, err := t()
		if err != nil {
			return *new(B), err
		}
		return f(a)()
	})
}

func Tap[A, B any](t Task[A], f function.Mapping[A, Task[B]]) Task[A] {
	return func() (A, error) {
		a, err := t()
		if err != nil {
			return a, err
		}
		_, err = f(a)()
		return a, err
	}
}

func Then[A any, B any](t Task[A], f function.MapLazyErr[A, B]) Task[B] {
	return FlatMap(t, func(a A) Task[B] {
		return func() (B, error) { return f(a) }
	})
}

func Sequence[A any](ts ...Task[A]) Task[[]A] {
	return (func() ([]A, error) {
		results := make([]A, len(ts))
		errchan := make(chan error)
		var wg sync.WaitGroup
		for i, t := range ts {
			wg.Add(1)
			go func(i int) {
				a, err := t()
				if err != nil {
					errchan <- err
				}
				results[i] = a
				wg.Done()
			}(i)
		}

		go func() {
			wg.Wait()
			close(errchan)
		}()

		err := <-errchan
		if err != nil {
			return nil, err
		}

		return results, nil
	})
}

func Traverse[A, B any](ts []A, f function.Mapping[A, Task[B]]) Task[[]B] {
	return (func() ([]B, error) {
		tasks := make([]Task[B], len(ts))
		for i, t := range ts {
			tasks[i] = f(t)
		}
		return Sequence(tasks...)()
	})
}
