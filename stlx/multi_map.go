package stlx

import "sync"

type MultiMap[K comparable, V any] struct {
	mu sync.RWMutex
	*OrderedMap[K, []V]
}

func NewMultiMap[K comparable, V any]() *MultiMap[K, V] {
	return &MultiMap[K, V]{
		OrderedMap: NewMap[K, []V](),
	}
}

func (m *MultiMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	item, ok := m.Get(key)
	if !ok {
		item = []V{}
	}
	item = append(item, value)
	m.OrderedMap.Set(key, item)
}

func (m *MultiMap[K, V]) GetLast(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	item, ok := m.Get(key)
	if !ok {
		var zero V
		return zero, false
	}
	if len(item) == 0 {
		var zero V
		return zero, false
	}
	return item[len(item)-1], true
}
