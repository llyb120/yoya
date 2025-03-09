package stlx

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func marshalMap[K any, V any](mp jsonMap[K, V]) ([]byte, error) {
	mp.lock()
	defer mp.unlock()

	var buf bytes.Buffer
	buf.WriteByte('{')

	var reultErr error
	first := true
	mp.foreach(func(key K, value V) bool {
		if !first {
			buf.WriteByte(',')
		}
		first = false

		// 序列化键
		keyBytes, err := json.Marshal(key)
		if err != nil {
			reultErr = err
			return false
		}
		buf.Write(keyBytes)

		buf.WriteByte(':')

		// 序列化值
		valueBytes, err := json.Marshal(value)
		if err != nil {
			reultErr = err
			return false
		}
		buf.Write(valueBytes)
		return true
	})

	buf.WriteByte('}')
	return buf.Bytes(), reultErr
}

// UnmarshalJSON 实现json.Unmarshaler接口
func unmarshalMap[K any, V any](mp jsonMap[K, V], data []byte) error {
	mp.lock()
	defer mp.unlock()

	// 清空现有数据
	mp.clear()

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
		mp.set(key, value)
	}

	// 确保结束是一个对象
	if t, err := dec.Token(); err != nil {
		return err
	} else if delim, ok := t.(json.Delim); !ok || delim != '}' {
		return fmt.Errorf("expected }, got %v", t)
	}

	return nil
}

// MarshalJSON 实现json.Marshaler接口
func marshalCollection[T any](col jsonCollection[T]) ([]byte, error) {
	col.rlock()
	defer col.runlock()
	return json.Marshal(col.vals())
}

// UnmarshalJSON 实现json.Unmarshaler接口
func unmarshalCollection[T any](col jsonCollection[T], data []byte) error {
	col.lock()
	defer col.unlock()

	col.clear()

	var elements []T
	if err := json.Unmarshal(data, &elements); err != nil {
		return err
	}
	for _, element := range elements {
		col.add(element)
	}

	return nil
}
