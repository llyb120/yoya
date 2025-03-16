package syncx

import (
	"fmt"
	"time"
)

func async[T any](handler func() (T, error)) *T {
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

		result, err := handler()
		if err != nil {
			future.err = err
		} else {
			*ptrResult = result
			future.result = result
		}
		// if len(result) > 0 {
		// 	for _, r := range result {
		// 		val := r.Interface()
		// 		if err, ok := val.(error); ok {
		// 			future.err = err
		// 		} else if value, ok := val.(T); ok {
		// 			*ptrResult = value
		// 			future.result = value
		// 		}
		// 	}
		// }
	}()

	return ptrResult
}

func Async_0_0(fn func()) func() *any {
	return func() *any {
		var asyncFn = async[any](func() (any, error) {
			fn()
			return nil, nil
		})
		return asyncFn
	}
}

func Async_0_1[T any](fn func() T) func() *T {
	return func() *T {
		var asyncFn = async[T](func() (T, error) {
			return fn(), nil
		})
		return asyncFn
	}
}

func Async_0_2[R0 any, R1 error](fn func() (R0, R1)) func() *R0 {
	return func() *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn()
		})
		return asyncFn
	}
}

func Async_1_0[T any](fn func(T)) func(T) *any {
	return func(t T) *any {
		var asyncFn = async[any](func() (any, error) {
			fn(t)
			return nil, nil
		})
		return asyncFn
	}
}

func Async_1_1[P0, R0 any](fn func(P0) R0) func(P0) *R0 {
	return func(p0 P0) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0), nil
		})
		return asyncFn
	}
}

func Async_1_2[P0, P1, R0 any, R1 error](fn func(P0, P1) (R0, error)) func(P0, P1) *R0 {
	return func(p0 P0, p1 P1) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1)
		})
		return asyncFn
	}
}

func Async_2_0[R0, R1 any](fn func(R0, R1)) func(R0, R1) *any {
	return func(r0 R0, r1 R1) *any {
		var asyncFn = async[any](func() (any, error) {
			fn(r0, r1)
			return nil, nil
		})
		return asyncFn
	}
}

func Async_2_1[P0, P1, R0 any](fn func(P0, P1) R0) func(P0, P1) *R0 {
	return func(p0 P0, p1 P1) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1), nil
		})
		return asyncFn
	}
}

func Async_2_2[P0, P1, R0 any, R1 error](fn func(P0, P1) (R0, R1)) func(P0, P1) *R0 {
	return func(p0 P0, p1 P1) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1)
		})
		return asyncFn
	}
}
func Async_3_0[P0, P1, P2 any](fn func(P0, P1, P2)) func(P0, P1, P2) *any {
	return func(p0 P0, p1 P1, p2 P2) *any {
		var asyncFn = async[any](func() (any, error) {
			fn(p0, p1, p2)
			return nil, nil
		})
		return asyncFn
	}
}

func Async_3_1[P0, P1, P2, R0 any](fn func(P0, P1, P2) R0) func(P0, P1, P2) *R0 {
	return func(p0 P0, p1 P1, p2 P2) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1, p2), nil
		})
		return asyncFn
	}
}

func Async_3_2[P0, P1, P2, R0 any, R1 error](fn func(P0, P1, P2) (R0, R1)) func(P0, P1, P2) *R0 {
	return func(p0 P0, p1 P1, p2 P2) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1, p2)
		})
		return asyncFn
	}
}

func Async_4_0[P0, P1, P2, P3 any](fn func(P0, P1, P2, P3)) func(P0, P1, P2, P3) *any {
	return func(p0 P0, p1 P1, p2 P2, p3 P3) *any {
		var asyncFn = async[any](func() (any, error) {
			fn(p0, p1, p2, p3)
			return nil, nil
		})
		return asyncFn
	}
}

func Async_4_1[P0, P1, P2, P3, R0 any](fn func(P0, P1, P2, P3) R0) func(P0, P1, P2, P3) *R0 {
	return func(p0 P0, p1 P1, p2 P2, p3 P3) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1, p2, p3), nil
		})
		return asyncFn
	}
}

func Async_4_2[P0, P1, P2, P3, R0 any, R1 error](fn func(P0, P1, P2, P3) (R0, R1)) func(P0, P1, P2, P3) *R0 {
	return func(p0 P0, p1 P1, p2 P2, p3 P3) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1, p2, p3)
		})
		return asyncFn
	}
}

func Async_5_0[P0, P1, P2, P3, P4 any](fn func(P0, P1, P2, P3, P4)) func(P0, P1, P2, P3, P4) *any {
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) *any {
		var asyncFn = async[any](func() (any, error) {
			fn(p0, p1, p2, p3, p4)
			return nil, nil
		})
		return asyncFn
	}
}

func Async_5_1[P0, P1, P2, P3, P4, R0 any](fn func(P0, P1, P2, P3, P4) R0) func(P0, P1, P2, P3, P4) *R0 {
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1, p2, p3, p4), nil
		})
		return asyncFn
	}
}

func Async_5_2[P0, P1, P2, P3, P4, R0 any, R1 error](fn func(P0, P1, P2, P3, P4) (R0, R1)) func(P0, P1, P2, P3, P4) *R0 {
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1, p2, p3, p4)
		})
		return asyncFn
	}
}
