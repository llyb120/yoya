package errx

import (
	"fmt"
	"runtime"
)

func Try(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, 4096)
			stackLen := runtime.Stack(stack, false)
			err = fmt.Errorf("panic: %v\nstack: %s", r, stack[:stackLen])
		}
	}()
	return fn()
}

func TryDo[T any](fn func() (T, error)) (v T, err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, 4096)
			stackLen := runtime.Stack(stack, false)
			err = fmt.Errorf("panic: %v\nstack: %s", r, stack[:stackLen])
		}
	}()
	v, err = fn()
	return
}
