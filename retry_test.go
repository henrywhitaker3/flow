package flow_test

import (
	"context"
	"errors"
	"testing"

	"github.com/henrywhitaker3/flow"
	"github.com/stretchr/testify/assert"
)

func TestItRetriesTheFuncIfItErrors(t *testing.T) {
	calls := 0
	do := func(ctx context.Context) (struct{}, error) {
		defer func() { calls++ }()
		if calls < 2 {
			return struct{}{}, errors.New("bongo")
		}
		return struct{}{}, nil
	}

	retry := flow.Retry(do, 3)
	_, err := retry(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 3, calls)
}

func TestItReturnsTheErrorIfItHitsTheRetyLimit(t *testing.T) {
	calls := 0
	do := func(ctx context.Context) (struct{}, error) {
		defer func() { calls++ }()
		if calls < 2 {
			return struct{}{}, errors.New("bongo")
		}
		return struct{}{}, nil
	}

	retry := flow.Retry(do, 2)
	_, err := retry(context.Background())
	assert.Equal(t, "bongo", err.Error())
	assert.Equal(t, 2, calls)
}

func TestItOnlyCallsOnceIfItPassesFirstTime(t *testing.T) {
	calls := 0
	do := func(ctx context.Context) (struct{}, error) {
		defer func() { calls++ }()
		return struct{}{}, nil
	}

	retry := flow.Retry(do, 3)
	_, err := retry(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 1, calls)
}
