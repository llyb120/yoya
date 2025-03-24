package objx

import (
	"fmt"

	reflect "github.com/goccy/go-reflect"
)

func Assign[T comparable, K0 any, K1 any](dest map[T]K0, source map[T]K1) {
	var zero0 K0
	var zero1 K1
	isSameType := reflect.TypeOf(zero0) == reflect.TypeOf(zero1)
	if isSameType {
		for k, v := range source {
			dest[k] = any(v).(K0)
		}
		return
	} else {
		// 一次性转过来
		var mp map[T]K0
		if err := Cast(source, &mp); err != nil {
			fmt.Println("err:", err)
			return
		}
		for k, v := range mp {
			dest[k] = v
		}
		return
	}
}
