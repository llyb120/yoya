package stlx

type ArrayList[T comparable] struct {
	list []T
	mu   lock
}

func NewList[T comparable](args ...any) *ArrayList[T] {
	list := &ArrayList[T]{
		list: make([]T, 0, 16),
	}
	for _, arg := range args {
		switch v := arg.(type) {
		case []T:
			list.AddAll(v)
		case T:
			list.Add(v)
		case Collection[T]:
			v.For(func(item T) bool {
				list.add(item)
				return true
			})
		}
	}
	return list
}

func NewSyncList[T comparable](args ...any) *ArrayList[T] {
	list := NewList[T](args...)
	list.mu.sync = true
	return list
}

func (list *ArrayList[T]) Add(item T) {
	list.mu.Lock()
	defer list.mu.Unlock()

	list.add(item)
}

func (list *ArrayList[T]) AddAll(items []T) {
	list.mu.Lock()
	defer list.mu.Unlock()

	list.list = append(list.list, items...)
}

func (list *ArrayList[T]) Get(index int) (T, bool) {
	list.mu.RLock()
	defer list.mu.RUnlock()

	if index < 0 || index >= len(list.list) {
		var zero T
		return zero, false
	}
	return list.list[index], true
}

func (list *ArrayList[T]) Set(index int, item T) {
	list.mu.Lock()
	defer list.mu.Unlock()

	if index < 0 {
		return
	}
	if index >= len(list.list) {
		list.list = append(list.list, make([]T, index-len(list.list)+1)...)
	}
	list.list[index] = item
}

func (list *ArrayList[T]) Len() int {
	list.mu.RLock()
	defer list.mu.RUnlock()

	return len(list.list)
}

func (list *ArrayList[T]) Clear() {
	list.mu.Lock()
	defer list.mu.Unlock()

	list.clear()
}

func (list *ArrayList[T]) For(fn func(item T) bool) {
	list.mu.RLock()
	defer list.mu.RUnlock()

	list.foreach(fn)
}

func (list *ArrayList[T]) Vals() []T {
	list.mu.RLock()
	defer list.mu.RUnlock()

	return list.vals()
}

func (list *ArrayList[T]) Del(item T) {
	list.mu.Lock()
	defer list.mu.Unlock()

	for i, v := range list.list {
		if v == item {
			list.list = append(list.list[:i], list.list[i+1:]...)
			return
		}
	}
}

func (list *ArrayList[T]) Has(item T) bool {
	list.mu.RLock()
	defer list.mu.RUnlock()

	for _, v := range list.list {
		if v == item {
			return true
		}
	}
	return false
}
