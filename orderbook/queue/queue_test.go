package queue

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	q := New[int]()
	q.Iterate(func(val *int) IteratorAction {
		require.Fail(t, "queue should be empty")
		return IAStop
	})

	q.Add(1)

	k := 0
	q.Iterate(func(val *int) IteratorAction {
		*val *= 10
		k += *val
		return IAStop
	})
	require.Equal(t, 10, k)

	q.Add(2)
	q.Add(3)

	k = 0
	q.Iterate(func(val *int) IteratorAction {
		k += *val
		return IAStop
	})
	require.Equal(t, 10, k)

	k = 0
	q.Iterate(func(val *int) IteratorAction {
		k += *val
		return IARemoveAndContinue
	})
	require.Equal(t, 15, k)

	q.Iterate(func(val *int) IteratorAction {
		require.Fail(t, "queue should be empty")
		return IAStop
	})
}
