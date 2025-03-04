package collection

type Map[K comparable, V any] interface {
	Set(key K, value V)
	Get(key K) (V, bool)
	Del(key K) V
	Size() int
	Keys() []K
	Vals() []V
	Clear()
	ForEach(fn func(key K, value V) bool)
}
