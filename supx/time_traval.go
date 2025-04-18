package supx

import (
	"log"
	"runtime"
	"sync"
)

type TimeLeapOption int

const (
	Async TimeLeapOption = iota
)

type TimeLeapAble interface {
	Mark(func())
}

type timeLeap struct {
	mu    sync.Mutex
	funcs []func()
	async bool
}

func (t *timeLeap) Mark(f func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.funcs = append(t.funcs, f)
}

func TimeLeap(opts ...TimeLeapOption) (TimeLeapAble, func()) {
	t := &timeLeap{
		funcs: make([]func(), 0, 8),
	}

	for _, opt := range opts {
		switch opt {
		case Async:
			t.async = true
		}
	}
	return t, func() {
		t.mu.Lock()
		defer t.mu.Unlock()
		if t.async {
			var wg sync.WaitGroup
			for _, f := range t.funcs {
				wg.Add(1)
				go func(f func()) {
					defer func() {
						if r := recover(); r != nil {
							// 打印调用栈
							buf := make([]byte, 1024)
							buf = buf[:runtime.Stack(buf, false)]
							log.Printf("time leap panic: %v\n%s", r, buf)
						}
					}()
					defer wg.Done()
					f()
				}(f)
			}
			wg.Wait()
		} else {
			for _, f := range t.funcs {
				f()
			}
		}
	}
}
