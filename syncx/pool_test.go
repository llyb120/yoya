package syncx

import (
	"fmt"
	"testing"
)

func TestPool(t *testing.T) {
	var pool = PoolOption[[]string]{
		New: func() []string {
			return []string{}
		},
	}.Build()

	res, rec := pool.Get()
	defer rec()
	res = append(res, "1", "2", "3")
	fmt.Println(res)
}

func BenchmarkPool(b *testing.B) {
	var pool = PoolOption[[]string]{
		New: func() []string {
			return []string{}
		},
		Finalizer: func(v *[]string) {
			*v = (*v)[:0]
		},
	}.Build()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, rec := pool.Get()
		defer rec()
		res = append(res, "1", "2", "3")
	}

	res, _ := pool.Get()
	fmt.Println(res)
}
