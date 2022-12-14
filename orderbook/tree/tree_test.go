package tree

import (
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
