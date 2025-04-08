package syncx

import (
	"sync"
)

// 增强的池，可以设置终结函数，用于对象的回收
type pool[T any] struct {
	p sync.Pool
	PoolOption[T]
	size     int
	prepared chan T
}

type Pool[T any] interface {
	Get() (T, func())
}

type PoolOption[T any] struct {
	Finalizer func(*T)
	New       func() T
}

func (opt PoolOption[T]) Build() *pool[T] {
	return &pool[T]{
		p: sync.Pool{New: func() any {
			return opt.New()
		}},
		PoolOption: opt,
	}
}

// func NewPool[T any](opts PoolOption[T]) Pool[T] {
// 	var ptr = &pool[T]{
// 		Pool: sync.Pool{New: func() any {
// 			return opts.New()
// 		}},
// 		PoolOption: opts,
// 	}
// 	return ptr
// }

func (p *pool[T]) Get() (T, func()) {
	item := p.p.Get().(T)
	ptr := &item
	return item, func() {
		p.put(ptr)
	}
}

func (p *pool[T]) put(v *T) {
	if p.Finalizer != nil {
		p.Finalizer(v)
	}
	p.p.Put(*v)
}
