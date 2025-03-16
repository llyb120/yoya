package stlx

import (
	"sync"
)

type lock struct {
	sync bool
	mu   sync.RWMutex
}

func (l *lock) Lock() {
	if l.sync {
		l.mu.Lock()
	}
}

func (l *lock) Unlock() {
	if l.sync {
		l.mu.Unlock()
	}
}

func (l *lock) RLock() {
	if l.sync {
		l.mu.RLock()
	}
}

func (l *lock) RUnlock() {
	if l.sync {
		l.mu.RUnlock()
	}
}
