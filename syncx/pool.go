package syncx

import (
	"sync"
)

// 增强的池，可以设置终结函数，用于对象的回收
type pool[T any] struct {
	sync.Pool
	PoolOption[T]
	size     int
	prepared chan T
}

type Pool[T any] interface {
	Get() T
	Put(v T)
}

type PoolOption[T any] struct {
	Finalizer func(*T)
	New       func() T
}

func NewPool[T any](opts PoolOption[T]) Pool[T] {
	var ptr = &pool[T]{
		Pool: sync.Pool{New: func() any {
			return opts.New()
		}},
		PoolOption: opts,
	}
	return ptr
}

func (p *pool[T]) Get() T {
	return p.Pool.Get().(T)
}

func (p *pool[T]) Put(v T) {
	if p.Finalizer != nil {
		p.Finalizer(&v)
	}
	p.Pool.Put(v)
}
