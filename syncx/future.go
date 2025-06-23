package syncx

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	_ "unsafe"
)

type futureContext[T any] struct {
	data T
	wg   sync.WaitGroup
	err  error
}

type Future[T any] func() T

type FutureError func() error

type FutureAble interface {
	GetType() reflect.Type
}
type FutureCallAble interface {
	ToFunc(fn any) func(...any) (any, FutureError)
}

func (f Future[T]) GetType() reflect.Type {
	var t T
	return reflect.TypeOf(t)
}

func (f Future[T]) ToFunc(fn any) func(...any) (any, FutureError) {
	_fn := Async2[T](fn)
	return func(args ...any) (any, FutureError) {
		return _fn(args...)
	}
}

func (f Future[T]) MarshalJSON() ([]byte, error) {
	res := f()
	return encode(res)
}

func Async2[T any](fn any) func(...any) (Future[T], FutureError) {
	fv := reflect.ValueOf(fn)
	ft := fv.Type()
	return func(args ...any) (Future[T], FutureError) {
		var futureCtx = &futureContext[T]{}
		//future := &future{exprtime: time.Now().Add(5 * time.Minute)}
		var zero T
		futureCtx.data = zero
		futureCtx.wg.Add(1)
		// futureCache.Store(fn, futureCtx)

		future := func() T {
			futureCtx.wg.Wait()
			return futureCtx.data
		}
		futureError := func() error {
			futureCtx.wg.Wait()
			return futureCtx.err
		}

		// log.Println("future", (reflect.ValueOf(future).Pointer()))
		// runtime.SetFinalizer(&future, func(f any) {
		// 	log.Println("future destroyed", (reflect.ValueOf(f).Elem().Pointer()))
		// })

		if handler, ok := fn.(func() (T, error)); ok {
			go func() {
				defer func() {
					// 打印调用栈
					if r := recover(); r != nil {
						buf := make([]byte, 1024)
						n := runtime.Stack(buf, false)
						futureCtx.err = fmt.Errorf("future panic: %v\n%s", r, buf[:n])
					}
					futureCtx.wg.Done()
				}()
				result, err := handler()
				if err != nil {
					futureCtx.err = err
				} else {
					futureCtx.data = result
				}
			}()
		} else if handler, ok := fn.(func(...any) (T, error)); ok {
			go func() {
				defer func() {
					// 打印调用栈
					if r := recover(); r != nil {
						buf := make([]byte, 1024)
						n := runtime.Stack(buf, false)
						futureCtx.err = fmt.Errorf("future panic: %v\n%s", r, buf[:n])
					}
					futureCtx.wg.Done()
				}()
				result, err := handler(args...)
				if err != nil {
					futureCtx.err = err
				} else {
					futureCtx.data = result
				}
			}()
		} else {
			go func() {
				defer func() {
					// 打印调用栈
					if r := recover(); r != nil {
						buf := make([]byte, 1024)
						n := runtime.Stack(buf, false)
						futureCtx.err = fmt.Errorf("future panic: %v\n%s", r, buf[:n])
					}
					futureCtx.wg.Done()
					//futureCtx.done.Store(true)
				}()

				in := make([]reflect.Value, len(args))
				for i, arg := range args {
					in[i] = reflect.ValueOf(arg)
				}

				// 类型检查
				if len(args) != ft.NumIn() {
					futureCtx.err = fmt.Errorf("参数数量不匹配")
					return
				}
				for i := 0; i < ft.NumIn(); i++ {
					if i >= len(in) || !in[i].Type().AssignableTo(ft.In(i)) {
						futureCtx.err = fmt.Errorf("参数类型不匹配")
						return
					}
				}

				result := fv.Call(in)
				if len(result) > 0 {
					for _, r := range result {
						val := r.Interface()
						if err, ok := val.(error); ok {
							futureCtx.err = err
						} else if value, ok := val.(T); ok {
							futureCtx.data = value
						}
					}
				}
			}()
		}

		return future, futureError
	}
}

// ----------------------------- 以下为快捷方法定义 -----------------------------

func Async2_0_0[T any](fn func()) func() (Future[T], FutureError) {
	var asyncFn = Async2[T](func() (T, error) {
		fn()
		var zero T
		return zero, nil
	})
	return func() (Future[T], FutureError) {
		return asyncFn()
	}
}

func Async2_0_1[T any](fn func() T) func() (Future[T], FutureError) {
	var asyncFn = Async2[T](func() (T, error) {
		return fn(), nil
	})
	return func() (Future[T], FutureError) {
		return asyncFn()
	}
}

func Async2_0_2[R0 any, R1 error](fn func() (R0, R1)) func() (Future[R0], FutureError) {
	var asyncFn = Async2[R0](func() (R0, error) {
		return fn()
	})
	return func() (Future[R0], FutureError) {
		return asyncFn()
	}
}

func Async2_1_0[T any](fn func(T)) func(T) (Future[any], FutureError) {
	var asyncFn = Async2[any](func(args ...any) (any, error) {
		fn(args[0].(T))
		return nil, nil
	})
	return func(t T) (Future[any], FutureError) {
		return asyncFn(t)
	}
}

func Async2_1_1[P0, R0 any](fn func(P0) R0) func(P0) (Future[R0], FutureError) {
	var asyncFn = Async2[R0](func(args ...any) (R0, error) {
		return fn(args[0].(P0)), nil
	})
	return func(p0 P0) (Future[R0], FutureError) {
		return asyncFn(p0)
	}
}

func Async2_1_2[P0, P1, R0 any, R1 error](fn func(P0, P1) (R0, error)) func(P0, P1) (Future[R0], FutureError) {
	var asyncFn = Async2[R0](func(args ...any) (R0, error) {
		return fn(args[0].(P0), args[1].(P1))
	})
	return func(p0 P0, p1 P1) (Future[R0], FutureError) {
		return asyncFn(p0, p1)
	}
}

func Async2_2_0[P0, P1 any](fn func(P0, P1)) func(P0, P1) (Future[any], FutureError) {
	var asyncFn = Async2[any](func(args ...any) (any, error) {
		fn(args[0].(P0), args[1].(P1))
		return nil, nil
	})
	return func(p0 P0, p1 P1) (Future[any], FutureError) {
		return asyncFn(p0, p1)
	}
}

func Async2_2_1[P0, P1, R0 any](fn func(P0, P1) R0) func(P0, P1) (Future[R0], FutureError) {
	var asyncFn = Async2[R0](func(args ...any) (R0, error) {
		return fn(args[0].(P0), args[1].(P1)), nil
	})
	return func(p0 P0, p1 P1) (Future[R0], FutureError) {
		return asyncFn(p0, p1)
	}
}

func Async2_2_2[P0, P1, R0 any, R1 error](fn func(P0, P1) (R0, R1)) func(P0, P1) (Future[R0], FutureError) {
	var asyncFn = Async2[R0](func(args ...any) (R0, error) {
		return fn(args[0].(P0), args[1].(P1))
	})
	return func(p0 P0, p1 P1) (Future[R0], FutureError) {
		return asyncFn(p0, p1)
	}
}
func Async2_3_0[P0, P1, P2 any](fn func(P0, P1, P2)) func(P0, P1, P2) (Future[any], FutureError) {
	var asyncFn = Async2[any](func(args ...any) (any, error) {
		fn(args[0].(P0), args[1].(P1), args[2].(P2))
		return nil, nil
	})
	return func(p0 P0, p1 P1, p2 P2) (Future[any], FutureError) {
		return asyncFn(p0, p1, p2)
	}
}

func Async2_3_1[P0, P1, P2, R0 any](fn func(P0, P1, P2) R0) func(P0, P1, P2) (Future[R0], FutureError) {
	var asyncFn = Async2[R0](func(args ...any) (R0, error) {
		return fn(args[0].(P0), args[1].(P1), args[2].(P2)), nil
	})
	return func(p0 P0, p1 P1, p2 P2) (Future[R0], FutureError) {
		return asyncFn(p0, p1, p2)
	}
}

func Async2_3_2[P0, P1, P2, R0 any, R1 error](fn func(P0, P1, P2) (R0, R1)) func(P0, P1, P2) (Future[R0], FutureError) {
	var asyncFn = Async2[R0](func(args ...any) (R0, error) {
		return fn(args[0].(P0), args[1].(P1), args[2].(P2))
	})
	return func(p0 P0, p1 P1, p2 P2) (Future[R0], FutureError) {
		return asyncFn(p0, p1, p2)
	}
}

func Async2_4_0[P0, P1, P2, P3 any](fn func(P0, P1, P2, P3)) func(P0, P1, P2, P3) (Future[any], FutureError) {
	var asyncFn = Async2[any](func(args ...any) (any, error) {
		fn(args[0].(P0), args[1].(P1), args[2].(P2), args[3].(P3))
		return nil, nil
	})
	return func(p0 P0, p1 P1, p2 P2, p3 P3) (Future[any], FutureError) {
		return asyncFn(p0, p1, p2, p3)
	}
}

func Async2_4_1[P0, P1, P2, P3, R0 any](fn func(P0, P1, P2, P3) R0) func(P0, P1, P2, P3) (Future[R0], FutureError) {
	var asyncFn = Async2[R0](func(args ...any) (R0, error) {
		return fn(args[0].(P0), args[1].(P1), args[2].(P2), args[3].(P3)), nil
	})
	return func(p0 P0, p1 P1, p2 P2, p3 P3) (Future[R0], FutureError) {
		return asyncFn(p0, p1, p2, p3)
	}
}

func Async2_4_2[P0, P1, P2, P3, R0 any, R1 error](fn func(P0, P1, P2, P3) (R0, R1)) func(P0, P1, P2, P3) (Future[R0], FutureError) {
	var asyncFn = Async2[R0](func(args ...any) (R0, error) {
		return fn(args[0].(P0), args[1].(P1), args[2].(P2), args[3].(P3))
	})
	return func(p0 P0, p1 P1, p2 P2, p3 P3) (Future[R0], FutureError) {
		return asyncFn(p0, p1, p2, p3)
	}
}

func Async2_5_0[P0, P1, P2, P3, P4 any](fn func(P0, P1, P2, P3, P4)) func(P0, P1, P2, P3, P4) (Future[any], FutureError) {
	var asyncFn = Async2[any](func(args ...any) (any, error) {
		fn(args[0].(P0), args[1].(P1), args[2].(P2), args[3].(P3), args[4].(P4))
		return nil, nil
	})
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) (Future[any], FutureError) {
		return asyncFn(p0, p1, p2, p3, p4)
	}
}

func Async2_5_1[P0, P1, P2, P3, P4, R0 any](fn func(P0, P1, P2, P3, P4) R0) func(P0, P1, P2, P3, P4) (Future[R0], FutureError) {
	var asyncFn = Async2[R0](func(args ...any) (R0, error) {
		return fn(args[0].(P0), args[1].(P1), args[2].(P2), args[3].(P3), args[4].(P4)), nil
	})
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) (Future[R0], FutureError) {
		return asyncFn(p0, p1, p2, p3, p4)
	}
}

func Async2_5_2[P0, P1, P2, P3, P4, R0 any, R1 error](fn func(P0, P1, P2, P3, P4) (R0, R1)) func(P0, P1, P2, P3, P4) (Future[R0], FutureError) {
	var asyncFn = Async2[R0](func(args ...any) (R0, error) {
		return fn(args[0].(P0), args[1].(P1), args[2].(P2), args[3].(P3), args[4].(P4))
	})
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) (Future[R0], FutureError) {
		return asyncFn(p0, p1, p2, p3, p4)
	}
}

//go:linkname encode github.com/llyb120/yoya/supx.encode
func encode(v any) ([]byte, error)
