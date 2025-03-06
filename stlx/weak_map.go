package stlx

import (
	"runtime"

	"github.com/llyb120/gotool/syncx"
)

// WeakMap 实现了一个键为弱引用的映射
type WeakMap[K comparable, V any] struct {
	mu   syncx.Lock
	data map[K]V
}

// NewWeakMap 创建一个新的 WeakMap
func NewWeakMap[K comparable, V any]() *WeakMap[K, V] {
	return &WeakMap[K, V]{
		data: make(map[K]V),
	}
}

// Set 存储键值对，并为键设置终结器
func (wm *WeakMap[K, V]) Set(key K, value V) {
	// 确保 key 是指针类型
	wm.mu.Lock()
	defer wm.mu.Unlock()

	// 创建一个包装器来避免直接引用 wm
	wrapper := &keyWrapper[K, V]{key: key, weakMap: wm}
	wm.data[key] = value

	// 为包装器设置终结器
	runtime.SetFinalizer(wrapper, func(w *keyWrapper[K, V]) {
		w.weakMap.mu.Lock()
		delete(w.weakMap.data, w.key)
		w.weakMap.mu.Unlock()
	})

}

// keyWrapper 用于避免在终结器中直接引用 WeakMap
type keyWrapper[K comparable, V any] struct {
	key     K
	weakMap *WeakMap[K, V]
}

// Load 获取键对应的值
func (wm *WeakMap[K, V]) Get(key K) (V, bool) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	value, ok := wm.data[key]
	return value, ok
}

// Del 删除键值对
func (wm *WeakMap[K, V]) Del(key K) V {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	if val, ok := wm.data[key]; ok {
		delete(wm.data, key)
		return val
	}
	var zero V
	return zero
}

// Len 返回当前映射中的元素数量
func (wm *WeakMap[K, V]) Len() int {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	return len(wm.data)
}

func (wm *WeakMap[K, V]) Clear() {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	wm.data = make(map[K]V)
}

func (wm *WeakMap[K, V]) For(fn func(key K, value V) bool) {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	for key, value := range wm.data {
		if !fn(key, value) {
			break
		}
	}
}

func (wm *WeakMap[K, V]) Keys() []K {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	keys := make([]K, 0, len(wm.data))
	for key := range wm.data {
		keys = append(keys, key)
	}
	return keys
}

func (wm *WeakMap[K, V]) Vals() []V {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	values := make([]V, 0, len(wm.data))
	for _, value := range wm.data {
		values = append(values, value)
	}
	return values
}
