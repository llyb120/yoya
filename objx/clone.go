package objx

import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

// DeepClone 深拷贝任意对象，包括私有字段
// 支持所有Go类型：基本类型、指针、切片、映射、结构体、接口、通道、函数等
// 自动处理循环引用问题
func DeepClone[T any](src T) (T, error) {
	var result T
	err := safeCall(func() error {
		visited := make(map[uintptr]reflect.Value)
		srcValue := reflect.ValueOf(src)
		clonedValue := deepCloneValue(srcValue, visited)

		if clonedValue.IsValid() && clonedValue.Type().AssignableTo(reflect.TypeOf(result)) {
			result = clonedValue.Interface().(T)
		}
		return nil
	})
	return result, err
}

// DeepCloneAny 深拷贝任意类型的对象
func DeepCloneAny(src any) (any, error) {
	var result any
	err := safeCall(func() error {
		visited := make(map[uintptr]reflect.Value)
		srcValue := reflect.ValueOf(src)
		clonedValue := deepCloneValue(srcValue, visited)

		if clonedValue.IsValid() {
			result = clonedValue.Interface()
		}
		return nil
	})
	return result, err
}

// MustDeepClone 深拷贝，如果出错则panic
func MustDeepClone[T any](src T) T {
	result, err := DeepClone(src)
	if err != nil {
		panic(err)
	}
	return result
}

// MustDeepCloneAny 深拷贝任意类型，如果出错则panic
func MustDeepCloneAny(src any) any {
	result, err := DeepCloneAny(src)
	if err != nil {
		panic(err)
	}
	return result
}

// safeCall 安全调用函数，捕获panic并转换为error
func safeCall(fn func() error) (err error) {
	defer func() {
		if r := recover(); r != nil {
			stack := make([]byte, 4096)
			stackLen := runtime.Stack(stack, false)
			err = fmt.Errorf("panic: %v\nstack: %s", r, stack[:stackLen])
		}
	}()
	return fn()
}

// deepCloneValue 递归深拷贝reflect.Value
func deepCloneValue(src reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	if !src.IsValid() {
		return reflect.Value{}
	}

	srcType := src.Type()

	switch src.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128, reflect.String:
		// 基本类型直接复制
		return src

	case reflect.Array:
		return cloneArray(src, visited)

	case reflect.Slice:
		return cloneSlice(src, visited)

	case reflect.Map:
		return cloneMap(src, visited)

	case reflect.Ptr:
		return clonePtr(src, visited)

	case reflect.Struct:
		return cloneStruct(src, visited)

	case reflect.Interface:
		return cloneInterface(src, visited)

	case reflect.Chan:
		return cloneChan(src, visited)

	case reflect.Func:
		// 函数类型直接返回原值（函数无法深拷贝）
		return src

	case reflect.UnsafePointer:
		// unsafe.Pointer 直接复制
		return src

	default:
		// 其他类型直接返回零值
		return reflect.Zero(srcType)
	}
}

// cloneArray 拷贝数组
func cloneArray(src reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	srcType := src.Type()
	dst := reflect.New(srcType).Elem()

	for i := 0; i < src.Len(); i++ {
		elemClone := deepCloneValue(src.Index(i), visited)
		if elemClone.IsValid() {
			dst.Index(i).Set(elemClone)
		}
	}

	return dst
}

// cloneSlice 拷贝切片
func cloneSlice(src reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	if src.IsNil() {
		return reflect.Zero(src.Type())
	}

	srcType := src.Type()
	dst := reflect.MakeSlice(srcType, src.Len(), src.Cap())

	for i := 0; i < src.Len(); i++ {
		elemClone := deepCloneValue(src.Index(i), visited)
		if elemClone.IsValid() {
			dst.Index(i).Set(elemClone)
		}
	}

	return dst
}

// cloneMap 拷贝映射
func cloneMap(src reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	if src.IsNil() {
		return reflect.Zero(src.Type())
	}

	srcType := src.Type()
	dst := reflect.MakeMap(srcType)

	for _, key := range src.MapKeys() {
		value := src.MapIndex(key)

		// 对于从unexported字段获取的键和值，我们需要创建新的可用值
		var keyToUse, valueToUse reflect.Value

		// 处理键
		if key.CanInterface() {
			keyToUse = reflect.ValueOf(key.Interface())
		} else {
			keyToUse = key
		}

		// 处理值
		if value.CanInterface() {
			valueToUse = reflect.ValueOf(value.Interface())
		} else {
			valueToUse = value
		}

		// 深拷贝键和值
		keyClone := deepCloneValue(keyToUse, visited)
		valueClone := deepCloneValue(valueToUse, visited)

		if keyClone.IsValid() && valueClone.IsValid() {
			dst.SetMapIndex(keyClone, valueClone)
		}
	}

	return dst
}

// clonePtr 拷贝指针
func clonePtr(src reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	if src.IsNil() {
		return reflect.Zero(src.Type())
	}

	// 检查循环引用
	addr := src.Pointer()
	if cached, exists := visited[addr]; exists {
		return cached
	}

	srcType := src.Type()
	elemType := srcType.Elem()

	// 创建新的指针
	dst := reflect.New(elemType)
	visited[addr] = dst

	// 递归拷贝指针指向的值
	elemClone := deepCloneValue(src.Elem(), visited)
	if elemClone.IsValid() {
		dst.Elem().Set(elemClone)
	}

	return dst
}

// cloneStruct 拷贝结构体，包括私有字段
func cloneStruct(src reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	srcType := src.Type()
	dst := reflect.New(srcType).Elem()

	for i := 0; i < src.NumField(); i++ {
		srcField := src.Field(i)
		dstField := dst.Field(i)
		fieldType := srcType.Field(i)

		// 跳过不可访问的字段
		if !srcField.IsValid() {
			continue
		}

		// 如果是可导出字段且可设置，直接拷贝
		if fieldType.IsExported() && dstField.CanSet() {
			fieldClone := deepCloneValue(srcField, visited)
			if fieldClone.IsValid() {
				dstField.Set(fieldClone)
			}
		} else {
			// 私有字段使用unsafe拷贝
			if dst.CanAddr() {
				dstFieldPtr := unsafe.Pointer(dst.UnsafeAddr() + fieldType.Offset)

				if src.CanAddr() {
					// 源可寻址，使用unsafe访问
					srcFieldPtr := unsafe.Pointer(src.UnsafeAddr() + fieldType.Offset)
					srcFieldValue := reflect.NewAt(fieldType.Type, srcFieldPtr).Elem()
					dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()

					fieldClone := deepCloneValue(srcFieldValue, visited)
					if fieldClone.IsValid() {
						dstFieldValue.Set(fieldClone)
					}
				} else {
					// 源不可寻址，根据字段类型处理
					switch fieldType.Type.Kind() {
					case reflect.Bool:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(srcField.Bool())
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Int:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(int(srcField.Int()))
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Int8:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(int8(srcField.Int()))
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Int16:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(int16(srcField.Int()))
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Int32:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(int32(srcField.Int()))
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Int64:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(srcField.Int())
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Uint:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(uint(srcField.Uint()))
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Uint8:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(uint8(srcField.Uint()))
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Uint16:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(uint16(srcField.Uint()))
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Uint32:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(uint32(srcField.Uint()))
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Uint64:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(srcField.Uint())
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Uintptr:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(uintptr(srcField.Uint()))
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Float32:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(float32(srcField.Float()))
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Float64:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(srcField.Float())
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Complex64:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(complex64(srcField.Complex()))
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.Complex128:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(srcField.Complex())
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					case reflect.String:
						dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
						newValue := reflect.ValueOf(srcField.String())
						if newValue.IsValid() {
							dstFieldValue.Set(newValue)
						}
					default:
						// 复杂类型，尝试递归拷贝
						fieldClone := deepCloneValue(srcField, visited)
						if fieldClone.IsValid() {
							dstFieldValue := reflect.NewAt(fieldType.Type, dstFieldPtr).Elem()
							// 使用unsafe设置字段值
							unsafeSetFieldValue(dstFieldValue, fieldClone)
						}
					}
				}
			}
		}
	}

	return dst
}

// unsafeSetFieldValue 使用unsafe设置字段值，绕过Go的访问限制
func unsafeSetFieldValue(field reflect.Value, value reflect.Value) {
	// 如果目标字段和值的类型相同，则使用unsafe直接设置
	if field.Type() == value.Type() {
		if field.CanAddr() && value.CanAddr() {
			// 获取目标字段的指针
			fieldPtr := unsafe.Pointer(field.UnsafeAddr())
			// 获取值的指针
			valuePtr := unsafe.Pointer(value.UnsafeAddr())
			// 计算字段大小
			size := field.Type().Size()
			// 直接复制内存
			if size > 0 {
				copyMemory(fieldPtr, valuePtr, size)
			}
		} else {
			// 回退到使用标准反射
			field.Set(value)
		}
	} else {
		// 类型不同，回退到使用标准反射
		field.Set(value)
	}
}

// copyMemory 复制内存
func copyMemory(dst, src unsafe.Pointer, size uintptr) {
	// 转为切片进行复制
	dstSlice := (*[1 << 30]byte)(dst)[:size:size]
	srcSlice := (*[1 << 30]byte)(src)[:size:size]
	copy(dstSlice, srcSlice)
}

// cloneInterface 拷贝接口
func cloneInterface(src reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	if src.IsNil() {
		return reflect.Zero(src.Type())
	}

	// 获取接口中的实际值
	elem := src.Elem()
	elemClone := deepCloneValue(elem, visited)

	if !elemClone.IsValid() {
		return reflect.Zero(src.Type())
	}

	// 创建新的接口值
	dst := reflect.New(src.Type()).Elem()
	dst.Set(elemClone)

	return dst
}

// cloneChan 拷贝通道
func cloneChan(src reflect.Value, visited map[uintptr]reflect.Value) reflect.Value {
	if src.IsNil() {
		return reflect.Zero(src.Type())
	}

	srcType := src.Type()
	// 创建相同类型和缓冲区大小的新通道
	dst := reflect.MakeChan(srcType, src.Cap())

	// 注意：这里不复制通道中的数据，因为通道的内容是动态的
	// 如果需要复制通道中的数据，需要额外的逻辑来处理
	// 通道的拷贝主要是创建一个新的相同类型的通道

	return dst
}
