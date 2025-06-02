package stlx

import (
	"encoding/json"
	"log"

	"github.com/llyb120/yoya/internal"
)

type Object[T any] interface {
	Put()
}

type object[T any] struct {
	mu      lock
	objects []any
	mp      map[string]T // 缓存
}

func NewObject[T any]() *object[T] {
	return &object[T]{
		objects: make([]any, 4),
		mp:      make(map[string]T),
	}
}

func NewSyncObject[T any]() *object[T] {
	o := NewObject[T]()
	o.mu.sync = true
	return o
}

func (o *object[T]) Put(v any) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.objects = append(o.objects, v)
}

func (o *object[T]) Get(key string) (T, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.convert()
	v, ok := o.mp[key]
	return v, ok
}

func (o *object[T]) ToMap() map[string]T {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.convert()
	return o.mp
}

func (o *object[T]) convert() {
	if len(o.objects) == 0 {
		return
	}
	for _, v := range o.objects {
		var innerMap map[string]T
		err := internal.Cast(&innerMap, v)
		if err != nil {
			log.Println("cast error", err)
			continue
		}
		for k, v := range innerMap {
			o.mp[k] = v
		}
	}
	o.objects = nil
}

func (o *object[T]) MarshalJSON() ([]byte, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.convert()
	return json.Marshal(o.mp)
}
