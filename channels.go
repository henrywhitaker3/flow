package flow

import (
	"context"
	"runtime"
)

type Work[T, U any] func(T) (U, error)

type Response[T any] struct {
	output chan T
	err    chan error
}

func (r *Response[T]) Output() T {
	return <-r.output
}

func (r *Response[T]) Err() error {
	return <-r.err
}

type Channel[T, U any] struct {
	ch chan request[T, U]
	cb Work[T, U]
}

func NewChannel[T, U any](ctx context.Context, bufferSize int, cb Work[T, U]) *Channel[T, U] {
	channel := &Channel[T, U]{
		ch: make(chan request[T, U], bufferSize),
		cb: cb,
	}
	for range runtime.NumCPU() {
		go channel.work(ctx)
	}
	return channel
}

type request[T, U any] struct {
	item     T
	response *Response[U]
}

func (c *Channel[T, U]) Push(item T) *Response[U] {
	response := &Response[U]{
		output: make(chan U, 1),
		err:    make(chan error, 1),
	}
	c.ch <- request[T, U]{
		item:     item,
		response: response,
	}
	return response
}

func (c *Channel[T, U]) work(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-c.ch:
			out, err := c.cb(req.item)
			req.response.output <- out
			req.response.err <- err
		}
	}
}
