package skiplist

import (
	"math/rand"
	"testing"
	"time"
)

func TestTree(t *testing.T) {
	list := New[int, int](func(a int, b int) int {
		if a < b {
			return -1
		}
		if a > b {
			return 1
		}
		return 0
	})

	list.Set(-1, -1)
	list.Set(-2, -2)
	list.Set(-3, -3)
	list.Set(1, 1)
	list.Set(3, 3)
	list.Set(2, 2)
}

func BenchmarkSet(b *testing.B) {
	newSkipList := func() *SkipList[int, int] {
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
			list := newSkipList()
			next := test.generator()
			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				x := next()
				list.Set(x, x)
			}
		})
	}
}
