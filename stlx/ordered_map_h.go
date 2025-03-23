package stlx

func (om *orderedMap[K, V]) clear() {
	om.keys = nil
	om.mp = make(map[K]V)
}

func (om *orderedMap[K, V]) set(key K, value V) {
	if _, exists := om.mp[key]; exists {
		// 如果键已存在，只更新值
		om.mp[key] = value
		return
	}

	// 添加到映射
	om.keys = append(om.keys, key)
	om.mp[key] = value
}

func (om *orderedMap[K, V]) get(key K) (V, bool) {
	v, ok := om.mp[key]
	return v, ok
}

func (om *orderedMap[K, V]) foreach(fn func(key K, value V) bool) {
	for _, key := range om.keys {
		if !fn(key, om.mp[key]) {
			break
		}
	}
}

func (om *orderedMap[K, V]) lock() {
	om.mu.Lock()
}

func (om *orderedMap[K, V]) unlock() {
	om.mu.Unlock()
}

func (om *orderedMap[K, V]) rlock() {
	om.mu.RLock()
}

func (om *orderedMap[K, V]) runlock() {
	om.mu.RUnlock()
}

func (om *orderedMap[K, V]) MarshalJSON() ([]byte, error) {
	return marshalMap[K, V](om)
}

func (om *orderedMap[K, V]) UnmarshalJSON(data []byte) error {
	return unmarshalMap[K, V](om, data)
}
