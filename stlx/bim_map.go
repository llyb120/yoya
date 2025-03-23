package stlx

type BimMap[K comparable, V comparable] struct {
	mu lock
	*orderedMap[K, V]
	fMap *orderedMap[V, Set[K]]
}

func NewBimMap[K comparable, V comparable]() *BimMap[K, V] {
	return &BimMap[K, V]{
		orderedMap: NewMap[K, V](),
		fMap:       NewMap[V, Set[K]](),
	}
}

func NewSyncBimMap[K comparable, V comparable]() *BimMap[K, V] {
	bm := NewBimMap[K, V]()
	bm.mu.sync = true
	return bm
}

func (m *BimMap[K, V]) Set(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.orderedMap.Set(key, value)
	ks, ok := m.fMap.Get(value)
	if !ok {
		ks = NewSet[K]()
		m.fMap.Set(value, ks)
	}
	ks.Add(key)
}

func (m *BimMap[K, V]) Get(key K) (V, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.orderedMap.Get(key)
}

func (m *BimMap[K, V]) GetByValue(value V) ([]K, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	set, ok := m.fMap.Get(value)
	if !ok {
		return nil, false
	}
	return set.Vals(), true
}

func (m *BimMap[K, V]) Del(key K) V {
	m.mu.Lock()
	defer m.mu.Unlock()
	value, ok := m.orderedMap.Get(key)
	if !ok {
		var zero V
		return zero
	}
	m.orderedMap.Del(key)
	ks, ok := m.fMap.Get(value)
	if !ok {
		return value
	}
	ks.Del(key)
	if ks.Len() == 0 {
		m.fMap.Del(value)
	}
	return value
}

func (m *BimMap[K, V]) DelByValue(value V) []K {
	m.mu.Lock()
	defer m.mu.Unlock()
	keys, ok := m.fMap.Get(value)
	if !ok {
		var zero []K
		return zero
	}
	m.fMap.Del(value)
	ks := keys.Vals()
	for _, key := range ks {
		m.orderedMap.Del(key)
	}
	return ks
}
