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

	<-res.Done()

	assert.Equal(t, 5, res.Out())
	assert.Nil(t, res.Err())
}

func TestItSendsErrorsDownTheChannel(t *testing.T) {
	res := flow.Eventually(context.Background(), func(ctx context.Context) (int, error) {
		return 0, errors.New("bongo")
	})

	<-res.Done()

	assert.Equal(t, 0, res.Out())
	assert.NotNil(t, res.Err())
}

func TestItWorksInTheBackground(t *testing.T) {
	start := time.Now()
	res := flow.Eventually(context.Background(), func(ctx context.Context) (int, error) {
		time.Sleep(time.Millisecond)
		return 10, nil
	})

	afterCall := time.Since(start)
	assert.Less(t, afterCall, time.Millisecond)

	<-res.Done()

	assert.Greater(t, time.Since(start), time.Millisecond)
	assert.Equal(t, 10, res.Out())
	assert.Nil(t, res.Err())
}

func TestItCanWaitForDoneMultipleTimes(t *testing.T) {
	res := flow.Eventually(context.Background(), func(ctx context.Context) (int, error) {
		return 10, nil
	})

	<-res.Done()
	<-res.Done()
}

func TestItBlocksWhenCallingOutDirectly(t *testing.T) {
	start := time.Now()
	res := flow.Eventually(context.Background(), func(ctx context.Context) (int, error) {
		time.Sleep(time.Millisecond)
		return 10, nil
	})

	assert.Equal(t, 10, res.Out())

	afterCall := time.Since(start)
	assert.Greater(t, afterCall, time.Millisecond)
}

func TestItBlocksWhenCallingErrDirectly(t *testing.T) {
	start := time.Now()
	res := flow.Eventually(context.Background(), func(ctx context.Context) (int, error) {
		time.Sleep(time.Millisecond)
		return 10, nil
	})

	assert.Nil(t, res.Err())

	afterCall := time.Since(start)
	assert.Greater(t, afterCall, time.Millisecond)
}

func TestAResultGroupWaitsForAllResultsToResolve(t *testing.T) {
	group := flow.ResultGroup{}

	start := time.Now()
	instant := flow.Eventually(context.Background(), func(ctx context.Context) (int, error) {
		return 1, nil
	})
	slow := flow.Eventually(context.Background(), func(ctx context.Context) (int, error) {
		time.Sleep(time.Millisecond)
		return 2, nil
	})

	group.Add(instant)
	group.Add(slow)

	group.Wait()

	assert.Greater(t, time.Since(start), time.Millisecond)
}

func TestItErrorsWhenAddingAResultWhileWeAreWaiting(t *testing.T) {
	group := flow.ResultGroup{}

	instant := flow.Eventually(context.Background(), func(ctx context.Context) (int, error) {
		return 1, nil
	})
	slow := flow.Eventually(context.Background(), func(ctx context.Context) (int, error) {
		time.Sleep(time.Millisecond * 5)
		return 5, nil
	})

	group.Add(slow)
	go group.Wait()
	// Wait for the goroutine to start up
	time.Sleep(time.Millisecond)
	assert.ErrorIs(t, group.Add(instant), flow.ErrGroupAlreadyWaiting)
}
