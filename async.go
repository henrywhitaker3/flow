// Package flow
package flow

import (
	"context"
	"errors"
	"sync"
)

var (
	ErrGroupAlreadyWaiting = errors.New("resultgroup is already waiting")
)

type Result[T any] struct {
	// Channel to eventually receive the output on
	outCh <-chan T
	// Store the result from the output channel
	out T

	// Channel to eventually receive the error on
	errCh <-chan error
	// Store the error from the error channel
	err error

	// Whether the function has returned yet
	done bool
}

func (r *Result[T]) Err() error {
	if !r.done {
		<-r.Done()
	}
	return r.err
}

func (r *Result[T]) Out() T {
	if !r.done {
		<-r.Done()
	}
	return r.out
}

func (r *Result[T]) Done() <-chan struct{} {
	fin := make(chan struct{}, 1)

	if r.done {
		// If its done, send an empty struct so the call doesn't
		// block
		fin <- struct{}{}
	} else {
		// Wait for the results then send the finish signal down
		// the channel
		go func() {
			r.out = <-r.outCh
			r.err = <-r.errCh
			r.done = true
			fin <- struct{}{}
		}()
	}

	return fin
}

func Eventually[T any](ctx context.Context, f Effector[T]) *Result[T] {
	resCh := make(chan T, 1)
	errCh := make(chan error, 1)

	go func() {
		res, err := f(ctx)
		resCh <- res
		errCh <- err
	}()

	return &Result[T]{
		outCh: resCh,
		errCh: errCh,
		done:  false,
	}
}

type ResultGroup struct {
	// Whether everything has resolved
	done bool
	// Whether the group has started work yet
	working bool

	results []resultItem

	wg *sync.WaitGroup
}

type resultItem interface {
	Done() <-chan struct{}
}

// Add an item to the result group
// This will return an error if Wait has already
// been called. When Wait has finished, you can add
// more results to resolve.
func (r *ResultGroup) Add(res resultItem) error {
	if r.working {
		return ErrGroupAlreadyWaiting
	}
	r.results = append(r.results, res)
	return nil
}

// Wait blocks until every result has resolved
func (r *ResultGroup) Wait() {
	r.working = true
	defer func() { r.working = false }()
	if r.done {
		return
	}

	if r.wg == nil {
		r.wg = &sync.WaitGroup{}
	}

	for _, res := range r.results {
		r.wg.Add(1)
		go func(res resultItem) {
			defer r.wg.Done()
			<-res.Done()
		}(res)
	}

	r.wg.Wait()
	r.done = true
}
