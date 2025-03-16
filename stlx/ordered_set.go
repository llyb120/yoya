package stlx

type void struct{}

// OrderedSet 是一个协程安全的有序集合，按插入顺序维护元素
type OrderedSet[T comparable] struct {
	mp *OrderedMap[T, void]
}

// NewOrderedSet 创建一个新的有序集合
func NewSet[T comparable](args ...any) *OrderedSet[T] {
	set := &OrderedSet[T]{
		mp: NewMap[T, void](),
	}
	for _, arg := range args {
		switch v := arg.(type) {
		case []T:
			set.addAll(v)
		case T:
			set.Add(v)
		case Collection[T]:
			v.For(func(item T) bool {
				set.add(item)
				return true
			})
		}
	}
	return set
}

// Add 添加元素到集合
func (os *OrderedSet[T]) Add(element T) {
	os.add(element)
}

// Del 从集合中移除元素
func (os *OrderedSet[T]) Del(element T) {
	os.mp.Del(element)
}

// Has 检查元素是否在集合中
func (os *OrderedSet[T]) Has(element T) bool {
	_, ok := os.mp.Get(element)
	return ok
}

// Size 返回集合大小
func (os *OrderedSet[T]) Len() int {
	return os.mp.Len()
}

// Clear 清空集合
func (os *OrderedSet[T]) Clear() {
	os.mp.Clear()
}

// Vals 返回所有元素的切片
func (os *OrderedSet[T]) Vals() []T {
	os.rlock()
	defer os.runlock()
	return os.vals()
}

// Each 遍历集合中的所有元素
func (os *OrderedSet[T]) For(fn func(element T) bool) {
	os.foreach(fn)
}
