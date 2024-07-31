package flow_test

import (
	"testing"
	"time"

	"github.com/henrywhitaker3/flow"
	"github.com/stretchr/testify/require"
)

func TestItExpiresAnItemFromTheStore(t *testing.T) {
	store := flow.NewExpiringStore[string]()
	defer store.Close()
	store.Put("bongo", "bingo", time.Millisecond)
	_, ok := store.Get("bongo")
	require.True(t, ok)
	time.Sleep(time.Millisecond * 2)
	_, ok = store.Get("bongo")
	require.False(t, ok)
}

func TestItWaitsUntilEverythingHasExpiredBeforeClosing(t *testing.T) {
	store := flow.NewExpiringStore[string]()

	start := time.Now()
	store.Put("bongo", "bingo", time.Millisecond*5)
	// This call should block until after 5ms
	store.Close()
	end := time.Since(start)
	require.Greater(t, end, time.Millisecond*5)
}

func TestItCallsCallbacksWhenExpiringAnItem(t *testing.T) {
	store := flow.NewExpiringStore[string]()
	word := "ahkhge"
	called := false
	store.Put("bongo", word, time.Millisecond, func(key string, val string) {
		require.Equal(t, word, val)
		called = true
	})
	store.Close()
	require.True(t, called)
}

func TestItDoesntBlockWhenDeleting(t *testing.T) {
	store := flow.NewExpiringStore[string]()
	word := "ahkhge"
	called := false
	store.Put("bongo", word, time.Millisecond, func(key string, val string) {
		require.Equal(t, word, val)
		called = true
	})
	store.Delete("bongo")
	store.Close()
	require.False(t, called)
}
