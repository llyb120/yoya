package objx

import (
	"fmt"
	"testing"
	"time"

	"github.com/llyb120/yoya/syncx"
)

func TestWalk(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	var item = map[string][]Person{
		"name": {{Name: "张三", Age: 30}, {Name: "李四", Age: 25}},
	}
	now := time.Now()
	Walk(item, func(s, k any, v any) syncx.AsyncFn {
		if k == "Name" {
			return syncx.Async_0_1(func() string {
				time.Sleep(5 * time.Second)
				return "fuck"
			})()
		}
		return nil
	})

	fmt.Println(item)
	fmt.Println(time.Since(now))
}
