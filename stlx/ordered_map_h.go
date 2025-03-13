package stlx

func (om *OrderedMap[K, V]) clear() {
	om.keys = nil
	om.mp = make(map[K]V)
	om.indexes = make(map[K]int)
}

func (om *OrderedMap[K, V]) set(key K, value V) {
	if _, exists := om.mp[key]; exists {
		// 如果键已存在，只更新值
		om.mp[key] = value
		return
	}

	// 添加到映射
	om.keys = append(om.keys, key)
	om.mp[key] = value
	om.indexes[key] = len(om.keys) - 1
}

func (om *OrderedMap[K, V]) foreach(fn func(key K, value V) bool) {
	for _, key := range om.keys {
		if !fn(key, om.mp[key]) {
			break
		}
	}
}

func (om *OrderedMap[K, V]) lock() {
	om.mu.Lock()
}

func (om *OrderedMap[K, V]) unlock() {
	om.mu.Unlock()
}

func (om *OrderedMap[K, V]) rlock() {
	om.mu.RLock()
}

func (om *OrderedMap[K, V]) runlock() {
	om.mu.RUnlock()
}

func (om *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	return marshalMap[K, V](om)
}

func (om *OrderedMap[K, V]) UnmarshalJSON(data []byte) error {
	return unmarshalMap[K, V](om, data)
}
