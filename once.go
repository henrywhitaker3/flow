package flow

import "sync"

type OncePer[T comparable] struct {
	onces map[T]*sync.Once
}

func NewOncePer[T comparable]() *OncePer[T] {
	return &OncePer[T]{
		onces: make(map[T]*sync.Once),
	}
}

func (o *OncePer[T]) Do(key T, f func()) {
	once, ok := o.onces[key]
	if !ok {
		once = &sync.Once{}
		o.onces[key] = once
	}

	once.Do(f)
}

func (o *OncePer[T]) Reset(key T) {
	delete(o.onces, key)
}
