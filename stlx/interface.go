package stlx

type Map[K any, V any] interface {
	Set(key K, value V)
	Get(key K) (V, bool)
	Del(key K) V
	Len() int
	Keys() []K
	Vals() []V
	Clear()
	For(fn func(key K, value V) bool)
}

type Collection[T any] interface {
	Add(item T)
	Len() int
	Clear()
	Has(item T) bool
	Vals() []T
	For(fn func(item T) bool)
}

type Set[T comparable] interface {
	Collection[T]
	Del(item T)
}

type List[T comparable] interface {
	Collection[T]
	Set(pos int, value T)
}

type innerLock interface {
	lock()
	unlock()
	rlock()
	runlock()
}

type jsonMap[K any, V any] interface {
	Map[K, V]
	innerLock
	set(key K, value V)
	clear()
	foreach(fn func(key K, value V) bool)
}

type jsonCollection[T any] interface {
	Collection[T]
	innerLock
	clear()
	add(item T)
	vals() []T
	foreach(fn func(item T) bool)
}
