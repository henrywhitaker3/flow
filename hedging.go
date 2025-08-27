package flow

import (
	"context"
	"net/http"
)

func Hedge[T any](ctx context.Context, f Effector[T], count int) (T, error) {
	ops := []context.CancelFunc{}

	respCh := make(chan T, 1)
	errCh := make(chan error, 1)

	for range count {
		ctx, cancel := context.WithCancel(ctx)
		ops = append(ops, cancel)
		go func() {
			resp, err := f(ctx)
			respCh <- resp
			errCh <- err
		}()
	}

	resp := <-respCh
	err := <-errCh

	for _, cancel := range ops {
		cancel()
	}

	return resp, err
}

type HedgeClient struct {
	c     *http.Client
	count int
}

func NewHedgeClient(c *http.Client, count int) *HedgeClient {
	return &HedgeClient{
		c:     c,
		count: count,
	}
}

func (h *HedgeClient) Do(r *http.Request) (*http.Response, error) {
	return Hedge(r.Context(), func(ctx context.Context) (*http.Response, error) {
		return h.c.Do(r)
	}, h.count)
}
