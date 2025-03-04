package collection

import (
	"encoding/json"
	"sync"
)

// OrderedSet 是一个协程安全的有序集合，按插入顺序维护元素
type OrderedSet[T comparable] struct {
	mu       sync.RWMutex
	elements []T
	indexes  map[T]int
}

// NewOrderedSet 创建一个新的有序集合
func NewSet[T comparable]() Set[T] {
	return &OrderedSet[T]{
		indexes: make(map[T]int),
	}
}

// Add 添加元素到集合
func (os *OrderedSet[T]) Add(element T) {
	os.mu.Lock()
	defer os.mu.Unlock()

	if _, exists := os.indexes[element]; exists {
		return
	}

	os.elements = append(os.elements, element)
	os.indexes[element] = len(os.elements) - 1
}

// Del 从集合中移除元素
func (os *OrderedSet[T]) Del(element T) {
	os.mu.Lock()
	defer os.mu.Unlock()

	if pos, exists := os.indexes[element]; exists {
		delete(os.indexes, element)
		os.elements = append(os.elements[:pos], os.elements[pos+1:]...)
		// 更新索引
		for i := pos; i < len(os.elements); i++ {
			os.indexes[os.elements[i]] = i
		}
	}
}

// Has 检查元素是否在集合中
func (os *OrderedSet[T]) Has(element T) bool {
	os.mu.RLock()
	defer os.mu.RUnlock()

	_, exists := os.indexes[element]
	return exists
}

// Size 返回集合大小
func (os *OrderedSet[T]) Size() int {
	os.mu.RLock()
	defer os.mu.RUnlock()
	return len(os.elements)
}

// Clear 清空集合
func (os *OrderedSet[T]) Clear() {
	os.mu.Lock()
	defer os.mu.Unlock()

	os.elements = nil
	os.indexes = make(map[T]int)
}

// Vals 返回所有元素的切片
func (os *OrderedSet[T]) Vals() []T {
	os.mu.RLock()
	defer os.mu.RUnlock()

	result := make([]T, len(os.elements))
	copy(result, os.elements)
	return result
}

// Each 遍历集合中的所有元素
func (os *OrderedSet[T]) Each(fn func(element T) bool) {
	os.mu.RLock()
	defer os.mu.RUnlock()

	for _, element := range os.elements {
		if !fn(element) {
			break
		}
	}
}

// MarshalJSON 实现json.Marshaler接口
func (os *OrderedSet[T]) MarshalJSON() ([]byte, error) {
	os.mu.RLock()
	defer os.mu.RUnlock()

	return json.Marshal(os.elements)
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (os *OrderedSet[T]) UnmarshalJSON(data []byte) error {
	os.mu.Lock()
	defer os.mu.Unlock()

	os.elements = nil
	os.indexes = make(map[T]int)

	var elements []T
	if err := json.Unmarshal(data, &elements); err != nil {
		return err
	}

	for _, element := range elements {
		if _, exists := os.indexes[element]; !exists {
			os.elements = append(os.elements, element)
			os.indexes[element] = len(os.elements) - 1
		}
	}

	return nil
}
