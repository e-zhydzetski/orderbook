package tree

import (
	"github.com/e-zhydzetski/strips-tt/orderbook/memtable"
	"math/rand"
	"testing"
	"time"

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
	tree.Iterate(func(key int, val *int) memtable.IteratorAction {
		require.Fail(t, "tree should be empty")
		return memtable.IAStop
	})

	tree.Set(-1, -1)

	var vals []int
	tree.Iterate(func(key int, val *int) memtable.IteratorAction {
		*val *= 10
		vals = append(vals, *val)
		return memtable.IAStop
	})
	require.Equal(t, []int{-10}, vals)

	tree.Set(-2, -2)
	tree.Set(-3, -3)
	tree.Set(1, 1)
	tree.Set(3, 3)
	tree.Set(2, 2)

	vals = vals[:0]
	tree.Iterate(func(key int, val *int) memtable.IteratorAction {
		vals = append(vals, *val)
		return memtable.IAStop
	})
	require.Equal(t, []int{-3}, vals)

	vals = vals[:0]
	tree.Iterate(func(key int, val *int) memtable.IteratorAction {
		vals = append(vals, *val)
		return memtable.IARemoveAndContinue
	})
	require.Equal(t, []int{-3, -2, -10, 1, 2, 3}, vals) // value -10 has -1 key, order by key

	tree.Iterate(func(key int, val *int) memtable.IteratorAction {
		require.Fail(t, "tree should be empty")
		return memtable.IAStop
	})
}

func BenchmarkSet(b *testing.B) {
	newTree := func() *Tree[int, int] {
		return New[int, int](func(a int, b int) int {
			if a < b {
				return -1
			}
			if a > b {
				return 1
			}
			return 0
		})
	}

	tests := []struct {
		name      string
		generator func() func() int
	}{
		{
			"inc",
			func() func() int {
				x := 0
				return func() int {
					x++
					return x
				}
			},
		},
		{
			"dec",
			func() func() int {
				x := 0
				return func() int {
					x--
					return x
				}
			},
		},
		{
			"random",
			func() func() int {
				rand.Seed(time.Now().UnixNano())
				return func() int {
					return rand.Int()
				}
			},
		},
		{
			"const",
			func() func() int {
				return func() int {
					return 777
				}
			},
		},
	}

	for _, test := range tests {
		b.Run(test.name, func(b *testing.B) {
			tree := newTree()
			next := test.generator()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				x := next()
				tree.Set(x, x)
			}
		})
	}
}
