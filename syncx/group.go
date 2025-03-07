package syncx

import (
	"fmt"
	"runtime"
	"sync"
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

func (g *Group) Wait() error {
	g.wg.Wait()
	if g.eg != nil && len(g.eg.errors) > 0 {
		return g.eg
	}
	return nil
}
