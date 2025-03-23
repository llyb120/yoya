package stlx

import "encoding/json"

func (os *orderedSet[T]) add(element T) {
	os.mp.set(element, struct{}{})
}

func (os *orderedSet[T]) addAll(elements []T) {
	for _, element := range elements {
		os.add(element)
	}
}

func (os *orderedSet[T]) clear() {
	os.mp.clear()
}

func (os *orderedSet[T]) vals() []T {
	keys := make([]T, len(os.mp.keys))
	copy(keys, os.mp.keys)
	return keys
}

func (os *orderedSet[T]) foreach(fn func(element T) bool) {
	os.mp.foreach(func(key T, value void) bool {
		return fn(key)
	})
}

func (os *orderedSet[T]) lock() {
	os.mp.mu.Lock()
}

func (os *orderedSet[T]) unlock() {
	os.mp.mu.Unlock()
}

func (os *orderedSet[T]) rlock() {
	os.mp.mu.RLock()
}

func (os *orderedSet[T]) runlock() {
	os.mp.mu.RUnlock()
}

// MarshalJSON 实现json.Marshaler接口
func (os *orderedSet[T]) MarshalJSON() ([]byte, error) {
	os.lock()
	defer os.unlock()
	return json.Marshal(os.mp.keys)
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (os *orderedSet[T]) UnmarshalJSON(data []byte) error {
	return unmarshalCollection[T](os, data)
}
