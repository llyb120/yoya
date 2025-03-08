package syncx

import (
	"fmt"
	"reflect"
	"sync"
	"time"
)

type Future[T any] struct {
	timeout time.Duration
	result  T
	err     error
	wg      sync.WaitGroup
	done    bool
}

func (f *Future[T]) Get() (T, error) {
	if f.done {
		return f.result, f.err
	}

	// 创建一个带超时的channel
	done := make(chan struct{})
	go func() {
		f.wg.Wait()
		close(done)
	}()

	// 等待结果或超时
	select {
	case <-done:
		if f.err != nil {
			return f.result, f.err
		}
	case <-time.After(f.timeout):
		var zero T
		return zero, fmt.Errorf("future get timeout after %v", f.timeout)
	}

	return f.result, nil
}

func Async[T any](fn any, timeout time.Duration) func(...any) *Future[T] {
	fv := reflect.ValueOf(fn)
	ft := fv.Type()

	return func(args ...any) *Future[T] {
		future := &Future[T]{timeout: timeout}
		future.wg.Add(1)

		go func() {
			defer func() {
				future.wg.Done()
				future.done = true
				if r := recover(); r != nil {
					future.err = fmt.Errorf("future panic: %v", r)
				}
			}()

			in := make([]reflect.Value, len(args))
			for i, arg := range args {
				in[i] = reflect.ValueOf(arg)
			}

			// 类型检查
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
					} else {
						future.result = val.(T)
					}
				}
			}
		}()

		return future
	}
}
