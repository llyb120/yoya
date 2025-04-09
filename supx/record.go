package supx

import (
	"bytes"
	"encoding/json"
	"reflect"
	"sync"
)

type Record[T any] struct {
	Data       T
	ExtMap     map[string]any
	ExtObjects []any
	once       sync.Once
}

func NewRecord[T any](data T) *Record[T] {
	return &Record[T]{
		Data:       data,
		ExtMap:     make(map[string]any),
		ExtObjects: make([]any, 0, 4),
	}
}

func (r *Record[T]) Put(key string, value any) {
	r.init()
	r.ExtMap[key] = value
}

func (r *Record[T]) PutMap(m map[string]any) {
	r.init()
	for k, v := range m {
		r.ExtMap[k] = v
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
		r.ExtMap = make(map[string]any)
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
		if len(r.ExtMap) > 0 {
			mapBs, err := defaultJSONEncoder(r.ExtMap)
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
