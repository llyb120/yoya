package internal

var _converter = newConverter()

// 普通类型转换
func Cast(dest any, src any) error {
	// tp := reflect.TypeOf(*dest)
	// if tp.Kind() == reflect.Slice {
	// 	return _converter.ConvertSlice(src, dest)
	// }
	return _converter.Convert(src, dest)
}
