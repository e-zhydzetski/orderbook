package tree

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTree(t *testing.T) {
	tree := New[int, int](func(a int, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	tree.Iterate(func(key int, val *int) IteratorAction {
		require.Fail(t, "tree should be empty")
		return IAStop
	})

	tree.Set(1, 1)

	v := 0
	tree.Iterate(func(key int, val *int) IteratorAction {
		*val *= 10
		v += *val
		return IAStop
	})
	require.Equal(t, 10, v)

	tree.Set(2, 2)
	tree.Set(3, 3)

	v = 0
	tree.Iterate(func(key int, val *int) IteratorAction {
		v += *val
		return IAStop
	})
	require.Equal(t, 10, v)

	v = 0
	tree.Iterate(func(key int, val *int) IteratorAction {
		v += *val
		return IARemoveAndContinue
	})
	require.Equal(t, 15, v)

	tree.Iterate(func(key int, val *int) IteratorAction {
		require.Fail(t, "tree should be empty")
		return IAStop
	})
}
