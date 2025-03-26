package objx

import (
	"fmt"
	"testing"
)

func TestWalk(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}
	var item = map[string][]Person{
		"name": {{Name: "张三", Age: 30}, {Name: "李四", Age: 25}},
	}
	Walk(item, func(k any, v any) any {
		if k == "Name" {
			return "fuck"
		}
		return nil
	})

	fmt.Println(item)
}
