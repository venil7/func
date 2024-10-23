package result

import "github.com/venil7/func/function"

type Result[A any] struct {
	data A
	err  error
}

func Ok[A any](data A) Result[A] {
	return Result[A]{data, nil}
}

func Err[A any](err error) Result[A] {
	return Result[A]{data: *new(A), err: err}
}

func From[A any](f function.LazyErr[A]) Result[A] {
	data, err := f()
	if err == nil {
		return Ok(data)
	}
	return Err[A](err)
}

func (r Result[A]) IsOk() bool {
	return r.err == nil
}

func (r Result[A]) IsErr() bool {
	return !r.IsOk()
}

func (r Result[A]) Tuple() (A, error) {
	return r.data, r.err
}

func (r Result[A]) Match(f function.Mapping[A, any]) error {
	if r.IsOk() {
		f(r.data)
		return nil
	}
	return r.err
}

func Map[A any, B any](r *Result[A], f function.Mapping[A, B]) Result[B] {
	if r.IsOk() {
		return Ok(f(r.data))
	}
	return Err[B](r.err)
}

func FlatMap[A any, B any](r *Result[A], f function.Mapping[A, Result[B]]) Result[B] {
	if r.IsOk() {
		return f(r.data)
	}
	return Err[B](r.err)
}

func FlatMapTuple[A any, B any](r *Result[A], f function.LazyErr[B]) Result[B] {
	if r.IsOk() {
		return From(func() (B, error) {
			return f()
		})
	}
	return Err[B](r.err)
}
