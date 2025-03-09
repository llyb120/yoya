package stlx

func (sl *SkipList[T]) clear() {
	sl.header = newSkipListNode(*new(T), maxLevel)
	sl.level = 1
	sl.length = 0
}

func (sl *SkipList[T]) add(value T) {
	update := make([]*skipListNode[T], maxLevel)
	current := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && sl.less(current.forward[i].value, value) {
			current = current.forward[i]
		}
		update[i] = current
	}

	newLevel := sl.randomLevel()
	if newLevel > sl.level {
		for i := sl.level; i < newLevel; i++ {
			update[i] = sl.header
		}
		sl.level = newLevel
	}

	newNode := newSkipListNode(value, newLevel)
	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}

	sl.length++
}

func (sl *SkipList[T]) foreach(fn func(value T) bool) {
	current := sl.header.forward[0]
	for current != nil {
		if !fn(current.value) {
			break
		}
		current = current.forward[0]
	}
}

func (sl *SkipList[T]) vals() []T {
	result := make([]T, 0, sl.length)
	current := sl.header.forward[0]
	for current != nil {
		result = append(result, current.value)
		current = current.forward[0]
	}
	return result
}

// Lock 获取跳表的写锁
func (sl *SkipList[T]) lock() {
	sl.mu.Lock()
}

// Unlock 释放跳表的写锁
func (sl *SkipList[T]) unlock() {
	sl.mu.Unlock()
}

// RLock 获取跳表的读锁
func (sl *SkipList[T]) rlock() {
	sl.mu.RLock()
}

// RUnlock 释放跳表的读锁
func (sl *SkipList[T]) runlock() {
	sl.mu.RUnlock()
}

// MarshalJSON 实现json.Marshaler接口，将跳表序列化为JSON
func (sl *SkipList[T]) MarshalJSON() ([]byte, error) {
	return marshalCollection[T](sl)
}

// UnmarshalJSON 实现json.Unmarshaler接口，从JSON反序列化为跳表
func (sl *SkipList[T]) UnmarshalJSON(data []byte) error {
	return unmarshalCollection[T](sl, data)
}
