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
	Walk(item, func(s, k any, v any) any {
		if k == "Name" {
			var fn = syncx.Async_0_1(func() string {
				time.Sleep(2 * time.Second)
				return "ok"
			})
			res := fn()
			if err := syncx.Await(res); err != nil {
				return nil
			}
			return res
		}
		return nil
	}, Async, 1 * Level)

	fmt.Println(item)
	fmt.Println(time.Since(now))
}

func TestBreakWalk(t *testing.T) {
	type Person struct {
		Name  string
		Age   int
		Child []Person
	}
	var item = map[string][]Person{
		"name": {{Name: "张三", Age: 30, Child: []Person{{Name: "张三儿子", Age: 30}, {Name: "张三女儿", Age: 25}}}, {Name: "李四", Age: 25}},
	}
	Walk(item, func(s, k any, v any) any {
		fmt.Println(k, v)
		return nil
	})

	println("---------------------------------------")

	Walk(item, func(s, k any, v any) any {
		if k == "Child" {
			return BreakWalkSelf
		}
		fmt.Println(k, v)
		return nil
	})
}

func TestEnsure(t *testing.T) {
	var a any
	var b any
	var c any
	a = 1
	b = 2
	c = 3

	var d int
	var e int
	var f int
	Ensure(&d, a, &e, b, &f, c)
	fmt.Println(d, e, f)
}

func TestCollect(t *testing.T) {
	src := "name[age=10,name*=张三] id [1,2,3] name[age=10,name='张三']"
	selector := &selector{
		src: src,
	}
	nodes := selector.parse()
	fmt.Printf("%+v\n", nodes)

	data := map[string]any{
		"name": []any{
			map[string]any{
				"age":  10,
				"name": "张三",
				"id":   1,
			},
			map[string]any{
				"age":  10,
				"name": "张三",
				"id":   2,
			},
		},
	}
	results := Pick[string](data, "name [age=10,name='张三'] id")
	fmt.Printf("%+v\n", results)

	reuslt2 := Pick[any](data, "id")
	fmt.Printf("%+v\n", reuslt2)
}
