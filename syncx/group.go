package syncx

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"github.com/llyb120/yoya/errx"
	"github.com/llyb120/yoya/stlx"
	"github.com/petermattis/goid"
)

type Group struct {
	wg sync.WaitGroup
	eg errx.MultiError
}

var globalGroupHolder = stlx.NewSyncBimMap[int64, int64]()

func (g *Group) Go(fn func() error) {
	g.wg.Add(1)
	var parentGoid = goid.Get()
	go func() {
		defer g.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				stack := make([]byte, 4096)
				stackLen := runtime.Stack(stack, false)
				g.eg.Add(fmt.Errorf("panic: %v\nstack: %s", r, stack[:stackLen]))
			}
		}()
		// 存储协程id
		localGoid := goid.Get()
		globalGroupHolder.Set(localGoid, parentGoid)
		defer globalGroupHolder.Del(localGoid)
		// 调用
		err := fn()
		if err != nil {
			g.eg.Add(err)
		}
	}()
}

func (g *Group) Wait(timeout ...time.Duration) error {
	if len(timeout) > 0 {
		return g.waitWithTimeout(timeout[0])
	}
	return g.waitWithTimeout(0)
}

func (g *Group) waitWithTimeout(timeout time.Duration) error {
	if timeout <= 0 {
		// 无超时，直接等待
		g.wg.Wait()
		if g.eg.HasError() {
			return &g.eg
		}
		return nil
	}

	// 创建一个带超时的channel
	done := make(chan struct{})
	go func() {
		g.wg.Wait()
		close(done)
	}()

	// 等待结果或超时
	select {
	case <-done:
		if g.eg.HasError() {
			return &g.eg
		}
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("等待超时，超过 %v", timeout)
	}
}
