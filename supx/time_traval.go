package supx

import (
	"sync"
)

type TimeTravelAble interface {
	Mark(func())
}

type timeTravel struct {
	mu    sync.Mutex
	funcs []func()
}

func (t *timeTravel) Mark(f func()) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.funcs = append(t.funcs, f)
}

func TimeTravel() (TimeTravelAble, func()) {
	t := &timeTravel{
		funcs: make([]func(), 0, 8),
	}
	return t, func() {
		t.mu.Lock()
		defer t.mu.Unlock()
		for _, f := range t.funcs {
			f()
		}
	}
}
