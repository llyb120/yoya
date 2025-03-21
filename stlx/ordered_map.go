package stlx

// OrderedMap 是一个协程安全的有序映射，按插入顺序维护键值对
type OrderedMap[K comparable, V any] struct {
	mu      lock
	keys    []K
	mp      map[K]V
	indexes map[K]int
}

// NewOrderedMap 创建一个新的有序映射
func NewMap[K comparable, V any](args ...any) *OrderedMap[K, V] {
	om := &OrderedMap[K, V]{
		mp:      make(map[K]V),
		indexes: make(map[K]int),
	}
	for _, arg := range args {
		switch v := arg.(type) {
		case map[K]V:
			for k, v := range v {
				om.set(k, v)
			}
		case Map[K, V]:
			v.For(func(key K, value V) bool {
				om.set(key, value)
				return true
			})
		}
	}
	return om
}

func NewSyncMap[K comparable, V any](args ...any) *OrderedMap[K, V] {
	om := NewMap[K, V](args...)
	om.mu.sync = true
	return om
}

// Set 添加或更新键值对
func (om *OrderedMap[K, V]) Set(key K, value V) {
	om.mu.Lock()
	defer om.mu.Unlock()

	om.set(key, value)
}

// Get 获取键对应的值
func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	val, exists := om.mp[key]
	return val, exists
}

// Del 删除键值对
func (om *OrderedMap[K, V]) Del(key K) V {
	om.mu.Lock()
	defer om.mu.Unlock()

	index, exists := om.indexes[key]
	if !exists {
		var zero V
		return zero
	} else {
		val := om.mp[key]
		delete(om.mp, key)
		delete(om.indexes, key)
		om.keys = append(om.keys[:index], om.keys[index+1:]...)
		return val
	}
}

// Size 返回映射大小
func (om *OrderedMap[K, V]) Len() int {
	om.mu.RLock()
	defer om.mu.RUnlock()
	return len(om.keys)
}

// Keys 按插入顺序返回所有键
func (om *OrderedMap[K, V]) Keys() []K {
	om.mu.RLock()
	defer om.mu.RUnlock()

	keys := make([]K, len(om.keys))
	copy(keys, om.keys)
	return keys
}

// Vals 按插入顺序返回所有值
func (om *OrderedMap[K, V]) Vals() []V {
	om.mu.RLock()
	defer om.mu.RUnlock()

	values := make([]V, len(om.keys))
	for i, key := range om.keys {
		values[i] = om.mp[key]
	}
	return values
}

// Clear 清空映射
func (om *OrderedMap[K, V]) Clear() {
	om.mu.Lock()
	defer om.mu.Unlock()

	om.clear()
}

// For 按顺序遍历所有键值对
func (om *OrderedMap[K, V]) For(fn func(key K, value V) bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	for _, key := range om.keys {
		if !fn(key, om.mp[key]) {
			break
		}
	}
}
