package stlx

func (sm *SkipMap[K, V]) clear() {
	sm.header = newSkipListEntryNode[K, V](*new(K), *new(V), maxLevel)
	sm.level = 1
	sm.length = 0
}

// Set 添加或更新键值对
func (sm *SkipMap[K, V]) set(key K, value V) {

	// 用于记录每一层需要更新的节点
	update := make([]*skipListEntryNode[K, V], maxLevel)
	current := sm.header

	// 从最高层开始查找
	for i := sm.level - 1; i >= 0; i-- {
		// 在当前层查找最后一个小于 key 的节点
		for current.forward[i] != nil && sm.less(current.forward[i].key, key) {
			current = current.forward[i]
		}
		update[i] = current
	}

	// 移动到第0层的下一个节点
	current = current.forward[0]

	// 如果找到了键，则更新值
	if current != nil && current.key == key {
		current.value = value
		return
	}

	// 生成一个随机层数
	newLevel := sm.randomLevel()

	// 如果新层数大于当前层数，则更新头节点的 forward 指针
	if newLevel > sm.level {
		for i := sm.level; i < newLevel; i++ {
			update[i] = sm.header
		}
		sm.level = newLevel
	}

	// 创建新节点
	newNode := newSkipListEntryNode(key, value, newLevel)

	// 更新所有受影响的节点的 forward 指针
	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}

	sm.length++
}

func (sm *SkipMap[K, V]) foreach(fn func(key K, value V) bool) {
	current := sm.header.forward[0]

	for current != nil {
		if !fn(current.key, current.value) {
			break
		}
		current = current.forward[0]
	}

}

func (sm *SkipMap[K, V]) lock() {
	sm.mu.Lock()
}

func (sm *SkipMap[K, V]) unlock() {
	sm.mu.Unlock()
}

func (sm *SkipMap[K, V]) rlock() {
	sm.mu.RLock()
}

func (sm *SkipMap[K, V]) runlock() {
	sm.mu.RUnlock()
}

func (sm *SkipMap[K, V]) MarshalJSON() ([]byte, error) {
	return marshalMap[K, V](sm)
}

func (sm *SkipMap[K, V]) UnmarshalJSON(data []byte) error {
	return unmarshalMap[K, V](sm, data)
}
