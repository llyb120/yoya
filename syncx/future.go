package syncx

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	_ "unsafe"
)

type Future[T any] func() T

type FutureAble interface {
	GetType() reflect.Type
}
type FutureCallAble interface {
	ToFunc(fn func(args ...any) any) func(...any) (any, Future[error])
}

func (f Future[T]) GetType() reflect.Type {
	var t T
	return reflect.TypeOf(t)
}

func (f Future[T]) ToFunc(fn func(args ...any) T) func(...any) (any, Future[error]) {
	return func(args ...any) (any, Future[error]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r T
		var err error
		go func() {
			defer wg.Done()
			defer handlePanic(&err)
			r = fn(args...)
		}()
		wg.Wait()
		return func() T {
				return r
			}, func() error {
				return err
			}
	}
}

func (f Future[T]) MarshalJSON() ([]byte, error) {
	res := f()
	return json.Marshal(res)
}

func Mirai[T any]() Future[T] {
	return func() T {
		var zero T
		return zero
	}
}

// ----------------------------- 以下为快捷方法定义 -----------------------------

func handlePanic(args ...any) {
	if r := recover(); r != nil {
		for _, arg := range args {
			if errPtr, ok := arg.(*error); ok && errPtr != nil {
				buf := make([]byte, 1024)
				n := runtime.Stack(buf, false)
				// 复制 error 指针的值
				*errPtr = fmt.Errorf("future panic: %v\n%s", r, buf[:n])
				break
			}
		}
	}
}

func Async2_0_0(fn func()) func() Future[any] {
	return func() Future[any] {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer handlePanic()
			fn()
		}()
		return func() any {
			wg.Wait()
			return nil
		}
	}
}

func Async2_0_1[T any](fn func() T) func() Future[T] {
	return func() Future[T] {
		var wg sync.WaitGroup
		wg.Add(1)
		var r T
		go func() {
			defer wg.Done()
			defer handlePanic(&r)
			r = fn()
		}()
		return func() T {
			wg.Wait()
			return r
		}
	}
}

func Async2_0_2[R0 any, R1 any](fn func() (R0, R1)) func() (Future[R0], Future[R1]) {
	return func() (Future[R0], Future[R1]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1)
			r0, r1 = fn()
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}
	}
}

func Async2_0_3[R0 any, R1 any, R2 any](fn func() (R0, R1, R2)) func() (Future[R0], Future[R1], Future[R2]) {
	return func() (Future[R0], Future[R1], Future[R2]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		var r2 R2
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1, &r2)
			r0, r1, r2 = fn()
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}, func() R2 {
				wg.Wait()
				return r2
			}
	}
}

func Async2_0_4[R0 any, R1 any, R2 any, R3 any](fn func() (R0, R1, R2, R3)) func() (Future[R0], Future[R1], Future[R2], Future[R3]) {
	return func() (Future[R0], Future[R1], Future[R2], Future[R3]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		var r2 R2
		var r3 R3
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1, &r2, &r3)
			r0, r1, r2, r3 = fn()
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}, func() R2 {
				wg.Wait()
				return r2
			}, func() R3 {
				wg.Wait()
				return r3
			}
	}
}

func Async2_1_0[T any](fn func(T)) func(T) Future[any] {
	return func(t T) Future[any] {
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			defer handlePanic()
			fn(t)
		}()
		return func() any {
			wg.Wait()
			return nil
		}
	}
}

func Async2_1_1[P0, R0 any](fn func(P0) R0) func(P0) Future[R0] {
	return func(p0 P0) Future[R0] {
		var wg sync.WaitGroup
		wg.Add(1)
		var r R0
		go func() {
			defer wg.Done()
			defer handlePanic(&r)
			r = fn(p0)
		}()
		return func() R0 {
			wg.Wait()
			return r
		}
	}
}

func Async2_1_2[P0, P1, R0 any, R1 any](fn func(P0, P1) (R0, R1)) func(P0, P1) (Future[R0], Future[R1]) {
	return func(p0 P0, p1 P1) (Future[R0], Future[R1]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1)
			r0, r1 = fn(p0, p1)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}
	}
}

func Async2_1_3[P0, P1, P2, R0 any, R1 any, R2 any](fn func(P0, P1, P2) (R0, R1, R2)) func(P0, P1, P2) (Future[R0], Future[R1], Future[R2]) {
	return func(p0 P0, p1 P1, p2 P2) (Future[R0], Future[R1], Future[R2]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		var r2 R2
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1, &r2)
			r0, r1, r2 = fn(p0, p1, p2)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}, func() R2 {
				wg.Wait()
				return r2
			}
	}
}

func Async2_1_4[P0, R0 any, R1 any, R2 any, R3 any](fn func(P0) (R0, R1, R2, R3)) func(P0) (Future[R0], Future[R1], Future[R2], Future[R3]) {
	return func(p0 P0) (Future[R0], Future[R1], Future[R2], Future[R3]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		var r2 R2
		var r3 R3
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1, &r2, &r3)
			r0, r1, r2, r3 = fn(p0)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}, func() R2 {
				wg.Wait()
				return r2
			}, func() R3 {
				wg.Wait()
				return r3
			}
	}
}

func Async2_2_0[P0, P1 any](fn func(P0, P1)) func(P0, P1) Future[any] {
	return func(p0 P0, p1 P1) Future[any] {
		var wg sync.WaitGroup
		wg.Add(1)
		var r any
		go func() {
			defer wg.Done()
			defer handlePanic(&r)
			fn(p0, p1)
		}()
		return func() any {
			wg.Wait()
			return r
		}
	}
}

func Async2_2_1[P0, P1, R0 any](fn func(P0, P1) R0) func(P0, P1) Future[R0] {
	return func(p0 P0, p1 P1) Future[R0] {
		var wg sync.WaitGroup
		wg.Add(1)
		var r R0
		go func() {
			defer wg.Done()
			defer handlePanic(&r)
			r = fn(p0, p1)
		}()
		return func() R0 {
			wg.Wait()
			return r
		}
	}
}

func Async2_2_2[P0, P1, R0 any, R1 any](fn func(P0, P1) (R0, R1)) func(P0, P1) (Future[R0], Future[R1]) {
	return func(p0 P0, p1 P1) (Future[R0], Future[R1]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1)
			r0, r1 = fn(p0, p1)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}
	}
}

func Async2_2_3[P0, P1, P2, R0 any, R1 any, R2 any](fn func(P0, P1, P2) (R0, R1, R2)) func(P0, P1, P2) (Future[R0], Future[R1], Future[R2]) {
	return func(p0 P0, p1 P1, p2 P2) (Future[R0], Future[R1], Future[R2]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		var r2 R2
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1, &r2)
			r0, r1, r2 = fn(p0, p1, p2)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}, func() R2 {
				wg.Wait()
				return r2
			}
	}
}

func Async2_2_4[P0, P1, R0 any, R1 any, R2 any, R3 any](fn func(P0, P1) (R0, R1, R2, R3)) func(P0, P1) (Future[R0], Future[R1], Future[R2], Future[R3]) {
	return func(p0 P0, p1 P1) (Future[R0], Future[R1], Future[R2], Future[R3]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		var r2 R2
		var r3 R3
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1, &r2, &r3)
			r0, r1, r2, r3 = fn(p0, p1)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}, func() R2 {
				wg.Wait()
				return r2
			}, func() R3 {
				wg.Wait()
				return r3
			}
	}
}

func Async2_3_0[P0, P1, P2 any](fn func(P0, P1, P2)) func(P0, P1, P2) Future[any] {
	return func(p0 P0, p1 P1, p2 P2) Future[any] {
		var wg sync.WaitGroup
		wg.Add(1)
		var r any
		go func() {
			defer wg.Done()
			defer handlePanic(&r)
			fn(p0, p1, p2)
		}()
		return func() any {
			wg.Wait()
			return r
		}
	}
}

func Async2_3_1[P0, P1, P2, R0 any](fn func(P0, P1, P2) R0) func(P0, P1, P2) Future[R0] {
	return func(p0 P0, p1 P1, p2 P2) Future[R0] {
		var wg sync.WaitGroup
		wg.Add(1)
		var r R0
		go func() {
			defer wg.Done()
			defer handlePanic(&r)
			r = fn(p0, p1, p2)
		}()
		return func() R0 {
			wg.Wait()
			return r
		}
	}
}

func Async2_3_2[P0, P1, P2, R0 any, R1 any](fn func(P0, P1, P2) (R0, R1)) func(P0, P1, P2) (Future[R0], Future[R1]) {
	return func(p0 P0, p1 P1, p2 P2) (Future[R0], Future[R1]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1)
			r0, r1 = fn(p0, p1, p2)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}
	}
}

func Async2_3_3[P0, P1, P2, R0 any, R1 any, R2 any](fn func(P0, P1, P2) (R0, R1, R2)) func(P0, P1, P2) (Future[R0], Future[R1], Future[R2]) {
	return func(p0 P0, p1 P1, p2 P2) (Future[R0], Future[R1], Future[R2]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		var r2 R2
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1, &r2)
			r0, r1, r2 = fn(p0, p1, p2)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}, func() R2 {
				wg.Wait()
				return r2
			}
	}
}

func Async2_3_4[P0, P1, P2, R0 any, R1 any, R2 any, R3 any](fn func(P0, P1, P2) (R0, R1, R2, R3)) func(P0, P1, P2) (Future[R0], Future[R1], Future[R2], Future[R3]) {
	return func(p0 P0, p1 P1, p2 P2) (Future[R0], Future[R1], Future[R2], Future[R3]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		var r2 R2
		var r3 R3
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1, &r2, &r3)
			r0, r1, r2, r3 = fn(p0, p1, p2)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}, func() R2 {
				wg.Wait()
				return r2
			}, func() R3 {
				wg.Wait()
				return r3
			}
	}
}

func Async2_4_0[P0, P1, P2, P3 any](fn func(P0, P1, P2, P3)) func(P0, P1, P2, P3) Future[any] {
	return func(p0 P0, p1 P1, p2 P2, p3 P3) Future[any] {
		var wg sync.WaitGroup
		wg.Add(1)
		var r any
		go func() {
			defer wg.Done()
			defer handlePanic(&r)
			fn(p0, p1, p2, p3)
		}()
		return func() any {
			wg.Wait()
			return r
		}
	}
}

func Async2_4_1[P0, P1, P2, P3, R0 any](fn func(P0, P1, P2, P3) R0) func(P0, P1, P2, P3) Future[R0] {
	return func(p0 P0, p1 P1, p2 P2, p3 P3) Future[R0] {
		var wg sync.WaitGroup
		wg.Add(1)
		var r R0
		go func() {
			defer wg.Done()
			defer handlePanic(&r)
			r = fn(p0, p1, p2, p3)
		}()

		return func() R0 {
			wg.Wait()
			return r
		}
	}
}

func Async2_4_2[P0, P1, P2, P3, R0 any, R1 any](fn func(P0, P1, P2, P3) (R0, R1)) func(P0, P1, P2, P3) (Future[R0], Future[R1]) {
	return func(p0 P0, p1 P1, p2 P2, p3 P3) (Future[R0], Future[R1]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1)
			r0, r1 = fn(p0, p1, p2, p3)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}
	}
}

func Async2_4_3[P0, P1, P2, P3, R0 any, R1 any, R2 any](fn func(P0, P1, P2, P3) (R0, R1, R2)) func(P0, P1, P2, P3) (Future[R0], Future[R1], Future[R2]) {
	return func(p0 P0, p1 P1, p2 P2, p3 P3) (Future[R0], Future[R1], Future[R2]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		var r2 R2
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1, &r2)
			r0, r1, r2 = fn(p0, p1, p2, p3)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}, func() R2 {
				wg.Wait()
				return r2
			}
	}
}

func Async2_4_4[P0, P1, P2, P3, R0 any, R1 any, R2 any, R3 any](fn func(P0, P1, P2, P3) (R0, R1, R2, R3)) func(P0, P1, P2, P3) (Future[R0], Future[R1], Future[R2], Future[R3]) {
	return func(p0 P0, p1 P1, p2 P2, p3 P3) (Future[R0], Future[R1], Future[R2], Future[R3]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		var r2 R2
		var r3 R3
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1, &r2, &r3)
			r0, r1, r2, r3 = fn(p0, p1, p2, p3)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}, func() R2 {
				wg.Wait()
				return r2
			}, func() R3 {
				wg.Wait()
				return r3
			}
	}
}

func Async2_5_0[P0, P1, P2, P3, P4 any](fn func(P0, P1, P2, P3, P4)) func(P0, P1, P2, P3, P4) Future[any] {
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) Future[any] {
		var wg sync.WaitGroup
		wg.Add(1)
		var r any
		go func() {
			defer wg.Done()
			defer handlePanic(&r)
			fn(p0, p1, p2, p3, p4)
		}()

		return func() any {
			wg.Wait()
			return r
		}
	}
}

func Async2_5_1[P0, P1, P2, P3, P4, R0 any](fn func(P0, P1, P2, P3, P4) R0) func(P0, P1, P2, P3, P4) Future[R0] {
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) Future[R0] {
		var wg sync.WaitGroup
		wg.Add(1)
		var r R0
		go func() {
			defer wg.Done()
			defer handlePanic(&r)
			r = fn(p0, p1, p2, p3, p4)
		}()

		return func() R0 {
			wg.Wait()
			return r
		}
	}
}

func Async2_5_2[P0, P1, P2, P3, P4, R0 any, R1 any](fn func(P0, P1, P2, P3, P4) (R0, R1)) func(P0, P1, P2, P3, P4) (Future[R0], Future[R1]) {
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) (Future[R0], Future[R1]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1)
			r0, r1 = fn(p0, p1, p2, p3, p4)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}
	}
}

func Async2_5_3[P0, P1, P2, P3, P4, R0 any, R1 any, R2 any](fn func(P0, P1, P2, P3, P4) (R0, R1, R2)) func(P0, P1, P2, P3, P4) (Future[R0], Future[R1], Future[R2]) {
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) (Future[R0], Future[R1], Future[R2]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		var r2 R2
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1, &r2)
			r0, r1, r2 = fn(p0, p1, p2, p3, p4)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}, func() R2 {
				wg.Wait()
				return r2
			}
	}
}

func Async2_5_4[P0, P1, P2, P3, P4, R0 any, R1 any, R2 any, R3 any](fn func(P0, P1, P2, P3, P4) (R0, R1, R2, R3)) func(P0, P1, P2, P3, P4) (Future[R0], Future[R1], Future[R2], Future[R3]) {
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) (Future[R0], Future[R1], Future[R2], Future[R3]) {
		var wg sync.WaitGroup
		wg.Add(1)
		var r0 R0
		var r1 R1
		var r2 R2
		var r3 R3
		go func() {
			defer wg.Done()
			defer handlePanic(&r0, &r1, &r2, &r3)
			r0, r1, r2, r3 = fn(p0, p1, p2, p3, p4)
		}()
		return func() R0 {
				wg.Wait()
				return r0
			}, func() R1 {
				wg.Wait()
				return r1
			}, func() R2 {
				wg.Wait()
				return r2
			}, func() R3 {
				wg.Wait()
				return r3
			}
	}
}
