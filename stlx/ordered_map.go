package stlx

import "sort"

// orderedMap 是一个协程安全的有序映射，按插入顺序维护键值对
type orderedMap[K comparable, V any] struct {
	mu   lock
	keys []K
	mp   map[K]V
}

// NeworderedMap 创建一个新的有序映射
func NewMap[K comparable, V any](args ...any) *orderedMap[K, V] {
	om := &orderedMap[K, V]{
		mp: make(map[K]V),
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

func NewSyncMap[K comparable, V any](args ...any) *orderedMap[K, V] {
	om := NewMap[K, V](args...)
	om.mu.sync = true
	return om
}

// Set 添加或更新键值对
func (om *orderedMap[K, V]) Set(key K, value V) {
	om.mu.Lock()
	defer om.mu.Unlock()

	om.set(key, value)
}

// Get 获取键对应的值
func (om *orderedMap[K, V]) Get(key K) (V, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	return om.get(key)
}

// Del 删除键值对
func (om *orderedMap[K, V]) Del(key K) V {
	om.mu.Lock()
	defer om.mu.Unlock()

	var index = -1
	for i, k := range om.keys {
		if k == key {
			index = i
			break
		}
	}
	if index == -1 {
		var zero V
		return zero
	} else {
		val := om.mp[key]
		delete(om.mp, key)
		om.keys = append(om.keys[:index], om.keys[index+1:]...)
		return val
	}
}

// Size 返回映射大小
func (om *orderedMap[K, V]) Len() int {
	om.mu.RLock()
	defer om.mu.RUnlock()
	return len(om.keys)
}

// Keys 按插入顺序返回所有键
func (om *orderedMap[K, V]) Keys() []K {
	om.mu.RLock()
	defer om.mu.RUnlock()

	keys := make([]K, len(om.keys))
	copy(keys, om.keys)
	return keys
}

// Vals 按插入顺序返回所有值
func (om *orderedMap[K, V]) Vals() []V {
	om.mu.RLock()
	defer om.mu.RUnlock()

	values := make([]V, len(om.keys))
	for i, key := range om.keys {
		values[i] = om.mp[key]
	}
	return values
}

// Clear 清空映射
func (om *orderedMap[K, V]) Clear() {
	om.mu.Lock()
	defer om.mu.Unlock()

	om.clear()
}

// For 按顺序遍历所有键值对
func (om *orderedMap[K, V]) For(fn func(key K, value V) bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	for _, key := range om.keys {
		if !fn(key, om.mp[key]) {
			break
		}
	}
}

func (om *orderedMap[K, V]) SortByKey(fn func(a, b K) bool) {
	om.mu.Lock()
	defer om.mu.Unlock()

	sort.Slice(om.keys, func(i, j int) bool {
		return fn(om.keys[i], om.keys[j])
	})
}

func (om *orderedMap[K, V]) SortByValue(fn func(a, b V) bool) {
	om.mu.Lock()
	defer om.mu.Unlock()

	type pair struct {
		Index int
		Value V
	}
	values := make([]pair, len(om.keys))
	for i, key := range om.keys {
		values[i] = pair{Index: i, Value: om.mp[key]}
	}
	sort.Slice(values, func(i, j int) bool {
		return fn(values[i].Value, values[j].Value)
	})
	keys := make([]K, len(values))
	for i, v := range values {
		keys[i] = om.keys[v.Index]
	}
	om.keys = keys
}

func (om *orderedMap[K, V]) Fork() *orderedMap[K, V] {
	om.mu.RLock()
	defer om.mu.RUnlock()

	forkMap := NewMap[K, V](om)
	return forkMap
}

func (om *orderedMap[K, V]) Index(key K) int {
	om.mu.RLock()
	defer om.mu.RUnlock()

	for i, k := range om.keys {
		if k == key {
			return i
		}
	}
	return -1
}
