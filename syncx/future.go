package syncx

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/llyb120/gotool/errx"
)

type future[T any] struct {
	timeout time.Duration
	result  T
	err     error
	wg      sync.WaitGroup
	done    bool
}

func (f *future[T]) Get() (T, error) {
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

var futureHolder sync.Map

func Async[T any](fn any, timeout time.Duration) func(...any) *T {
	fv := reflect.ValueOf(fn)
	ft := fv.Type()

	return func(args ...any) *T {
		future := &future[any]{timeout: timeout}
		var zero T
		ptrResult := &zero
		future.wg.Add(1)
		futureHolder.Store(ptrResult, future)

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

func HasError(objs ...any) error {
	errs := hasError(objs...)
	if errs != nil && errs.HasError() {
		return errs
	}
	return nil
}

func IsFailed(objs ...any) bool {
	errs := hasError(objs...)
	if errs != nil && !errs.HasError() {
		return true
	}
	return false
}

func hasError(objs ...any) *errx.MultiError {
	var errs = &errx.MultiError{}
	if len(objs) > 1 {
		var g Group
		for _, obj := range objs {
			obj := obj
			g.Go(func() error {
				f, ok := futureHolder.Load(obj)
				if !ok {
					return nil
				}
				_, err := f.(*future[any]).Get()
				if err != nil {
					errs.Add(err)
				}
				return nil
			})
		}
		if err := g.Wait(); err != nil {
			errs.Add(err)
		}
	} else if len(objs) == 1 {
		f, ok := futureHolder.Load(objs[0])
		if !ok {
			return nil
		}
		_, err := f.(*future[any]).Get()
		if err != nil {
			errs.Add(err)
		}
	}
	if errs.HasError() {
		return errs
	}
	return nil
}
