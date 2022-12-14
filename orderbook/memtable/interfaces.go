package memtable

type IteratorAction byte

const (
	IAStop IteratorAction = iota
	IARemoveAndContinue
)

type MemTable[K any, V any] interface {
	Set(key K, value V)
	Iterate(f func(key K, val *V) IteratorAction)
}

type Queue[T any] interface {
	Add(val T)
	Iterate(f func(val *T) IteratorAction)
}
