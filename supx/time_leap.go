package supx

import (
	"log"
	"runtime"
	"sync"

	"github.com/llyb120/yoya/syncx"
)

type TimeLeapOption int

const (
	Async TimeLeapOption = iota
)

type TimeLeapAble interface {
	Leap(func())
}

type timeLeap struct {
	mu    sync.Mutex
	funcs []func()
	waits []any
	async bool
}

func (t *timeLeap) Leap(f func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.funcs = append(t.funcs, f)
}

func TimeLeap(opts ...any) (TimeLeapAble, func()) {
	t := &timeLeap{
		funcs: make([]func(), 0, 8),
		waits: make([]any, 0, 2),
	}

	for _, opt := range opts {
		switch tp := opt.(type) {
		case TimeLeapOption:
			switch tp {
			case Async:
				t.async = true
			}
		default:
			t.waits = append(t.waits, opt)
		}
	}
	return t, func() {
		t.mu.Lock()
		defer t.mu.Unlock()
		if len(t.waits) > 0 {
			if err := syncx.Await(t.waits...); err != nil {
				log.Printf("time leap await error: %v", err)
				return
			}
		}
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
