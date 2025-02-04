package flow

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestItProcessesItemsFromChannel(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ch := NewChannel(ctx, 10, func(item string) (string, error) {
		return strings.ToUpper(item), nil
	})

	first := "bongo"
	second := "bingo"
	firstResp := ch.Push(first)
	secondResp := ch.Push(second)

	require.Nil(t, firstResp.Err())
	require.Nil(t, secondResp.Err())

	require.Equal(t, strings.ToUpper(first), firstResp.Output())
	require.Equal(t, strings.ToUpper(second), secondResp.Output())
}
