package flow

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	ErrThrottled = errors.New("throttled")
)

func ThrottleIn[T any](f EffectorIn[T], every time.Duration) EffectorIn[T] {
	called := false
	once := &sync.Once{}

	return func(ctx context.Context, t T) error {
		once.Do(func() {
			ticker := time.NewTicker(every)
			go func() {
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return
					case <-ticker.C:
						called = false
					}
				}
			}()
		})

		if !called {
			called = true
			return f(ctx, t)
		}

		return ErrThrottled
	}
}

func ThrottleInSilently[T any](f EffectorIn[T], every time.Duration) EffectorIn[T] {
	tf := ThrottleIn[T](f, every)
	return func(ctx context.Context, t T) error {
		err := tf(ctx, t)
		if errors.Is(err, ErrThrottled) {
			return nil
		}
		return err
	}
}

// Return an Effector that, when called, will only fire once every duration
func Throttle[T any](f Effector[T], every time.Duration) Effector[T] {
	called := false
	once := sync.Once{}

	return func(ctx context.Context) (T, error) {
		// Run a loop here to turn called false after every x duration
		once.Do(func() {
			ticker := time.NewTicker(every)
			go func() {
				defer ticker.Stop()
				for {
					select {
					case <-ctx.Done():
						return

					case <-ticker.C:
						called = false
					}
				}
			}()
		})

		if !called {
			called = true
			return f(ctx)
		}

		var out T
		return out, ErrThrottled
	}
}

func SilentThrottle[T any](f Effector[T], every time.Duration) Effector[T] {
	var value T
	var err error

	tf := Throttle(f, every)
	return func(ctx context.Context) (T, error) {
		res, iErr := tf(ctx)
		if errors.Is(iErr, ErrThrottled) {
			return value, nil
		}

		value = res
		err = iErr

		return value, err
	}
}
