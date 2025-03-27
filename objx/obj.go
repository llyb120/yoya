package objx

import (
	"fmt"
	"runtime"

	reflect "github.com/goccy/go-reflect"
	"github.com/llyb120/yoya/refx"
	"github.com/llyb120/yoya/syncx"
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

var Unchanged = &struct{}{}

type walkFunc = func(s any, k any, v any) any
type asyncWalkFunc = func(s any, k any, v any) syncx.AsyncFn

func Walk[T walkFunc | asyncWalkFunc](dest any, fn T) {
	var isAsync bool
	if _, ok := any(fn).(asyncWalkFunc); ok {
		isAsync = true
	}
	var g *walkPool
	if isAsync {
		g = newWalkPool(runtime.GOMAXPROCS(0))
		defer g.Destroy()
	}
	walk(dest, fn, g)
	if isAsync {
		g.Wait()
	}
}

func walk[T any](dest any, fn T, wg *walkPool) {
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
			if wg != nil {
				wg.Go(func() {
					f := any(fn).(asyncWalkFunc)
					var ref reflect.Value
					if v.CanAddr() {
						ref = v.Addr()
					} else {
						ref = reflect.New(v.Type())
						ref.Elem().Set(v)
					}
					res := f(ref, kk, vv.Interface())
					err := syncx.Await(res)
					if err != nil {
						return
					}
					if res != Unchanged && res != nil {
						wg.Lock()
						defer wg.Unlock()
						refx.UnsafeSetFieldValue(vv, reflect.ValueOf(res).Elem(), true)
					}
				})
			} else {
				f := any(fn).(walkFunc)
				var ref reflect.Value
				if v.CanAddr() {
					ref = v.Addr()
				} else {
					ref = reflect.New(v.Type())
					ref.Elem().Set(v)
				}
				res := f(ref, kk, vv.Interface())
				if res != Unchanged && res != nil {
					refx.UnsafeSetFieldValue(vv, reflect.ValueOf(res), true)
				}
			}
			walk(vv, fn, wg)
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			i := i
			vv := v.Index(i)
			if wg != nil {
				wg.Go(func() {
					f := any(fn).(asyncWalkFunc)
					var ref reflect.Value
					if v.CanAddr() {
						ref = v.Addr()
					} else {
						ref = reflect.New(v.Type())
						ref.Elem().Set(v)
					}
					res := f(ref, i, vv.Interface())
					err := syncx.Await(res)
					if err != nil {
						return
					}
					if res != Unchanged && res != nil {
						wg.Lock()
						defer wg.Unlock()
						refx.UnsafeSetFieldValue(vv, reflect.ValueOf(res).Elem(), true)
					}
				})
			} else {
				f := any(fn).(walkFunc)
				var ref reflect.Value
				if v.CanAddr() {
					ref = v.Addr()
				} else {
					ref = reflect.New(v.Type())
					ref.Elem().Set(v)
				}
				res := f(ref, i, vv.Interface())
				if res != Unchanged && res != nil {
					refx.UnsafeSetFieldValue(vv, reflect.ValueOf(res), true)
				}
			}
			walk(vv, fn, wg)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			i := i
			vv := v.Field(i)
			if wg != nil {
				wg.Go(func() {
					f := any(fn).(asyncWalkFunc)
					var ref reflect.Value
					if v.CanAddr() {
						ref = v.Addr()
					} else {
						ref = reflect.New(v.Type())
						ref.Elem().Set(v)
					}
					res := f(ref, v.Type().Field(i).Name, vv.Interface())
					err := syncx.Await(res)
					if err != nil {
						return
					}
					if res != Unchanged && res != nil {
						wg.Lock()
						defer wg.Unlock()
						refx.UnsafeSetFieldValue(vv, reflect.ValueOf(res).Elem(), true)
					}
				})
			} else {
				f := any(fn).(walkFunc)
				var ref reflect.Value
				if v.CanAddr() {
					ref = v.Addr()
				} else {
					ref = reflect.New(v.Type())
					ref.Elem().Set(v)
				}
				res := f(ref, v.Type().Field(i).Name, vv.Interface())
				if res != Unchanged && res != nil {
					// 可以设置才设置
					refx.UnsafeSetFieldValue(vv, reflect.ValueOf(res), true)
				}
			}
			walk(vv, fn, wg)
		}
	}
}
