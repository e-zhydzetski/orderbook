package queue

import (
	"github.com/e-zhydzetski/strips-tt/orderbook/memtable"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	q := New[int]()
	q.Iterate(func(val *int) memtable.IteratorAction {
		require.Fail(t, "queue should be empty")
		return memtable.IAStop
	})

	q.Add(1)

	k := 0
	q.Iterate(func(val *int) memtable.IteratorAction {
		*val *= 10
		k += *val
		return memtable.IAStop
	})
	require.Equal(t, 10, k)

	q.Add(2)
	q.Add(3)

	k = 0
	q.Iterate(func(val *int) memtable.IteratorAction {
		k += *val
		return memtable.IAStop
	})
	require.Equal(t, 10, k)

	k = 0
	q.Iterate(func(val *int) memtable.IteratorAction {
		k += *val
		return memtable.IARemoveAndContinue
	})
	require.Equal(t, 15, k)

	q.Iterate(func(val *int) memtable.IteratorAction {
		require.Fail(t, "queue should be empty")
		return memtable.IAStop
	})
}
