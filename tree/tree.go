package tree

import (
	"sync"

	"github.com/e-zhydzetski/orderbook/memtable"
)

func New[K any, V any](compareFunc func(a K, b K) int) *Tree[K, V] {
	return &Tree[K, V]{
		compareFunc: compareFunc,
		nodePool: &sync.Pool{
			New: func() any {
				return new(Node[K, V])
			},
		},
	}
}

type Tree[K any, V any] struct {
	compareFunc func(a K, b K) int
	root        *Node[K, V]
	nodePool    *sync.Pool
}

type Node[K any, V any] struct {
	Key   K
	Value V
	Left  *Node[K, V]
	Right *Node[K, V]
}

func (t *Tree[K, V]) Set(key K, value V) {
	newNode := func() *Node[K, V] {
		nn := t.nodePool.Get().(*Node[K, V])
		nn.Key = key
		nn.Value = value
		nn.Right = nil
		nn.Left = nil
		return nn
	}

	if t.root == nil {
		t.root = newNode()
		return
	}
	cur := t.root
	for {
		switch t.compareFunc(key, cur.Key) {
		case -1:
			if cur.Left == nil {
				cur.Left = newNode()
				return
			}
			cur = cur.Left
		case 1:
			if cur.Right == nil {
				cur.Right = newNode()
				return
			}
			cur = cur.Right
		default:
			return // duplicate key, skip
		}
	}
}

// return remove and continue flags
func (t *Tree[K, V]) iter(node *Node[K, V], f func(key K, val *V) memtable.IteratorAction) (bool, bool) {
	if node.Left != nil {
		rem, cont := t.iter(node.Left, f)
		if rem {
			t.nodePool.Put(node.Left)
			node.Left = nil
		}
		if !cont {
			return false, false
		}
	}

	action := f(node.Key, &node.Value)
	if action == memtable.IAStop {
		return false, false
	}

	// cur node should be removed, but its right child maybe not

	if node.Right == nil { // no right child - just remove cur note and continue
		return true, true
	}

	rem, cont := t.iter(node.Right, f)
	if rem {
		// remove right child and cur node
		t.nodePool.Put(node.Right)
		node.Right = nil
		return true, cont
	}

	// not remove right child, so replace cur node with it
	rc := node.Right
	node.Key = rc.Key
	node.Value = rc.Value
	node.Left = rc.Left
	node.Right = rc.Right
	t.nodePool.Put(rc)

	return false, cont
}

// Iterate tree elements from min to max key, next element may be accessed only after current remove
// element value is mutable
func (t *Tree[K, V]) Iterate(f func(key K, val *V) memtable.IteratorAction) {
	if t.root == nil {
		return
	}

	rem, _ := t.iter(t.root, f)
	if rem {
		t.nodePool.Put(t.root)
		t.root = nil
	}
}
