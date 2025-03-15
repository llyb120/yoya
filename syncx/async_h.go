package syncx

func Async_0_0(fn func()) func() *any {
	var asyncFn = Async[any](fn)
	return func() *any {
		return asyncFn()
	}
}

func Async_0_1[T any](fn func() T) func() *T {
	var asyncFn = Async[T](fn)
	return func() *T {
		return asyncFn()
	}
}

func Async_0_2[R0, R1 any](fn func() (R0, R1)) func() *R0 {
	var asyncFn = Async[R0](fn)
	return func() *R0 {
		return asyncFn()
	}
}

func Async_1_0[T any](fn func(T)) func(T) *any {
	var asyncFn = Async[any](fn)
	return func(t T) *any {
		return asyncFn(t)
	}
}

func Async_1_1[P0, R0 any](fn func(P0) R0) func(P0) *R0 {
	var asyncFn = Async[R0](fn)
	return func(p0 P0) *R0 {
		return asyncFn(p0)
	}
}

func Async_1_2[P0, P1, R0 any, R1 error](fn func(P0, P1) (R0, error)) func(P0, P1) *R0 {
	var asyncFn = Async[R0](fn)
	return func(p0 P0, p1 P1) *R0 {
		return asyncFn(p0, p1)
	}
}

func Async_2_0[R0, R1 any](fn func(R0, R1)) func(R0, R1) *any {
	var asyncFn = Async[any](fn)
	return func(r0 R0, r1 R1) *any {
		return asyncFn(r0, r1)
	}
}

func Async_2_1[P0, P1, R0 any](fn func(P0, P1) R0) func(P0, P1) *R0 {
	var asyncFn = Async[R0](fn)
	return func(p0 P0, p1 P1) *R0 {
		return asyncFn(p0, p1)
	}
}

func Async_2_2[P0, P1, R0 any, R1 error](fn func(P0, P1) (R0, R1)) func(P0, P1) *R0 {
	var asyncFn = Async[R0](fn)
	return func(p0 P0, p1 P1) *R0 {
		return asyncFn(p0, p1)
	}
}
func Async_3_0[P0, P1, P2 any](fn func(P0, P1, P2)) func(P0, P1, P2) *any {
	var asyncFn = Async[any](fn)
	return func(p0 P0, p1 P1, p2 P2) *any {
		return asyncFn(p0, p1, p2)
	}
}

func Async_3_1[P0, P1, P2, R0 any](fn func(P0, P1, P2) R0) func(P0, P1, P2) *R0 {
	var asyncFn = Async[R0](fn)
	return func(p0 P0, p1 P1, p2 P2) *R0 {
		return asyncFn(p0, p1, p2)
	}
}

func Async_3_2[P0, P1, P2, R0 any, R1 error](fn func(P0, P1, P2) (R0, R1)) func(P0, P1, P2) *R0 {
	var asyncFn = Async[R0](fn)
	return func(p0 P0, p1 P1, p2 P2) *R0 {
		return asyncFn(p0, p1, p2)
	}
}

func Async_4_0[P0, P1, P2, P3 any](fn func(P0, P1, P2, P3)) func(P0, P1, P2, P3) *any {
	var asyncFn = Async[any](fn)
	return func(p0 P0, p1 P1, p2 P2, p3 P3) *any {
		return asyncFn(p0, p1, p2, p3)
	}
}

func Async_4_1[P0, P1, P2, P3, R0 any](fn func(P0, P1, P2, P3) R0) func(P0, P1, P2, P3) *R0 {
	var asyncFn = Async[R0](fn)
	return func(p0 P0, p1 P1, p2 P2, p3 P3) *R0 {
		return asyncFn(p0, p1, p2, p3)
	}
}

func Async_4_2[P0, P1, P2, P3, R0 any, R1 error](fn func(P0, P1, P2, P3) (R0, R1)) func(P0, P1, P2, P3) *R0 {
	var asyncFn = Async[R0](fn)
	return func(p0 P0, p1 P1, p2 P2, p3 P3) *R0 {
		return asyncFn(p0, p1, p2, p3)
	}
}

func Async_5_0[P0, P1, P2, P3, P4 any](fn func(P0, P1, P2, P3, P4)) func(P0, P1, P2, P3, P4) *any {
	var asyncFn = Async[any](fn)
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) *any {
		return asyncFn(p0, p1, p2, p3, p4)
	}
}

func Async_5_1[P0, P1, P2, P3, P4, R0 any](fn func(P0, P1, P2, P3, P4) R0) func(P0, P1, P2, P3, P4) *R0 {
	var asyncFn = Async[R0](fn)
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) *R0 {
		return asyncFn(p0, p1, p2, p3, p4)
	}
}

func Async_5_2[P0, P1, P2, P3, P4, R0 any, R1 error](fn func(P0, P1, P2, P3, P4) (R0, R1)) func(P0, P1, P2, P3, P4) *R0 {
	var asyncFn = Async[R0](fn)
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) *R0 {
		return asyncFn(p0, p1, p2, p3, p4)
	}
}
