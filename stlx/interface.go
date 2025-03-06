package stlx

type Map[K comparable, V any] interface {
	Set(key K, value V)
	Get(key K) (V, bool)
	Del(key K) V
	Len() int
	Keys() []K
	Vals() []V
	Clear()
	For(fn func(key K, value V) bool)
}

type Collection[T comparable] interface {
	Add(item T)
	Del(item T)
	Len() int
	Clear()
	Has(item T) bool
	Vals() []T
	For(fn func(item T) bool)
}

type Set[T comparable] Collection[T]

type List[T comparable] interface {
	Collection[T]
	Set(pos int, value T)
}
