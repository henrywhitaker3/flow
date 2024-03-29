package flow

import (
	"context"
	"time"
)

// Retries calling the effector x times if it fails
func Retry[T any](f Effector[T], times int) Effector[T] {
	return func(ctx context.Context) (T, error) {
		var out T
		var err error
		for i := 0; i < times; i++ {
			out, err = f(ctx)
			if err == nil {
				return out, nil
			}
		}
		return out, err
	}
}

// Does the same a Retry, but adds a delay after each attempt
func RetryDelay[T any](f Effector[T], times int, delay time.Duration) Effector[T] {
	return func(ctx context.Context) (T, error) {
		var out T
		var err error
		for i := 0; i < times; i++ {
			out, err = f(ctx)
			if err == nil {
				return out, nil
			}
			time.Sleep(delay)
		}
		return out, err
	}
}
