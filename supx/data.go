package supx

import (
	"reflect"
	"strings"
	"sync"
	_ "unsafe"
)

const dataKey = "$data"

type Data[T any] map[string]any

var fieldCache sync.Map // map[reflect.Type]map[string]struct{}

func cachedFieldSet(t reflect.Type) map[string]struct{} {
	// 处理指针类型
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if v, ok := fieldCache.Load(t); ok {
		return v.(map[string]struct{})
	}
	fields := make(map[string]struct{})
	if t.Kind() == reflect.Struct {
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			if !f.IsExported() {
				continue
			}
			name := f.Tag.Get("json")
			if name == "" {
				name = f.Name
			} else {
				// 处理如 "id,omitempty" 的情况
				if idx := strings.IndexByte(name, ','); idx != -1 {
					name = name[:idx]
				}
			}
			if name == "-" || name == "" {
				continue
			}
			fields[name] = struct{}{}
		}
	}
	fieldCache.Store(t, fields)
	return fields
}

func NewData[T any]() Data[T] {
	return make(Data[T])
}

func (d Data[T]) Data() *T {
	if data, ok := d[dataKey]; ok {
		return data.(*T)
	}
	return nil
}

func (d Data[T]) Set(data T) {
	d[dataKey] = &data
}

func (r Data[T]) GetType() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func (d Data[T]) Clone() Data[T] {
	cp, err := DeepCloneAny(d)
	if err != nil {
		return NewData[T]()
	}
	item, ok := cp.(Data[T])
	if !ok {
		return NewData[T]()
	}
	return item
}

// MarshalJSON 实现JSON序列化，优先序列化$data，然后补充其他属性
func (d Data[T]) MarshalJSON() ([]byte, error) {
	var dataJSON string
	var fieldSet map[string]struct{}

	// 处理$data
	if dataPtr := d.Data(); dataPtr != nil {
		dataBytes, err := defaultJSONEncoder(*dataPtr)
		if err != nil {
			return nil, err
		}
		dataJSON = string(dataBytes)
		fieldSet = cachedFieldSet(reflect.TypeOf(*dataPtr))
	} else {
		// data 为空也要获取类型字段缓存
		var zero T
		fieldSet = cachedFieldSet(reflect.TypeOf(zero))
	}

	// 处理额外字段
	extraFields := make(map[string]interface{})
	for key, val := range d {
		if key == dataKey {
			continue
		}
		if _, inData := fieldSet[key]; !inData {
			extraFields[key] = val
		}
	}

	// 如果$data不是JSON对象（例如字符串、数字等），直接构造map进行序列化
	if dataJSON != "" && (dataJSON[0] != '{') {
		result := make(map[string]interface{})
		if dataPtr := d.Data(); dataPtr != nil {
			result[dataKey] = *dataPtr
		}
		for k, v := range extraFields {
			result[k] = v
		}
		return defaultJSONEncoder(result)
	}

	// 合并JSON
	merged, err := mergeJSON(dataJSON, extraFields)
	if err != nil {
		return nil, err
	}
	return merged, nil
}

// 合并$data JSON和额外字段的助手
func mergeJSON(dataJSON string, extraFields map[string]interface{}) ([]byte, error) {
	// 如果没有$data，直接序列化额外字段或返回空对象
	if dataJSON == "" {
		if extraFields == nil || len(extraFields) == 0 {
			return []byte("{}"), nil
		}
		return defaultJSONEncoder(extraFields)
	}

	// 如果没有额外字段
	if extraFields == nil || len(extraFields) == 0 {
		return []byte(dataJSON), nil
	}

	// 序列化额外字段
	extraBytes, err := defaultJSONEncoder(extraFields)
	if err != nil {
		return nil, err
	}
	extraJSON := string(extraBytes)

	// 如果其中一个是空对象，直接返回另一个
	if dataJSON == "{}" {
		return []byte(extraJSON), nil
	}
	if extraJSON == "{}" {
		return []byte(dataJSON), nil
	}

	// 拼接：去掉dataJSON最后一个}和extraJSON第一个{，中间用逗号连接
	merged := dataJSON[:len(dataJSON)-1] + "," + extraJSON[1:]
	return []byte(merged), nil
}

// UnmarshalJSON 实现JSON反序列化
func (d *Data[T]) UnmarshalJSON(data []byte) error {
	if *d == nil {
		*d = NewData[T]()
	}

	// 解析到临时map，便于分类字段
	var raw map[string]interface{}
	if err := defaultJSONDecoder(data, &raw); err != nil {
		return err
	}

	// 获取字段缓存
	var zero T
	fieldSet := cachedFieldSet(reflect.TypeOf(zero))

	// 反序列化到实际类型（忽略错误，不影响额外字段处理）
	var typed T
	_ = defaultJSONDecoder(data, &typed)
	(*d)[dataKey] = &typed // 始终保证存在$data，并非nil

	// 分类额外字段
	for k, v := range raw {
		if k == dataKey {
			continue
		}
		if _, exists := fieldSet[k]; !exists {
			(*d)[k] = v
		}
	}

	return nil
}

func Get[T any](d map[string]any, key string) T {
	if data, ok := d[key]; ok {
		v, ok := data.(T)
		if !ok {
			return *new(T)
		}
		return v
	}

	return *new(T)
}

//go:linkname DeepCloneAny github.com/llyb120/yoya/objx.DeepCloneAny
func DeepCloneAny(src any) (any, error)
