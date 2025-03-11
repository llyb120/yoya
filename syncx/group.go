package syncx

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Group struct {
	wg   sync.WaitGroup
	eg   *groupError
	once sync.Once
}

type groupError struct {
	mu     sync.Mutex
	errors []error
}

func (g *groupError) Error() string {
	return fmt.Sprintf("%v", g.errors)
}

func (g *groupError) Append(err error) {
	g.mu.Lock()
	g.errors = append(g.errors, err)
	g.mu.Unlock()
}

func (g *Group) Go(fn func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		defer func() {
			if r := recover(); r != nil {
				stack := make([]byte, 4096)
				stackLen := runtime.Stack(stack, false)
				g.append(fmt.Errorf("panic: %v\nstack: %s", r, stack[:stackLen]))
			}
		}()
		err := fn()
		if err != nil {
			g.append(err)
		}
	}()
}

func (g *Group) append(err error) {
	g.once.Do(func() {
		g.eg = &groupError{}
	})
	g.eg.Append(err)
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
		if g.eg != nil && len(g.eg.errors) > 0 {
			return g.eg
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
		if g.eg != nil && len(g.eg.errors) > 0 {
			return g.eg
		}
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("等待超时，超过 %v", timeout)
	}
}
