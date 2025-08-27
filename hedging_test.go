package flow

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type dummyHedged struct {
	do func(context.Context)
}

func (d dummyHedged) Do(ctx context.Context) (struct{}, error) {
	if d.do != nil {
		d.do(ctx)
	}
	return struct{}{}, nil
}

func TestItFiresOffMultipleRequests(t *testing.T) {
	count := &atomic.Int32{}

	_, err := Hedge(context.Background(), dummyHedged{
		do: func(ctx context.Context) {
			count.Add(1)
		},
	}.Do, 3)
	require.Nil(t, err)
	// Cancelling doesn't affect th do func, so just wait for them all to run
	time.Sleep(time.Millisecond)
	require.Equal(t, int32(3), count.Load())
}

func TestItCancelsOtherRequestsAfterOneFinshes(t *testing.T) {
	iters := &atomic.Int32{}
	hits := &atomic.Int32{}

	_, err := Hedge(context.Background(), dummyHedged{
		do: func(ctx context.Context) {
			waits := []time.Duration{0, 0, time.Second}
			iter := iters.Add(1)
			time.Sleep(waits[int(iter)-1])
			if ctx.Err() == nil {
				hits.Add(1)
			}
		},
	}.Do, 3)
	require.Nil(t, err)
	require.Equal(t, int32(1), hits.Load())
}
