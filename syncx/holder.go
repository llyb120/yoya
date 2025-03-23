package syncx

import (
	"sync"

	"github.com/petermattis/goid"
)

type Holder[V any] struct {
	*sync.RWMutex
	mp       map[int64]V
	once     sync.Once
	InitFunc func() V
}

func (h *Holder[V]) init() {
	h.once.Do(func() {
		h.RWMutex = &sync.RWMutex{}
		h.mp = make(map[int64]V)
	})
}

func (h *Holder[V]) Get() V {
	h.init()
	h.RLock()
	goid := goid.Get()
	if item, ok := h.mp[goid]; ok {
		h.RUnlock()
		return item
	} else {
		// 如果本协程没有，尝试去父协程找
		targetGoid := goid
		for {
			parentGoid, ok := globalGroupHolder.Get(targetGoid)
			if !ok {
				break
			}
			targetGoid = parentGoid
			// 如果可以在父协程找到
			if item, ok := h.mp[targetGoid]; ok {
				return item
			}
		}
	}
	h.RUnlock()
	if h.InitFunc == nil {
		var zero V
		return zero
	}
	item := h.InitFunc()
	h.Lock()
	h.mp[goid] = item
	h.Unlock()
	return item
}

func (h *Holder[V]) Set(value V) {
	h.init()
	h.Lock()
	defer h.Unlock()
	goid := goid.Get()
	h.mp[goid] = value
}

func (h *Holder[V]) Del() V {
	h.init()
	h.Lock()
	defer h.Unlock()
	goid := goid.Get()
	value := h.mp[goid]
	delete(h.mp, goid)
	return value
}
