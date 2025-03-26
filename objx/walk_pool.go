package objx

import (
	"fmt"
	"sync"
)

type walkPool struct {
	sync.WaitGroup
	mu    sync.Mutex
	limit int
	task  chan func()
	stop  chan struct{}
}

func newWalkPool(limit int) *walkPool {
	ptr := &walkPool{
		limit: limit,
		task:  make(chan func(), limit),
		stop:  make(chan struct{}),
	}
	for i := 0; i < limit; i++ {
		go func() {
			for {
				select {
				case <-ptr.stop:
					return
				case fn, ok := <-ptr.task:
					if ok {
						func() {
							defer ptr.WaitGroup.Done()
							defer func() {
								if r := recover(); r != nil {
									fmt.Printf("future panic: %v \n", r)
								}
							}()
							fn()
						}()
					}
				}
			}
		}()
	}
	return ptr
}

func (p *walkPool) Go(fn func()) {
	p.WaitGroup.Add(1)
	p.task <- fn
}

func (p *walkPool) Destroy() {
	close(p.task)
	close(p.stop)
}

func (p *walkPool) Lock() {
	p.mu.Lock()
}

func (p *walkPool) Unlock() {
	p.mu.Unlock()
}
