package stlx

import "encoding/json"

func (os *OrderedSet[T]) add(element T) {
	os.mp.set(element, struct{}{})
}

func (os *OrderedSet[T]) clear() {
	os.mp.clear()
}

func (os *OrderedSet[T]) vals() []T {
	keys := make([]T, len(os.mp.keys))
	copy(keys, os.mp.keys)
	return keys
}

func (os *OrderedSet[T]) foreach(fn func(element T) bool) {
	os.mp.foreach(func(key T, value void) bool {
		return fn(key)
	})
}

func (os *OrderedSet[T]) lock() {
	os.mp.mu.Lock()
}

func (os *OrderedSet[T]) unlock() {
	os.mp.mu.Unlock()
}

func (os *OrderedSet[T]) rlock() {
	os.mp.mu.RLock()
}

func (os *OrderedSet[T]) runlock() {
	os.mp.mu.RUnlock()
}

// MarshalJSON 实现json.Marshaler接口
func (os *OrderedSet[T]) MarshalJSON() ([]byte, error) {
	os.lock()
	defer os.unlock()
	return json.Marshal(os.mp.keys)
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (os *OrderedSet[T]) UnmarshalJSON(data []byte) error {
	return unmarshalCollection[T](os, data)
}
