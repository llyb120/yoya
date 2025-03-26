package syncx

import (
	"fmt"
	"sync"
	"time"

	"github.com/petermattis/goid"
)

var (
	futureHolder = &asyncHolder{
		mp:      make(map[int64]map[any]*future),
		indexMp: make(map[any]int64),
	}
)

type asyncHolder struct {
	mu sync.Mutex
	mp map[int64]map[any]*future
	// 某个future的协程索引
	indexMp map[any]int64
}

func (h *asyncHolder) save(ptrResult any, f *future) {
	h.mu.Lock()
	defer h.mu.Unlock()
	goid := goid.Get()
	if h.mp[goid] == nil {
		h.mp[goid] = make(map[any]*future)
	}
	h.mp[goid][ptrResult] = f
	h.indexMp[ptrResult] = goid
}

func (h *asyncHolder) contains(ptrResult any) bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	_, ok := h.indexMp[ptrResult]
	return ok
}

func (h *asyncHolder) loadAndDelete(ptrResult any) *future {
	h.mu.Lock()
	defer h.mu.Unlock()
	goid, ok := h.indexMp[ptrResult]
	if !ok {
		return nil
	}
	m, ok := h.mp[goid]
	if !ok {
		return nil
	}
	f, ok := m[ptrResult]
	if !ok {
		return nil
	}
	delete(m, ptrResult)
	if len(m) == 0 {
		delete(h.mp, goid)
	}
	delete(h.indexMp, ptrResult)
	return f
}

func (h *asyncHolder) loadAndDeleteWithGid() map[any]*future {
	h.mu.Lock()
	defer h.mu.Unlock()
	goid := goid.Get()
	mp := make(map[any]*future)
	for k, v := range h.mp[goid] {
		mp[k] = v
	}
	delete(h.mp, goid)
	return mp
}

func (h *asyncHolder) clean() {
	h.mu.Lock()
	defer h.mu.Unlock()
	if len(h.mp) == 0 {
		return
	}
	var newIndexMp = make(map[any]int64)
	for k, v := range h.indexMp {
		if h.mp[v] == nil {
			// 不存在，键不需要保留
			continue
		}
		f := h.mp[v][k]
		if f == nil || f.done.Load() || f.exprtime.Before(time.Now()) {
			delete(h.mp[v], k)
			if len(h.mp[v]) == 0 {
				delete(h.mp, v)
			}
			continue
		}
		newIndexMp[k] = v
	}
	h.indexMp = newIndexMp
}

func (h *asyncHolder) cleanGid() {
	goid := goid.Get()
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.mp, goid)
}

func async[T any](handler func() (T, error)) *T {
	future := &future{exprtime: time.Now().Add(5 * time.Minute)}
	var zero T
	ptrResult := &zero
	future.wg.Add(1)
	futureHolder.save(ptrResult, future)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				future.err = fmt.Errorf("future panic: %v", r)
			}
			future.wg.Done()
			future.done.Store(true)
		}()

		result, err := handler()
		if err != nil {
			future.err = err
		} else {
			*ptrResult = result
			future.result = result
		}
		// if len(result) > 0 {
		// 	for _, r := range result {
		// 		val := r.Interface()
		// 		if err, ok := val.(error); ok {
		// 			future.err = err
		// 		} else if value, ok := val.(T); ok {
		// 			*ptrResult = value
		// 			future.result = value
		// 		}
		// 	}
		// }
	}()

	return ptrResult
}

func Async_0_0(fn func()) func() *any {
	return func() *any {
		var asyncFn = async[any](func() (any, error) {
			fn()
			return nil, nil
		})
		return asyncFn
	}
}

func Async_0_1[T any](fn func() T) func() *T {
	return func() *T {
		var asyncFn = async[T](func() (T, error) {
			return fn(), nil
		})
		return asyncFn
	}
}

func Async_0_2[R0 any, R1 error](fn func() (R0, R1)) func() *R0 {
	return func() *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn()
		})
		return asyncFn
	}
}

func Async_1_0[T any](fn func(T)) func(T) *any {
	return func(t T) *any {
		var asyncFn = async[any](func() (any, error) {
			fn(t)
			return nil, nil
		})
		return asyncFn
	}
}

func Async_1_1[P0, R0 any](fn func(P0) R0) func(P0) *R0 {
	return func(p0 P0) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0), nil
		})
		return asyncFn
	}
}

func Async_1_2[P0, P1, R0 any, R1 error](fn func(P0, P1) (R0, error)) func(P0, P1) *R0 {
	return func(p0 P0, p1 P1) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1)
		})
		return asyncFn
	}
}

func Async_2_0[R0, R1 any](fn func(R0, R1)) func(R0, R1) *any {
	return func(r0 R0, r1 R1) *any {
		var asyncFn = async[any](func() (any, error) {
			fn(r0, r1)
			return nil, nil
		})
		return asyncFn
	}
}

func Async_2_1[P0, P1, R0 any](fn func(P0, P1) R0) func(P0, P1) *R0 {
	return func(p0 P0, p1 P1) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1), nil
		})
		return asyncFn
	}
}

func Async_2_2[P0, P1, R0 any, R1 error](fn func(P0, P1) (R0, R1)) func(P0, P1) *R0 {
	return func(p0 P0, p1 P1) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1)
		})
		return asyncFn
	}
}
func Async_3_0[P0, P1, P2 any](fn func(P0, P1, P2)) func(P0, P1, P2) *any {
	return func(p0 P0, p1 P1, p2 P2) *any {
		var asyncFn = async[any](func() (any, error) {
			fn(p0, p1, p2)
			return nil, nil
		})
		return asyncFn
	}
}

func Async_3_1[P0, P1, P2, R0 any](fn func(P0, P1, P2) R0) func(P0, P1, P2) *R0 {
	return func(p0 P0, p1 P1, p2 P2) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1, p2), nil
		})
		return asyncFn
	}
}

func Async_3_2[P0, P1, P2, R0 any, R1 error](fn func(P0, P1, P2) (R0, R1)) func(P0, P1, P2) *R0 {
	return func(p0 P0, p1 P1, p2 P2) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1, p2)
		})
		return asyncFn
	}
}

func Async_4_0[P0, P1, P2, P3 any](fn func(P0, P1, P2, P3)) func(P0, P1, P2, P3) *any {
	return func(p0 P0, p1 P1, p2 P2, p3 P3) *any {
		var asyncFn = async[any](func() (any, error) {
			fn(p0, p1, p2, p3)
			return nil, nil
		})
		return asyncFn
	}
}

func Async_4_1[P0, P1, P2, P3, R0 any](fn func(P0, P1, P2, P3) R0) func(P0, P1, P2, P3) *R0 {
	return func(p0 P0, p1 P1, p2 P2, p3 P3) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1, p2, p3), nil
		})
		return asyncFn
	}
}

func Async_4_2[P0, P1, P2, P3, R0 any, R1 error](fn func(P0, P1, P2, P3) (R0, R1)) func(P0, P1, P2, P3) *R0 {
	return func(p0 P0, p1 P1, p2 P2, p3 P3) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1, p2, p3)
		})
		return asyncFn
	}
}

func Async_5_0[P0, P1, P2, P3, P4 any](fn func(P0, P1, P2, P3, P4)) func(P0, P1, P2, P3, P4) *any {
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) *any {
		var asyncFn = async[any](func() (any, error) {
			fn(p0, p1, p2, p3, p4)
			return nil, nil
		})
		return asyncFn
	}
}

func Async_5_1[P0, P1, P2, P3, P4, R0 any](fn func(P0, P1, P2, P3, P4) R0) func(P0, P1, P2, P3, P4) *R0 {
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1, p2, p3, p4), nil
		})
		return asyncFn
	}
}

func Async_5_2[P0, P1, P2, P3, P4, R0 any, R1 error](fn func(P0, P1, P2, P3, P4) (R0, R1)) func(P0, P1, P2, P3, P4) *R0 {
	return func(p0 P0, p1 P1, p2 P2, p3 P3, p4 P4) *R0 {
		var asyncFn = async[R0](func() (R0, error) {
			return fn(p0, p1, p2, p3, p4)
		})
		return asyncFn
	}
}
