package flow

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestOncePerSetsUpNewOnceIfDoesntExist(t *testing.T) {
	once := NewOncePer[string]()
	called := 0
	once.Do("bongo", func() {
		called++
	})
	once.Do("bongo", func() {
		called++
	})
	once.Do("bingo", func() {
		called++
	})
	once.Do("bingo", func() {
		called++
	})
	require.Equal(t, 2, called)
	once.Reset("bongo")
	once.Do("bongo", func() {
		called++
	})
	require.Equal(t, 3, called)
}
