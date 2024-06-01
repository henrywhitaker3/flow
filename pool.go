package flow

import "sync"

type Pool[T any] struct {
	mu   *sync.RWMutex
	pool map[string]T
}

func NewPool[T any]() *Pool[T] {
	return &Pool[T]{
		mu:   &sync.RWMutex{},
		pool: map[string]T{},
	}
}

func (p *Pool[T]) Put(id string, val T) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.pool[id] = val
}

func (p *Pool[T]) Get(id string) (T, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	val, ok := p.pool[id]
	return val, ok
}
