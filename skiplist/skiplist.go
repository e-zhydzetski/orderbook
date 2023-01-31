package skiplist

import (
	"math/rand"
	"sync"

	"github.com/e-zhydzetski/orderbook/memtable"
)

func New[K any, V any](maxHeight int, compareFunc func(a K, b K) int) *SkipList[K, V] {
	nodePool := &sync.Pool{
		New: func() any {
			return &Node[K, V]{
				Next: make([]*Node[K, V], 0, maxHeight+1),
			}
		},
	}
	return &SkipList[K, V]{
		maxRandHeightMask: 2 << (maxHeight - 1),
		compareFunc:       compareFunc,
		head:              nodePool.Get().(*Node[K, V]),
		nodePool:          nodePool,
		tmpPrevNodes:      make([]*Node[K, V], 0, maxHeight+1),
	}
}

type SkipList[K any, V any] struct {
	maxRandHeightMask int
	compareFunc       func(a K, b K) int
	head              *Node[K, V]
	nodePool          *sync.Pool

	tmpPrevNodes []*Node[K, V]
}

type Node[K any, V any] struct {
	Key   K
	Value V
	Next  []*Node[K, V]
}

func (s *SkipList[K, V]) Upsert(key K, onInsert func() V, onUpdate func(val *V)) {
	level := 0
	//nolint:gosec // math random is OK
	for r := rand.Intn(s.maxRandHeightMask); r&1 == 1; r >>= 1 {
		level++
	}
	if level >= len(s.head.Next) {
		level = len(s.head.Next)
		// s.head.Next = append(s.head.Next, nil)
	}

	s.tmpPrevNodes = s.tmpPrevNodes[:level+1]

	cur := s.head
	for i := len(s.head.Next) - 1; i >= 0; i-- {
		for ; cur.Next[i] != nil; cur = cur.Next[i] {
			comp := s.compareFunc(cur.Next[i].Key, key)
			if comp > 0 {
				break
			}
			if comp == 0 {
				// update
				if onUpdate != nil {
					onUpdate(&cur.Next[i].Value)
				}
				// s.tmpPrevNodes = s.tmpPrevNodes[:0]
				return
			}
		}
		if i <= level {
			s.tmpPrevNodes[i] = cur
		}
	}

	// insert
	if onInsert == nil {
		return
	}

	nn := s.nodePool.Get().(*Node[K, V])
	nn.Key = key
	nn.Value = onInsert()
	nn.Next = nn.Next[:level+1]

	if level == len(s.head.Next) {
		// grow head
		s.head.Next = append(s.head.Next, nil)
		s.tmpPrevNodes[level] = s.head
	}

	for i := 0; i <= level; i++ {
		nn.Next[i] = s.tmpPrevNodes[i].Next[i]
		s.tmpPrevNodes[i].Next[i] = nn
	}
	// s.tmpPrevNodes = s.tmpPrevNodes[:0]
}

// Iterate tree elements from min to max key, next element may be accessed only after current remove
// element value is mutable
func (s *SkipList[K, V]) Iterate(f func(key K, val *V) memtable.IteratorAction) bool {
	if len(s.head.Next) == 0 {
		return false
	}

	var action memtable.IteratorAction
	for s.head.Next[0] != nil {
		cur := s.head.Next[0]

		action = f(cur.Key, &cur.Value)
		if action == memtable.IAStop {
			break
		}

		// remove cur and get next
		copy(s.head.Next, cur.Next)

		s.nodePool.Put(cur)
	}

	return action == memtable.IARemoveAndContinue
}
