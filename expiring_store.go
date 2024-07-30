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

type ExpiryCallbacks[T any] func(T)

func (e *ExpiringStore[T]) Put(id string, val T, exp time.Duration, callbacks ...ExpiryCallbacks[T]) {
	e.store.Put(id, val)
	e.waits.Add(1)
	cancel := make(chan struct{}, 1)
	e.cancelMutex.Lock()
	e.cancel[id] = cancel
	e.cancelMutex.Unlock()
	go func() {
		defer e.waits.Done()
		timer := time.NewTimer(exp)
		defer timer.Stop()
		select {
		case <-e.Closed():
			return
		case <-cancel:
			return
		case <-timer.C:
			item, ok := e.Get(id)
			if ok {
				e.store.Delete(id)
				for _, cb := range callbacks {
					cb(item)
				}
			}
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
