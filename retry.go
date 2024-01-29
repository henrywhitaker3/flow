package flow

import "context"

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
