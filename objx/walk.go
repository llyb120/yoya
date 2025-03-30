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

// type walkFunc = func(s any, k any, v any) any
// type asyncWalkFunc = func(s any, k any, v any) syncx.AsyncFn

type walkOption int

var (
	Async walkOption = -1
	Level walkOption = 1
)

func test() {
	Walk(&struct{}{}, func(s, k, v any) any {
		return nil
	}, 10*Level)
}

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
func Walk(dest any, fn func(s any, k any, v any) any, opts ...walkOption) {
	var walkCtx = &walkContext{
		fn: fn,
	}
	for _, opt := range opts {
		if opt == Async {
			walkCtx.isAsync = true
		}
		if opt > 0 {
			walkCtx.level = int(opt)
		}
	}
	if walkCtx.isAsync {
		walkCtx.wg = newWalkPool(runtime.GOMAXPROCS(0))
		defer walkCtx.wg.Destroy()
	}
	walkCtx.walk(dest, 0)
	if walkCtx.isAsync {
		walkCtx.wg.Wait()
	}
}

type walkContext struct {
	level   int
	isAsync bool
	wg      *walkPool
	fn      func(s any, k any, v any) any
}

func (w *walkContext) doFunc(ref reflect.Value, k any, v reflect.Value) any {
	var res any
	if kk, ok := k.(reflect.Value); ok {
		k = kk.Interface()
	}
	if w.isAsync {
		w.wg.Go(func() {
			res = w.fn(ref.Interface(), k, v.Interface())
			err := syncx.Await(res)
			if err != nil {
				return
			}
			if res != Unchanged && res != nil {
				w.wg.Lock()
				defer w.wg.Unlock()
				refx.UnsafeSetFieldValue(v, reflect.ValueOf(res).Elem(), true)
			}
		})
	} else {
		res = w.fn(ref.Interface(), k, v.Interface())
		if res != Unchanged && res != nil {
			refx.UnsafeSetFieldValue(v, reflect.ValueOf(res), true)
		}
	}
	return res
}

func (w *walkContext) walk(dest any, level int) {
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
			vv := v.MapIndex(k)
			res := w.doFunc(ref, k, vv)
			if res == BreakWalkSelf {
				continue
			}
			if res == BreakWalk {
				return
			}
			if w.level <= 0 || level+1 <= w.level {
				w.walk(vv, level+1)
			}
		}
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			i := i
			vv := v.Index(i)
			res := w.doFunc(ref, i, vv)
			if res == BreakWalkSelf {
				continue
			}
			if res == BreakWalk {
				return
			}
			if w.level <= 0 || level+1 <= w.level {
				w.walk(vv, level+1)
			}
		}
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			i := i
			vv := v.Field(i)
			k := v.Type().Field(i).Name
			res := w.doFunc(ref, k, vv)
			if res == BreakWalkSelf {
				continue
			}
			if res == BreakWalk {
				return
			}
			if w.level <= 0 || level+1 <= w.level {
				w.walk(vv, level+1)
			}
		}
	}
}
