package stlx

type void struct{}

// orderedSet 是一个协程安全的有序集合，按插入顺序维护元素
type orderedSet[T comparable] struct {
	mp *orderedMap[T, void]
}

// NeworderedSet 创建一个新的有序集合
func NewSet[T comparable](args ...any) *orderedSet[T] {
	set := &orderedSet[T]{
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

func NewSyncSet[T comparable](args ...any) *orderedSet[T] {
	set := NewSet[T](args...)
	set.mp.mu.sync = true
	return set
}

// Add 添加元素到集合
func (os *orderedSet[T]) Add(element T) {
	os.add(element)
}

// Del 从集合中移除元素
func (os *orderedSet[T]) Del(element T) {
	os.mp.Del(element)
}

// Has 检查元素是否在集合中
func (os *orderedSet[T]) Has(element T) bool {
	_, ok := os.mp.Get(element)
	return ok
}

// Size 返回集合大小
func (os *orderedSet[T]) Len() int {
	return os.mp.Len()
}

// Clear 清空集合
func (os *orderedSet[T]) Clear() {
	os.mp.Clear()
}

// Vals 返回所有元素的切片
func (os *orderedSet[T]) Vals() []T {
	os.rlock()
	defer os.runlock()
	return os.vals()
}

// Each 遍历集合中的所有元素
func (os *orderedSet[T]) For(fn func(element T) bool) {
	os.foreach(fn)
}
