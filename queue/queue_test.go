package queue

import (
	"github.com/e-zhydzetski/orderbook/memtable"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	q := New[int]()
	q.Iterate(func(val *int) memtable.IteratorAction {
		require.Fail(t, "queue should be empty")
		return memtable.IAStop
	})

	q.Add(-1)

	var vals []int
	q.Iterate(func(val *int) memtable.IteratorAction {
		*val *= 10
		vals = append(vals, *val)
		return memtable.IAStop
	})
	require.Equal(t, []int{-10}, vals)

	q.Add(-2)
	q.Add(-3)
	q.Add(1)
	q.Add(3)
	q.Add(2)

	vals = vals[:0]
	q.Iterate(func(val *int) memtable.IteratorAction {
		vals = append(vals, *val)
		return memtable.IAStop
	})
	require.Equal(t, []int{-10}, vals)

	vals = vals[:0]
	q.Iterate(func(val *int) memtable.IteratorAction {
		vals = append(vals, *val)
		return memtable.IARemoveAndContinue
	})
	require.Equal(t, []int{-10, -2, -3, 1, 3, 2}, vals)

	q.Iterate(func(val *int) memtable.IteratorAction {
		require.Fail(t, "queue should be empty")
		return memtable.IAStop
	})
}
