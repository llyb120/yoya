package collection

type Set[T comparable] interface {
	Add(item T)
	Del(item T)
	Len() int
	Clear()
	Has(item T) bool
	Vals() []T
	Each(fn func(item T) bool)
}
