package memtable

type IteratorAction byte

const (
	IAStop IteratorAction = iota
	IARemoveAndContinue
)

type MemTable[K any, V any] interface {
	Set(key K, value V)
	// Iterate returns true if last iterator action was continue
	Iterate(f func(key K, val *V) IteratorAction) bool
}

type Queue[T any] interface {
	Add(val T)
	// Iterate returns true if last iterator action was continue
	Iterate(f func(val *T) IteratorAction) bool
}
