package objx

import (
	"fmt"
	"github.com/llyb120/yoya/syncx"
	"runtime"

	reflect "github.com/goccy/go-reflect"
)

func Assign[T comparable, K0 any, K1 any](dest map[T]K0, source map[T]K1) {
	var zero0 K0
	var zero1 K1
	isSameType := reflect.TypeOf(zero0) == reflect.TypeOf(zero1)
	if isSameType {
		for k, v := range source {
			dest[k] = any(v).(K0)
		}
		return
	} else {
		// 一次性转过来
		var mp map[T]K0
		if err := Cast(&mp, source); err != nil {
			fmt.Println("err:", err)
			return
		}
		for k, v := range mp {
			dest[k] = v
		}
		return
	}
}

type walkFunc = func(k any, v any) any

var Unchanged = &struct{}{}

func Walk(dest any, fn walkFunc) {
	var g = newWalkPool(runtime.GOMAXPROCS(0))
	defer g.Destroy()
	walk(dest, fn, g)
	g.Wait()
}

func walk(dest any, fn walkFunc, wg *walkPool) {
	var v reflect.Value
	var ok bool
	if v, ok = dest.(reflect.Value); !ok {
		v = reflect.ValueOf(dest)
	}
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Map:
		for _, k := range v.MapKeys() {
			kk := k.Interface()
			vv := v.MapIndex(k)
			res := fn(kk, vv.Interface())
			canAwait := syncx.CanAwait(res)
			wg.Go(func() {
				_ = syncx.Await(res)
				if res != Unchanged && res != nil {
					wg.Lock()
					defer wg.Unlock()
					if canAwait {
						vv.Set(reflect.ValueOf(res).Elem())
					} else {
						vv.Set(reflect.ValueOf(res))
					}
				}
			})
			walk(vv, fn, wg)
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			i := i
			vv := v.Index(i)
			res := fn(i, vv.Interface())
			canAwait := syncx.CanAwait(res)
			wg.Go(func() {
				_ = syncx.Await(res)
				if res != Unchanged && res != nil {
					wg.Lock()
					defer wg.Unlock()
					if canAwait {
						vv.Set(reflect.ValueOf(res).Elem())
					} else {
						vv.Set(reflect.ValueOf(res))
					}
				}
			})
			walk(vv, fn, wg)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			i := i
			vv := v.Field(i)
			res := fn(v.Type().Field(i).Name, vv.Interface())
			canAwait := syncx.CanAwait(res)
			wg.Go(func() {
				_ = syncx.Await(res)
				if res != Unchanged && res != nil {
					wg.Lock()
					defer wg.Unlock()
					if canAwait {
						vv.Set(reflect.ValueOf(res).Elem())
					} else {
						vv.Set(reflect.ValueOf(res))
					}
				}
			})
			walk(vv, fn, wg)
		}
	}
}
