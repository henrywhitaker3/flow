package flow

import "context"

type Effector[T any] func(context.Context) (T, error)
