package syncx

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type future[T any] struct {
	timeout  time.Duration
	exprtime time.Time
	result   T
	err      error
	wg       sync.WaitGroup
	done     atomic.Bool
}

func (f *future[T]) Get() (T, error) {
	if f.done.Load() {
		return f.result, f.err
	}

	// 创建一个带超时的channel
	done := make(chan struct{})
	go func() {
		f.wg.Wait()
		close(done)
	}()

	// 等待结果或超时
	if f.timeout > 0 {
		select {
		case <-done:
			if f.err != nil {
				return f.result, f.err
			}
		case <-time.After(f.timeout):
			var zero T
			return zero, fmt.Errorf("future get timeout after %v", f.timeout)
		}
	} else {
		// 没有超时，等待结果
		<-done
		if f.err != nil {
			return f.result, f.err
		}
	}

	return f.result, nil
}

var futureHolder sync.Map

func Async[T any](fn any, timeout ...time.Duration) func(...any) *T {
	fv := reflect.ValueOf(fn)
	ft := fv.Type()
	var t time.Duration
	if len(timeout) > 0 {
		t = timeout[0]
	}
	if t <= 0 {
		t = 0
	}
	return func(args ...any) *T {
		future := &future[any]{timeout: t, exprtime: time.Now().Add(2 * t)}
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
	return hasError(objs...)
}

func hasError(objs ...any) error {
	if len(objs) > 1 {
		var g Group
		for _, obj := range objs {
			obj := obj
			g.Go(func() error {
				f, ok := futureHolder.LoadAndDelete(obj)
				if !ok {
					return nil
				}
				_, err := f.(*future[any]).Get()
				return err
			})
		}
		if err := g.Wait(); err != nil {
			return err
		}
	} else if len(objs) == 1 {
		f, ok := futureHolder.LoadAndDelete(objs[0])
		if !ok {
			return nil
		}
		_, err := f.(*future[any]).Get()
		return err
	}
	return nil
}

// 清理因失败而过期的future
func clearFutures() {
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
