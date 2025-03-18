package stlx

func (list *ArrayList[T]) add(item T) {
	list.list = append(list.list, item)
}

func (list *ArrayList[T]) clear() {
	list.list = list.list[:0]
}

func (list *ArrayList[T]) vals() []T {
	var cp = make([]T, len(list.list))
	copy(cp, list.list)
	return cp
}

func (list *ArrayList[T]) foreach(fn func(item T) bool) {
	for _, item := range list.list {
		if !fn(item) {
			break
		}
	}
}

func (list *ArrayList[T]) lock() {
	list.mu.Lock()
}

func (list *ArrayList[T]) unlock() {
	list.mu.Unlock()
}

func (list *ArrayList[T]) rlock() {
	list.mu.RLock()
}

func (list *ArrayList[T]) runlock() {
	list.mu.RUnlock()
}

func (list *ArrayList[T]) MarshalJSON() ([]byte, error) {
	return marshalCollection[T](list)
}

func (list *ArrayList[T]) UnmarshalJSON(data []byte) error {
	return unmarshalCollection[T](list, data)
}
