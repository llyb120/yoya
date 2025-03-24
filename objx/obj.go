package objx

import (
	reflect "github.com/goccy/go-reflect"
)

func Assign[T comparable, K0 any, K1 any](dest map[T]K0, source map[T]K1) {
	var zero0 K0
	var zero1 K1
	isSameType := reflect.TypeOf(zero0) == reflect.TypeOf(zero1)
	for k, v := range source {
		v := v
		if isSameType {
			dest[k] = any(v).(K0)
		} else {
			var zero K0
			err := Cast(v, &zero)
			if err != nil {
				continue
			}
			dest[k] = zero
		}
	}
}
