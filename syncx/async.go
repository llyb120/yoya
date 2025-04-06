package syncx

import (
	"context"
	"fmt"

	"sync"
	"sync/atomic"
	"time"

	"reflect"

	rf "github.com/goccy/go-reflect"
)

type future struct {
	exprtime time.Time
	result   any
	err      error
	wg       sync.WaitGroup
	done     atomic.Bool
}

func (f *future) Get(timeout time.Duration) (any, error) {
	if f.done.Load() {
		return f.result, f.err
	}

	// 等待结果或超时
	if timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// 创建一个完成信号通道
		done := make(chan struct{})

		go func() {
			f.wg.Wait()
			close(done)
		}()

		// 等待结果或超时
		select {
		case <-done:
			// 任务已完成，返回结果
			return f.result, f.err
		case <-ctx.Done():
			// 超时了，返回零值和超时错误
			return nil, fmt.Errorf("获取结果超时: %w", ctx.Err())
		}
	} else {
		// 没有超时，等待结果
		f.wg.Wait()
		if f.err != nil {
			return f.result, f.err
		}
	}

	return f.result, nil
}

func Async[T any](fn any) func(...any) *T {
	fv := rf.ValueOf(fn)
	ft := fv.Type()
	return func(args ...any) (ptrResult *T) {
		future := &future{exprtime: time.Now().Add(5 * time.Minute)}
		var zero T
		ptrResult = &zero
		future.wg.Add(1)
		saveFuture(ptrResult, future)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					future.err = fmt.Errorf("future panic: %v", r)
				}
				future.wg.Done()
				future.done.Store(true)
			}()

			in := make([]rf.Value, len(args))
			for i, arg := range args {
				in[i] = rf.ValueOf(arg)
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

			result := fv.Call(in)
			if len(result) > 0 {
				for _, r := range result {
					val := r.Interface()
					if err, ok := val.(error); ok {
						future.err = err
					} else if value, ok := val.(T); ok {
						*ptrResult = value
						future.result = value
					}
				}
			}
		}()

		return ptrResult
	}
}

func AsyncReflect(fn reflect.Value, outType reflect.Type) func(...any) any {
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
							} else if underlyingValue.CanConvert(outType) {
								newInstance.Set(underlyingValue.Convert(outType))
							} else {
								future.err = fmt.Errorf("无法将类型 %v 转换为 %v", underlyingValue.Type(), outType)
								return
							}
						} else if r.Type().AssignableTo(outType) {
							// 直接赋值
							newInstance.Set(r)
						} else if r.CanConvert(outType) {
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

// func CanAwait(obj any) bool {
// 	return futureHolder.contains(obj)
// }

func Await(objs ...any) error {
	var timeout time.Duration = 0

	// 检查最后一个参数是否为超时时间
	if len(objs) > 0 {
		var ok bool
		if timeout, ok = objs[len(objs)-1].(time.Duration); ok {
			objs = objs[:len(objs)-1]
		}
	}

	// 如果有需要展开的
	var futures []*future = make([]*future, 0, len(objs))
	for _, e := range objs {
		// 如果是数组，展开
		tp := rf.TypeOf(e)
		if tp == nil {
			continue
		}
		if tp.Kind() == rf.Array || tp.Kind() == rf.Slice {
			val := rf.ValueOf(e)
			for i := 0; i < val.Len(); i++ {
				f := loadFuture(val.Index(i).Interface())
				if f != nil {
					futures = append(futures, f)
				}
			}
			continue
		}
		// 如果是map
		if tp.Kind() == rf.Map {
			val := rf.ValueOf(e)
			for _, vk := range val.MapKeys() {
				f := loadFuture(val.MapIndex(vk).Interface())
				if f != nil {
					futures = append(futures, f)
				}
			}
			continue
		}
		// 其余的情况
		f := loadFuture(e)
		if f != nil {
			futures = append(futures, f)
		}
	}

	// 如果没有参数，等待所有
	if len(futures) == 0 {
		return nil
	}

	if len(futures) > 1 {
		var g Group
		for _, f := range futures {
			f := f
			g.Go(func() error {
				_, err := f.Get(timeout)
				return err
			})
		}
		return g.Wait(timeout)

	} else if len(futures) == 1 {
		f := futures[0]
		_, err := f.Get(timeout)
		return err
	}

	return nil
}
