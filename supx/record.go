package supx

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sync"
)

type Record[T any] struct {
	Data       T
	Ext        map[string]any
	ExtObjects []any
	once       sync.Once
}

func NewRecord[T any](data T) *Record[T] {
	return &Record[T]{
		Data:       data,
		Ext:        make(map[string]any),
		ExtObjects: make([]any, 0, 4),
	}
}

func (r *Record[T]) Put(key string, value any) {
	r.init()
	r.Ext[key] = value
}

func (r *Record[T]) GetExt(key string) any {
	r.init()
	return r.Ext[key]
}

func (r *Record[T]) PutMap(m map[string]any) {
	r.init()
	for k, v := range m {
		r.Ext[k] = v
	}
}

func (r *Record[T]) PubObject(obj any) {
	r.init()
	r.ExtObjects = append(r.ExtObjects, obj)
}

func (r Record[T]) GetType() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func (r *Record[T]) init() {
	r.once.Do(func() {
		r.Ext = make(map[string]any)
	})
}

func (r *Record[T]) MarshalJSON() ([]byte, error) {
	buff := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buff.Reset()
		bufferPool.Put(buff)
	}()
	// 构造主数据
	mainBs, err := defaultJSONEncoder(r.Data)
	if err != nil {
		return nil, err
	}
	buff.WriteByte('{')
	var hasMain bool
	if len(mainBs) > 2 {
		mainBs = mainBs[1 : len(mainBs)-1]
		buff.Write(mainBs)
		hasMain = true
	}
	if hasMain {
		// 构造额外数据
		if len(r.Ext) > 0 {
			mapBs, err := defaultJSONEncoder(r.Ext)
			if err != nil {
				return nil, err
			}
			if len(mapBs) > 2 {
				mapBs = mapBs[1 : len(mapBs)-1]
				buff.WriteByte(',')
				buff.Write(mapBs)
			}
		}
		// 如果有额外的对象
		if len(r.ExtObjects) > 0 {
			for _, v := range r.ExtObjects {
				objBs, err := defaultJSONEncoder(v)
				if err != nil {
					return nil, err
				}
				if len(objBs) > 2 {
					objBs = objBs[1 : len(objBs)-1]
					buff.WriteByte(',')
					buff.Write(objBs)
				}
			}
		}
	}
	buff.WriteByte('}')
	return buff.Bytes(), nil
}

func (r *Record[T]) UnmarshalJSON(data []byte) error {
	r.init()

	// 解析为临时映射
	var rawMap map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMap); err != nil {
		return err
	}

	// 创建一个新的T类型实例
	var target T

	// 获取目标类型的字段信息
	targetType := reflect.TypeOf(target)
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}

	// 如果是基本类型或非结构体类型，直接解析整个JSON到Data
	if targetType.Kind() != reflect.Struct {
		if err := json.Unmarshal(data, &target); err != nil {
			return err
		}
		r.Data = target
		return nil
	}

	// 创建一个map，用于存储目标类型的字段名到JSON字段名的映射
	fieldMap := make(map[string]string)
	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			fieldMap[field.Name] = field.Name
		} else {
			// 处理json标签，提取字段名
			tagParts := bytes.Split([]byte(jsonTag), []byte(","))
			if len(tagParts[0]) > 0 {
				fieldMap[field.Name] = string(tagParts[0])
			} else {
				fieldMap[field.Name] = field.Name
			}
		}
	}

	// 创建一个只包含目标类型字段的映射
	extData := make(map[string]json.RawMessage)

	// 遍历原始映射，区分哪些字段属于目标类型，哪些是额外字段
	for k, v := range rawMap {
		found := false
		for _, jsonKey := range fieldMap {
			if k == jsonKey {
				found = true
				break
			}
		}
		if !found {
			extData[k] = v
		}
	}

	if err := json.Unmarshal(data, &target); err != nil {
		return err
	}
	r.Data = target

	// 解析额外字段到Ext
	for k, v := range extData {
		var anyValue any
		if err := json.Unmarshal(v, &anyValue); err != nil {
			return err
		}
		r.Ext[k] = anyValue
	}

	return nil
}

type JSONEncoder func(v any) ([]byte, error)

var defaultJSONEncoder JSONEncoder = json.Marshal

func SetJsonEncoder(encoder JSONEncoder) {
	defaultJSONEncoder = encoder
}

// 不能使用这个反序列化

var bufferPool = sync.Pool{
	New: func() any {
		return bytes.NewBuffer(nil)
	},
}
