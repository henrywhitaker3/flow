package flow_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/henrywhitaker3/flow"
	"github.com/stretchr/testify/assert"
)

func TestItSendsAResultDownTheChannel(t *testing.T) {
	res := flow.Eventually(context.Background(), func(ctx context.Context) (int, error) {
		return 5, nil
	})
	out := <-res
	assert.Equal(t, 5, out.Out())
	assert.Nil(t, out.Err())
}

func TestItSendsErrorsDownTheChannel(t *testing.T) {
	res := flow.Eventually(context.Background(), func(ctx context.Context) (int, error) {
		return 0, errors.New("bongo")
	})
	out := <-res
	assert.Equal(t, 0, out.Out())
	assert.NotNil(t, out.Err())
}

func TestItWorksInTheBackground(t *testing.T) {
	start := time.Now()
	res := flow.Eventually(context.Background(), func(ctx context.Context) (int, error) {
		time.Sleep(time.Millisecond)
		return 10, nil
	})

	afterCall := time.Since(start)
	assert.Less(t, afterCall, time.Millisecond)

	out := <-res

	assert.Greater(t, time.Since(start), time.Millisecond)
	assert.Equal(t, 10, out.Out())
	assert.Nil(t, out.Err())
}
