package skiplist

import (
	"fmt"
	"github.com/e-zhydzetski/orderbook/memtable"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
	"time"
)

func TestSkipList(t *testing.T) {
	list := New[int, int](10, func(a int, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})
	list.Iterate(func(key int, val *int) memtable.IteratorAction {
		require.Fail(t, "list should be empty")
		return memtable.IAStop
	})

	list.Upsert(-1, func() int { return -1 }, nil)

	var vals []int
	list.Iterate(func(key int, val *int) memtable.IteratorAction {
		*val *= 10
		vals = append(vals, *val)
		return memtable.IAStop
	})
	require.Equal(t, []int{-10}, vals)

	list.Upsert(-1, nil, func(val *int) {
		*val *= 10
	})

	list.Upsert(-2, func() int { return -2 }, nil)
	list.Upsert(-3, func() int { return -3 }, nil)
	list.Upsert(1, func() int { return 1 }, nil)
	list.Upsert(3, func() int { return 3 }, nil)
	list.Upsert(2, func() int { return 2 }, nil)

	vals = vals[:0]
	list.Iterate(func(key int, val *int) memtable.IteratorAction {
		vals = append(vals, *val)
		return memtable.IAStop
	})
	require.Equal(t, []int{-3}, vals)

	vals = vals[:0]
	list.Iterate(func(key int, val *int) memtable.IteratorAction {
		vals = append(vals, *val)
		return memtable.IARemoveAndContinue
	})
	require.Equal(t, []int{-3, -2, -100, 1, 2, 3}, vals) // value -10 has -1 key, order by key

	list.Iterate(func(key int, val *int) memtable.IteratorAction {
		require.Fail(t, "list should be empty")
		return memtable.IAStop
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
						list.Upsert(x, func() int { return x }, nil)
					}
				})
			}
		})
	}
}
