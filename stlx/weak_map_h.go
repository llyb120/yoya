package stlx

import "runtime"

func (wm *WeakMap[K, V]) clear() {
	wm.data = make(map[K]V)
}

func (wm *WeakMap[K, V]) set(key K, value V) {
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

func (wm *WeakMap[K, V]) foreach(fn func(key K, value V) bool) {
	for key, value := range wm.data {
		if !fn(key, value) {
			break
		}
	}
}

func (wm *WeakMap[K, V]) lock() {
	wm.mu.Lock()
}

func (wm *WeakMap[K, V]) unlock() {
	wm.mu.Unlock()
}

func (wm *WeakMap[K, V]) rlock() {
	wm.mu.RLock()
}

func (wm *WeakMap[K, V]) runlock() {
	wm.mu.RUnlock()
}

func (wm *WeakMap[K, V]) MarshalJSON() ([]byte, error) {
	return marshalMap[K, V](wm)
}

func (wm *WeakMap[K, V]) UnmarshalJSON(data []byte) error {
	return unmarshalMap[K, V](wm, data)
}
