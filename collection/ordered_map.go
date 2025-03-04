package collection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sync"
)

// OrderedMap 是一个协程安全的有序映射，按插入顺序维护键值对
type OrderedMap[K comparable, V any] struct {
	mu      sync.RWMutex
	keys    []K
	values  []V
	indexes map[K]int
}

// NewOrderedMap 创建一个新的有序映射
func NewOrderedMap[K comparable, V any]() Map[K, V] {
	return &OrderedMap[K, V]{
		indexes: make(map[K]int),
	}
}

// Set 添加或更新键值对
func (om *OrderedMap[K, V]) Set(key K, value V) {
	om.mu.Lock()
	defer om.mu.Unlock()

	om.doSet(key, value)
}

func (om *OrderedMap[K, V]) doSet(key K, value V) {
	if index, exists := om.indexes[key]; exists {
		// 如果键已存在，只更新值
		om.values[index] = value
		return
	}

	// 添加到映射
	om.keys = append(om.keys, key)
	om.values = append(om.values, value)
	om.indexes[key] = len(om.keys) - 1
}

// Get 获取键对应的值
func (om *OrderedMap[K, V]) Get(key K) (V, bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	if index, exists := om.indexes[key]; exists {
		return om.values[index], true
	}

	var zero V
	return zero, false
}

// Del 删除键值对
func (om *OrderedMap[K, V]) Del(key K) V {
	om.mu.Lock()
	defer om.mu.Unlock()

	pos, exists := om.indexes[key]
	if !exists {
		var zero V
		return zero
	} else {
		delete(om.indexes, key)
		val := om.values[pos]
		om.keys = append(om.keys[:pos], om.keys[pos+1:]...)
		om.values = append(om.values[:pos], om.values[pos+1:]...)
		return val
	}
}

// Size 返回映射大小
func (om *OrderedMap[K, V]) Size() int {
	om.mu.RLock()
	defer om.mu.RUnlock()
	return len(om.keys)
}

// Keys 按插入顺序返回所有键
func (om *OrderedMap[K, V]) Keys() []K {
	om.mu.RLock()
	defer om.mu.RUnlock()

	keys := make([]K, 0, len(om.keys))
	copy(keys, om.keys)
	return keys
}

// Vals 按插入顺序返回所有值
func (om *OrderedMap[K, V]) Vals() []V {
	om.mu.RLock()
	defer om.mu.RUnlock()

	values := make([]V, 0, len(om.values))
	copy(values, om.values)
	return values
}

// Clear 清空映射
func (om *OrderedMap[K, V]) Clear() {
	om.mu.Lock()
	defer om.mu.Unlock()

	om.doClear()
}

func (om *OrderedMap[K, V]) doClear() {
	om.keys = nil
	om.values = nil
	om.indexes = make(map[K]int)
}

// MarshalJSON 实现json.Marshaler接口
func (om *OrderedMap[K, V]) MarshalJSON() ([]byte, error) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	var buf bytes.Buffer
	buf.WriteByte('{')

	first := true
	for i, key := range om.keys {
		if !first {
			buf.WriteByte(',')
		}
		first = false

		// 序列化键
		keyBytes, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		buf.Write(keyBytes)

		buf.WriteByte(':')

		// 序列化值
		valueBytes, err := json.Marshal(om.values[i])
		if err != nil {
			return nil, err
		}
		buf.Write(valueBytes)
	}

	buf.WriteByte('}')
	return buf.Bytes(), nil
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (om *OrderedMap[K, V]) UnmarshalJSON(data []byte) error {
	om.mu.Lock()
	defer om.mu.Unlock()

	// 清空现有数据
	om.doClear()

	dec := json.NewDecoder(bytes.NewReader(data))

	// 确保开始是一个对象
	if t, err := dec.Token(); err != nil {
		return err
	} else if delim, ok := t.(json.Delim); !ok || delim != '{' {
		return fmt.Errorf("expected {, got %v", t)
	}

	// 读取键值对
	for dec.More() {
		// 读取键
		var key K
		keyToken, err := dec.Token()
		if err != nil {
			return err
		}

		// 如果键是字符串类型，需要特殊处理
		if keyStr, ok := keyToken.(string); ok {
			if err := json.Unmarshal([]byte(`"`+keyStr+`"`), &key); err != nil {
				return err
			}
		} else {
			if err := json.Unmarshal([]byte(fmt.Sprintf("%v", keyToken)), &key); err != nil {
				return err
			}
		}

		// 读取值
		var value V
		if err := dec.Decode(&value); err != nil {
			return err
		}

		// 添加到有序映射
		om.doSet(key, value)
	}

	// 确保结束是一个对象
	if t, err := dec.Token(); err != nil {
		return err
	} else if delim, ok := t.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expected }, got %v", t)
	}

	return nil
}

// ForEach 按顺序遍历所有键值对
func (om *OrderedMap[K, V]) ForEach(fn func(key K, value V) bool) {
	om.mu.RLock()
	defer om.mu.RUnlock()

	for i, key := range om.keys {
		if !fn(key, om.values[i]) {
			break
		}
	}
}
