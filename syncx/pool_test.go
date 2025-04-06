package syncx

import (
	"fmt"
	"testing"
)

func TestPool(t *testing.T) {
	var pool = NewPool(PoolOption[[]string]{
		New: func() []string {
			return []string{}
		},
	})

	res, rec := pool.Get()
	defer rec()
	res = append(res, "1", "2", "3")
	fmt.Println(res)
}

func BenchmarkPool(b *testing.B) {
	var pool = NewPool(PoolOption[[]string]{
		New: func() []string {
			return []string{}
		},
		Finalizer: func(v *[]string) {
			*v = (*v)[:0]
		},
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, rec := pool.Get()
		defer rec()
		res = append(res, "1", "2", "3")
	}

	res, _ := pool.Get()
	fmt.Println(res)
}
