package stlx

import "sync"

type BiMap[K comparable, V comparable] struct {
	mu sync.RWMutex
	*OrderedMap[K, V]
	fMap *OrderedMap[V, K]
}

func NewBiMap[K comparable, V comparable]() *BiMap[K, V] {
	return &BiMap[K, V]{
		OrderedMap: NewMap[K, V](),
		fMap:       NewMap[V, K](),
	}
}

func (m *BiMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.OrderedMap.Set(key, value)
	m.fMap.Set(value, key)
}

func (m *BiMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.OrderedMap.Get(key)
}

func (m *BiMap[K, V]) GetByValue(value V) (K, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.fMap.Get(value)
}

func (m *BiMap[K, V]) Del(key K) V {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok := m.OrderedMap.Get(key)
	if !ok {
		var zero V
		return zero
	}
	m.OrderedMap.Del(key)
	m.fMap.Del(value)
	return value
}

func (m *BiMap[K, V]) DelByValue(value V) K {
	m.mu.Lock()
	defer m.mu.Unlock()
	key, ok := m.fMap.Get(value)
	if !ok {
		var zero K
		return zero
	}
	m.fMap.Del(value)
	m.OrderedMap.Del(key)
	return key
}
