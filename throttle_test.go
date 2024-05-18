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

	runDo := flow.Throttle(do, time.Second)
	out, err := runDo(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 1, out)
}

func TestItThrottlesWhenCalledBeforeDurationHasPassed(t *testing.T) {
	do := func(ctx context.Context) (int, error) {
		return 1, nil
	}

	ctx := context.Background()

	runDo := flow.Throttle(do, time.Millisecond)
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

	runDo := flow.Throttle(do, time.Millisecond)
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

func TestItSilentlyReturnsTheStartValueWhenCalledSilently(t *testing.T) {
	i := 1
	do := func(ctx context.Context) (int, error) {
		defer func() { i++ }()
		return i, nil
	}

	ctx := context.Background()

	run := flow.SilentThrottle(do, time.Millisecond)

	out, err := run(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, out)
	out, err = run(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 1, out)

	time.Sleep(time.Millisecond * 2)

	out, err = run(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 2, out)
}
