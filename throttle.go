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
