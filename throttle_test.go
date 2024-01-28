package flow_test

import (
	"context"
	"testing"
	"time"

	"github.com/henrywhitaker3/flow"
	"github.com/stretchr/testify/assert"
)

func TestItReturnsTheFirstTimeItIsCalled(t *testing.T) {
	do := func(ctx context.Context) (int, error) {
		return 1, nil
	}

	runDo := flow.Throttle[int](do, time.Second)
	out, err := runDo(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 1, out)
}

func TestItThrottlesWhenCalledBeforeDurationHasPassed(t *testing.T) {
	do := func(ctx context.Context) (int, error) {
		return 1, nil
	}

	ctx := context.Background()

	runDo := flow.Throttle[int](do, time.Millisecond)
	out, err := runDo(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, out)

	out, err = runDo(ctx)
	assert.ErrorIs(t, err, flow.ErrThrottled)
	assert.Empty(t, out)
}

func TestItLetsYouCallItAgainAfterDurationHasPassed(t *testing.T) {
	do := func(ctx context.Context) (int, error) {
		return 1, nil
	}

	ctx := context.Background()

	runDo := flow.Throttle[int](do, time.Millisecond)
	out, err := runDo(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, out)

	out, err = runDo(ctx)
	assert.ErrorIs(t, err, flow.ErrThrottled)
	assert.Empty(t, out)

	time.Sleep(time.Millisecond * 2)

	out, err = runDo(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, out)
}
