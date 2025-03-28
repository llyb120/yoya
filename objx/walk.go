package objx

import (
	"runtime"

	reflect "github.com/goccy/go-reflect"
	"github.com/llyb120/yoya/refx"
	"github.com/llyb120/yoya/syncx"
)

// 这里的空对象是为了对齐内存，不能使用struct{}，这样指针会指向相同的部分从而导致指令出错
type walkCommand struct {
	flag bool
	_    struct{}
}

// 值未发生变化
var Unchanged = &walkCommand{}

// 停止遍历
var BreakWalk = &walkCommand{}

// 停止遍历子元素
var BreakWalkSelf = &walkCommand{}

type walkFunc = func(s any, k any, v any) any
type asyncWalkFunc = func(s any, k any, v any) syncx.AsyncFn

// 遍历任意对象
// 因为map和字段的问题，遍历的顺序无法预测，但从外到内可以保证(先序遍历)
//
// 遍历函数使用三个参数，分别是
//  1. 当前遍历的父元素
//  2. 当前遍历的key
//  3. 当前遍历的value
//
// 如果遍历函数返回nil，则不进行任何操作
// 如果遍历函数返回Unchanged，则不进行任何操作
//
// 流程控制语句：
//
//	如果遍历函数返回BreakWalk，则停止遍历
//	如果遍历函数返回BreakWalkSelf，则停止遍历当前元素
//	因异步执行时的调用顺序无法保证，所以异步函数不支持流程控制
//
// 异步遍历：
//
//	如果函数定义为asyncWalkFunc，则遍历函数会异步执行
//
// 除此之外的任何返回值都会被设置到当前遍历的元素上
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
	for v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem()
	}
	var ref reflect.Value
	if v.CanAddr() {
		ref = v.Addr()
	} else {
		ref = reflect.New(v.Type())
		ref.Elem().Set(v)
	}
	switch v.Kind() {
	case reflect.Map:
		for _, k := range v.MapKeys() {
			kk := k.Interface()
			vv := v.MapIndex(k)
			var res any
			if wg != nil {
				wg.Go(func() {
					f := any(fn).(asyncWalkFunc)
					res = f(ref.Interface(), kk, vv.Interface())
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
				res = f(ref.Interface(), kk, vv.Interface())
				if res != Unchanged && res != nil {
					refx.UnsafeSetFieldValue(vv, reflect.ValueOf(res), true)
				}
			}
			if res == BreakWalkSelf {
				continue
			}
			if res == BreakWalk {
				return
			}
			walk(vv, fn, wg)
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			i := i
			vv := v.Index(i)
			var res any
			if wg != nil {
				wg.Go(func() {
					f := any(fn).(asyncWalkFunc)
					res = f(ref.Interface(), i, vv.Interface())
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
				res = f(ref.Interface(), i, vv.Interface())
				if res != Unchanged && res != nil {
					refx.UnsafeSetFieldValue(vv, reflect.ValueOf(res), true)
				}
			}
			if res == BreakWalkSelf {
				continue
			}
			if res == BreakWalk {
				return
			}
			walk(vv, fn, wg)
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			i := i
			vv := v.Field(i)
			var res any
			if wg != nil {
				wg.Go(func() {
					f := any(fn).(asyncWalkFunc)
					res = f(ref.Interface(), v.Type().Field(i).Name, vv.Interface())
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
				res = f(ref.Interface(), v.Type().Field(i).Name, vv.Interface())
				if res != Unchanged && res != nil {
					// 可以设置才设置
					refx.UnsafeSetFieldValue(vv, reflect.ValueOf(res), true)
				}
			}
			if res == BreakWalkSelf {
				continue
			}
			if res == BreakWalk {
				return
			}
			walk(vv, fn, wg)
		}
	}
}
