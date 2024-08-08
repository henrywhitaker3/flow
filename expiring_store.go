package flow

import (
	"sync"
	"time"
)

type ExpiringStore[T any] struct {
	store  *Store[T]
	closed chan struct{}

	cancelMutex *sync.Mutex
	cancel      map[string]chan struct{}

	waits  *sync.WaitGroup
	closer *sync.Once
}

func NewExpiringStore[T any]() *ExpiringStore[T] {
	return &ExpiringStore[T]{
		store:       NewStore[T](),
		closed:      make(chan struct{}, 1),
		cancelMutex: &sync.Mutex{},
		cancel:      map[string]chan struct{}{},
		waits:       &sync.WaitGroup{},
		closer:      &sync.Once{},
	}
}

// Close the store, this will block until all items in the
// store have been expired before it returns, so be careful
// over the max expiry time you use when adding items to the
// store
func (e *ExpiringStore[T]) Close() {
	e.closer.Do(func() {
		e.waits.Wait()
		e.closed <- struct{}{}
	})
}

func (e *ExpiringStore[T]) Closed() <-chan struct{} {
	return e.closed
}

type ExpiryCallback[T any] func(key string, val T)

func (e *ExpiringStore[T]) Put(id string, val T, exp time.Duration, callbacks ...ExpiryCallback[T]) {
	e.store.Put(id, val)
	e.expireAfter(id, val, exp, callbacks...)
}

func (e *ExpiringStore[T]) expireAfter(id string, val T, exp time.Duration, callbacks ...ExpiryCallback[T]) {
	cancel := make(chan struct{}, 1)
	e.cancelMutex.Lock()
	e.cancel[id] = cancel
	e.cancelMutex.Unlock()

	timer := time.After(exp)
	e.waits.Add(1)
	go func() {
		defer e.waits.Done()
		select {
		case <-e.Closed():
			return
		case <-cancel:
			return
		case <-timer:
			for _, cb := range callbacks {
				cb(id, val)
			}
			e.store.Delete(id)
		}
	}()
}

func (e *ExpiringStore[T]) Get(id string) (T, bool) {
	return e.store.Get(id)
}

func (e *ExpiringStore[T]) Delete(id string) {
	e.store.Delete(id)
	cancel, ok := e.cancel[id]
	if ok {
		cancel <- struct{}{}
		e.cancelMutex.Lock()
		delete(e.cancel, id)
		e.cancelMutex.Unlock()
	}
}
