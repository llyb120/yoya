package collection

type Set[T comparable] interface {
	Add(item T)
	Del(item T)
	Size() int
	Clear()
	Has(item T) bool
	Vals() []T
	ForEach(fn func(item T) bool)
}
