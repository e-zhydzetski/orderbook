package memtable

type IteratorAction byte

const (
	IAStop IteratorAction = iota
	IARemoveAndContinue
)

type MemTable[K any, V any] interface {
	Upsert(key K, onInsert func() V, onUpdate func(val *V))
	// Iterate returns true if last iterator action was continue
	Iterate(f func(key K, val *V) IteratorAction) bool
}

type Queue[T any] interface {
	PushHead(val T)
	// Iterate returns true if last iterator action was continue
	Iterate(f func(val *T) IteratorAction) bool
}
