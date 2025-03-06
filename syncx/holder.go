package syncx

import "sync"

type Holder[T any] struct {
	mu    sync.RWMutex
	mp    map[string]T
	newFn func() T
}

func NewHolder[T any](fn func() T) *Holder[T] {
	return &Holder[T]{newFn: fn, mp: make(map[string]T)}
}

func (h *Holder[T]) Get(key string) T {
	if h.newFn == nil {
		var zero T
		return zero
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	if item, ok := h.mp[key]; ok {
		return item
	}
	item := h.newFn()
	h.mp[key] = item
	return item
}

func (h *Holder[T]) Set(key string, value T) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.mp[key] = value
}

func (h *Holder[T]) Del(key string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.mp, key)
}
