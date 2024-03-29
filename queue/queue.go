package queue

import (
	"sync"

	"github.com/e-zhydzetski/orderbook/memtable"
)

type Node[T any] struct {
	Data T
	Next *Node[T]
}

func New[T any]() *Queue[T] {
	return &Queue[T]{
		nodePool: &sync.Pool{
			New: func() any {
				return new(Node[T])
			},
		},
	}
}

type Queue[T any] struct {
	Head     *Node[T]
	Tail     *Node[T]
	nodePool *sync.Pool
}

func (q *Queue[T]) Add(val T) {
	newNode := func() *Node[T] {
		nn := q.nodePool.Get().(*Node[T])
		nn.Data = val
		nn.Next = nil
		return nn
	}

	if q.Head == nil {
		q.Head = newNode()
		q.Tail = q.Head
		return
	}
	q.Head.Next = newNode()
	q.Head = q.Head.Next
}

func (q *Queue[T]) Iterate(f func(val *T) memtable.IteratorAction) {
	cur := q.Tail
	for cur != nil {
		action := f(&cur.Data)
		if action == memtable.IAStop {
			break
		}
		// remove from queue
		q.Tail = cur.Next
		if q.Tail == nil {
			q.Head = nil
		}
		q.nodePool.Put(cur)
		// and get next
		cur = q.Tail
	}
}
