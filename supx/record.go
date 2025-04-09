package supx

import (
	"bytes"
	"encoding/json"
	"sync"
)

var (
	ObjectKey = "__OBJECT__"
)

type Record[T any] interface {
	Put(key string, value any)
	Get(key ...string) any
	MarshalJSON() ([]byte, error)
}

type record[T any] struct {
	data         T
	extraMap     map[string]any
	extraObjects []any
}

func NewRecord[T any](data T) Record[T] {
	return &record[T]{
		data:         data,
		extraMap:     make(map[string]any),
		extraObjects: make([]any, 0, 4),
	}
}

func (r *record[T]) Put(key string, value any) {
	if key == ObjectKey {
		r.extraObjects = append(r.extraObjects, value)
	} else {
		r.extraMap[key] = value
	}
}

func (r *record[T]) Get(key ...string) any {
	if len(key) == 0 {
		return r.data
	}
	if key[0] == ObjectKey {
		return r.extraObjects
	}
	return r.extraMap[key[0]]
}

func (r record[T]) MarshalJSON() ([]byte, error) {
	buff := bufferPool.Get().(*bytes.Buffer)
	defer func() {
		buff.Reset()
		bufferPool.Put(buff)
	}()
	// 构造主数据
	mainBs, err := defaultJSONEncoder(r.data)
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
		if len(r.extraMap) > 0 {
			mapBs, err := defaultJSONEncoder(r.extraMap)
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
		if len(r.extraObjects) > 0 {
			for _, v := range r.extraObjects {
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
