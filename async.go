package flow

import "context"

type Result[T any] struct {
	out T
	err error
}

func (r *Result[T]) Err() error {
	return r.err
}

func (r *Result[T]) Out() T {
	return r.out
}

func Eventually[T any](ctx context.Context, f Effector[T]) <-chan Result[T] {
	out := make(chan Result[T], 1)

	go func() {
		res, err := f(ctx)

		out <- Result[T]{
			out: res,
			err: err,
		}
	}()

	return out
}
