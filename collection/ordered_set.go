package collection

import (
	"encoding/json"
)

type void struct{}

// OrderedSet 是一个协程安全的有序集合，按插入顺序维护元素
type OrderedSet[T comparable] struct {
	mp *OrderedMap[T, void]
}

// NewOrderedSet 创建一个新的有序集合
func NewSet[T comparable]() *OrderedSet[T] {
	return &OrderedSet[T]{
		mp: NewMap[T, void](),
	}
}

// Add 添加元素到集合
func (os *OrderedSet[T]) Add(element T) {
	os.mp.Set(element, struct{}{})
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
	return os.mp.Keys()
}

// Each 遍历集合中的所有元素
func (os *OrderedSet[T]) Each(fn func(element T) bool) {
	os.mp.Each(func(key T, value void) bool {
		return fn(key)
	})
}

// MarshalJSON 实现json.Marshaler接口
func (os *OrderedSet[T]) MarshalJSON() ([]byte, error) {
	os.mp.mu.Lock()
	defer os.mp.mu.Unlock()
	return json.Marshal(os.mp.keys)
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (os *OrderedSet[T]) UnmarshalJSON(data []byte) error {
	os.mp.mu.Lock()
	defer os.mp.mu.Unlock()

	os.Clear()

	var elements []T
	if err := json.Unmarshal(data, &elements); err != nil {
		return err
	}
	for _, element := range elements {
		os.mp.Set(element, struct{}{})
	}

	return nil
}
