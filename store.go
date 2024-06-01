package flow

import "sync"

type Store[T any] struct {
	mu    *sync.RWMutex
	store map[string]T
}

func NewStore[T any]() *Store[T] {
	return &Store[T]{
		mu:    &sync.RWMutex{},
		store: map[string]T{},
	}
}

func (p *Store[T]) Put(id string, val T) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.store[id] = val
}

func (p *Store[T]) Get(id string) (T, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	val, ok := p.store[id]
	return val, ok
}

func (p *Store[T]) Delete(id string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.store, id)
}
