package tree

func New[K any, V any](compareFunc func(a K, b K) int) *Tree[K, V] {
	return &Tree[K, V]{
		compareFunc: compareFunc,
	}
}

type Tree[K any, V any] struct {
	compareFunc func(a K, b K) int
	root        *Node[K, V]
}

type Node[K any, V any] struct {
	Key   K
	Value V
	Left  *Node[K, V]
	Right *Node[K, V]
}

func (t *Tree[K, V]) Set(key K, value V) {
	nn := &Node[K, V]{
		Key:   key,
		Value: value,
	}

	if t.root == nil {
		t.root = nn
		return
	}
	cur := t.root
	for {
		switch t.compareFunc(key, cur.Key) {
		case -1:
			if cur.Left == nil {
				cur.Left = nn
				return
			}
			cur = cur.Left
		case 1:
			if cur.Right == nil {
				cur.Right = nn
				return
			}
			cur = cur.Right
		default:
			return // duplicate key, skip
		}
	}
}

// return remove and continue flags
func iter[K any, V any](node *Node[K, V], f func(key K, val *V, removeAndContinue func())) (bool, bool) {
	if node == nil {
		return true, true
	}
	rem, cont := iter(node.Left, f)
	if rem {
		node.Left = nil
	}
	if !cont {
		return false, false
	}

	remAndCont := false
	f(node.Key, &node.Value, func() {
		remAndCont = true
	})
	if !remAndCont {
		return false, false
	}

	// cur node should be removed, but its right child maybe not

	rem, cont = iter(node.Right, f)
	if rem {
		// remove right child and cur node
		node.Right = nil
		return true, cont
	}

	// not remove right child, so replace cur node with it
	rc := node.Right
	node.Key = rc.Key
	node.Value = rc.Value
	node.Left = rc.Left
	node.Right = rc.Right

	return false, cont
}

// Iterate tree elements from min to max key, next element may be accessed only after current remove
// element value is mutable
func (t *Tree[K, V]) Iterate(f func(key K, val *V, removeAndContinue func())) {
	rem, _ := iter(t.root, f)
	if rem {
		t.root = nil
	}
}
