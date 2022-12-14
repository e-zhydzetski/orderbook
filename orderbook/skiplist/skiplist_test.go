package skiplist

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func TestTree(t *testing.T) {
	list := New[int, int](10, func(a int, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	list.Iterate(func(key int, val *int) IteratorAction {
		require.Fail(t, "list should be empty")
		return IAStop
	})

	list.Set(-1, -1)

	var vals []int
	list.Iterate(func(key int, val *int) IteratorAction {
		*val *= 10
		vals = append(vals, *val)
		return IAStop
	})
	require.Equal(t, []int{-10}, vals)

	list.Set(-2, -2)
	list.Set(-3, -3)
	list.Set(1, 1)
	list.Set(3, 3)
	list.Set(2, 2)

	vals = vals[:0]
	list.Iterate(func(key int, val *int) IteratorAction {
		vals = append(vals, *val)
		return IAStop
	})
	require.Equal(t, []int{-3}, vals)

	vals = vals[:0]
	list.Iterate(func(key int, val *int) IteratorAction {
		vals = append(vals, *val)
		return IARemoveAndContinue
	})
	require.Equal(t, []int{-3, -2, -10, 1, 2, 3}, vals) // value -10 has -1 key, order by key

	list.Iterate(func(key int, val *int) IteratorAction {
		require.Fail(t, "list should be empty")
		return IAStop
	})
}

func BenchmarkSet(b *testing.B) {
	newSkipList := func(maxHeight int) *SkipList[int, int] {
		return New[int, int](maxHeight, func(a int, b int) int {
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

	for _, maxHeight := range []int{1, 2, 4, 8, 16, 32} {
		b.Run(fmt.Sprintf("mh-%d", maxHeight), func(b *testing.B) {
			for _, test := range tests {
				b.Run(test.name, func(b *testing.B) {
					list := newSkipList(maxHeight)
					next := test.generator()
					b.ResetTimer()

					for i := 0; i < b.N; i++ {
						x := next()
						list.Set(x, x)
					}
				})
			}
		})
	}
}
