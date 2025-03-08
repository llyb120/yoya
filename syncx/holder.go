package syncx

import "sync"

type Holder[K comparable, V any] struct {
	mu    sync.RWMutex
	mp    map[K]V
	newFn func() V
}

func NewHolder[K comparable, V any](fn func() V) *Holder[K, V] {
	return &Holder[K, V]{newFn: fn, mp: make(map[K]V)}
}

func (h *Holder[K, V]) Get(key K) V {
	if h.newFn == nil {
		var zero V
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

func (h *Holder[K, V]) Set(key K, value V) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.mp[key] = value
}

func (h *Holder[K, V]) Del(key K) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.mp, key)
}
