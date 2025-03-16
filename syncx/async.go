package syncx

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type future[T any] struct {
	exprtime time.Time
	result   T
	err      error
	wg       sync.WaitGroup
	done     atomic.Bool
}

func (f *future[T]) Get(timeout time.Duration) (T, error) {
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
			var zero T
			return zero, fmt.Errorf("获取结果超时: %w", ctx.Err())
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

var futureHolder sync.Map

func Async[T any](fn any) func(...any) *T {
	fv := reflect.ValueOf(fn)
	ft := fv.Type()
	return func(args ...any) *T {
		future := &future[any]{exprtime: time.Now().Add(5 * time.Minute)}
		var zero T
		ptrResult := &zero
		future.wg.Add(1)
		futureHolder.Store(ptrResult, future)

		go func() {
			defer func() {
				if r := recover(); r != nil {
					future.err = fmt.Errorf("future panic: %v", r)
				}
				future.wg.Done()
				future.done.Store(true)
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

func Await(objs ...any) error {
	var timeout time.Duration = 0

	// 检查最后一个参数是否为超时时间
	if len(objs) > 0 {
		var ok bool
		if timeout, ok = objs[len(objs)-1].(time.Duration); ok {
			objs = objs[:len(objs)-1]
		}
	}

	if len(objs) > 1 {
		var g Group
		for _, obj := range objs {
			obj := obj
			g.Go(func() error {
				f, ok := futureHolder.LoadAndDelete(obj)
				if !ok {
					return nil
				}
				_, err := f.(*future[any]).Get(0)
				return err
			})
		}
		return g.Wait(timeout)

	} else if len(objs) == 1 {
		f, ok := futureHolder.LoadAndDelete(objs[0])
		if !ok {
			return nil
		}
		_, err := f.(*future[any]).Get(timeout)
		return err
	}
	return nil
}

// 清理因失败而过期的future
func clearFutures() {
	defer func() {
		if r := recover(); r != nil {
			// 30s 后重新运行
			time.AfterFunc(30*time.Second, func() {
				go clearFutures()
			})
		}
	}()
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for range ticker.C {
		var deletedKeys []any
		now := time.Now()
		futureHolder.Range(func(key, value any) bool {
			future := value.(*future[any])
			// 如果超时2倍以上且没有清理，则清理
			if future.exprtime.Before(now) && future.done.Load() {
				deletedKeys = append(deletedKeys, key)
			}
			return true
		})
		for _, key := range deletedKeys {
			futureHolder.Delete(key)
		}
	}
}
