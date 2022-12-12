package queue

import "testing"
import "github.com/stretchr/testify/require"

func TestQueue(t *testing.T) {
	q := New[int]()
	q.Iterate(func(val *int, removeAndContinue func()) {
		require.Fail(t, "queue should be empty")
	})

	q.Add(1)

	k := 0
	q.Iterate(func(val *int, removeAndContinue func()) {
		*val *= 10
		k += *val
	})
	require.Equal(t, 10, k)

	q.Add(2)
	q.Add(3)

	k = 0
	q.Iterate(func(val *int, removeAndContinue func()) {
		k += *val
	})
	require.Equal(t, 10, k)

	k = 0
	q.Iterate(func(val *int, removeAndContinue func()) {
		k += *val
		removeAndContinue()
	})
	require.Equal(t, 15, k)

	q.Iterate(func(val *int, removeAndContinue func()) {
		require.Fail(t, "queue should be empty")
	})
}
