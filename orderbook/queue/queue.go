package queue

type Node[T any] struct {
	Data T
	Next *Node[T]
}

func New[T any]() *Queue[T] {
	return &Queue[T]{}
}

type Queue[T any] struct {
	Head *Node[T]
	Tail *Node[T]
}

func (q *Queue[T]) Add(val T) {
	nn := &Node[T]{
		Data: val,
	}
	if q.Head == nil {
		q.Head = nn
		q.Tail = q.Head
		return
	}
	q.Head.Next = nn
	q.Head = nn
}

func (q *Queue[T]) Iterate(f func(val *T, removeAndContinue func())) {
	cur := q.Tail
	for cur != nil {
		remAndCont := false
		f(&cur.Data, func() {
			remAndCont = true
		})
		if !remAndCont {
			break
		}
		// remove from queue
		q.Tail = cur.Next
		if q.Tail == nil {
			q.Head = nil
		}
		// and get next
		cur = q.Tail
	}
}