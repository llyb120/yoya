package syncx

import (
	"fmt"
	"reflect"
	"time"
)

// 高性能的reflect
func AsyncHighReflect(fn reflect.Value, outType reflect.Type) func(...any) any {
	ft := fn.Type()
	return func(args ...any) (ptrResult any) {
		future := &future{exprtime: time.Now().Add(5 * time.Minute)}
		// var zero = reflect.New(outType).Interface()
		ptrRef := reflect.New(outType)
		ptrResult = ptrRef.Interface()
		future.wg.Add(1)
		saveFuture(ptrResult, future)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					future.err = fmt.Errorf("future panic: %v", r)
				}
				future.done.Store(true)
				future.wg.Done()
				deleteFuture(ptrResult)
			}()

			in := make([]reflect.Value, len(args))
			for i, arg := range args {
				in[i] = reflect.ValueOf(arg)
			}

			// 类型检查
			if len(args) != ft.NumIn() {
				future.err = fmt.Errorf("参数数量不匹配")
				return
			}
			for i := 0; i < ft.NumIn(); i++ {
				if i >= len(in) || !in[i].Type().AssignableTo(ft.In(i)) {
					future.err = fmt.Errorf("参数类型不匹配")
					return
				}
			}

			result := fn.Call(in)
			if len(result) > 0 {
				for i, r := range result {
					val := r.Interface()
					if i == 0 {
						if val == nil {
							future.result = nil
							continue
						}
						// 创建一个outType类型的新实例
						newInstance := reflect.New(outType).Elem()

						// 尝试从r中提取值并设置到newInstance中
						if r.Kind() == reflect.Interface {
							// 如果r是接口类型，获取其底层值
							underlyingValue := r.Elem()
							if underlyingValue.Type().AssignableTo(outType) {
								newInstance.Set(underlyingValue)
							} else if canConvert(underlyingValue, outType) {
								newInstance.Set(underlyingValue.Convert(outType))
							} else {
								future.err = fmt.Errorf("无法将类型 %v 转换为 %v", underlyingValue.Type(), outType)
								return
							}
						} else if r.Type().AssignableTo(outType) {
							// 直接赋值
							newInstance.Set(r)
						} else if canConvert(r, outType) {
							// 尝试类型转换
							newInstance.Set(r.Convert(outType))
						} else {
							future.err = fmt.Errorf("无法将类型 %v 转换为 %v", r.Type(), outType)
							return
						}

						// 设置到ptrRef和future.result
						ptrRef.Elem().Set(newInstance)
						future.result = newInstance.Interface()
					}
					if i == 1 {
						if err, ok := val.(error); ok || err == nil {
							future.err = err
						} else {
							future.err = fmt.Errorf("返回值类型不匹配")
						}
					}
				}
			}
		}()

		return ptrResult
	}
}

func canConvert(v reflect.Value, t reflect.Type) bool {
	vt := v.Type()
	if !vt.ConvertibleTo(t) {
		return false
	}
	// Converting from slice to array or to pointer-to-array can panic
	// depending on the value.
	switch {
	case vt.Kind() == reflect.Slice && t.Kind() == reflect.Array:
		if t.Len() > v.Len() {
			return false
		}
	case vt.Kind() == reflect.Slice && t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Array:
		n := t.Elem().Len()
		if n > v.Len() {
			return false
		}
	}
	return true
}
