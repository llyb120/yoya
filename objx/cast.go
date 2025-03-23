package objx

import (
	reflect "github.com/goccy/go-reflect"
)

var _converter = newConverter()

// 普通类型转换
func Cast[T any](src any, dest *T) error {
	tp := reflect.TypeOf(*dest)
	if tp.Kind() == reflect.Slice {
		return _converter.ConvertSlice(src, dest)
	}
	return _converter.Convert(src, dest)
}
