package syncx

import (
	"sync"

	"github.com/petermattis/goid"
)

type Holder[V any] struct {
	mu    sync.RWMutex
	mp    map[int64]V
	newFn func() V
}

func NewHolder[V any](fn func() V) *Holder[V] {
	return &Holder[V]{newFn: fn, mp: make(map[int64]V)}
}

func (h *Holder[V]) Get() V {
	h.mu.RLock()
	defer h.mu.RUnlock()
	goid := goid.Get()
	if item, ok := h.mp[goid]; ok {
		return item
	} else {
		// 如果本协程没有，尝试去父协程找
		targetGoid := goid
		for {
			parentGoid, ok := globalGroupHolder.Load(targetGoid)
			if !ok {
				break
			}
			targetGoid = parentGoid.(int64)
			// 如果可以在父协程找到
			if item, ok := h.mp[targetGoid]; ok {
				return item
			}
		}
	}
	if h.newFn == nil {
		var zero V
		return zero
	}
	item := h.newFn()
	h.mp[goid] = item
	return item
}

func (h *Holder[V]) Set(value V) {
	h.mu.Lock()
	defer h.mu.Unlock()
	goid := goid.Get()
	h.mp[goid] = value
}

func (h *Holder[V]) Del() {
	h.mu.Lock()
	defer h.mu.Unlock()
	goid := goid.Get()
	delete(h.mp, goid)
}
