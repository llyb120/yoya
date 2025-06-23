package internal

import (
	"reflect"
	"unsafe"
)

func UnsafeSetFieldValue(field reflect.Value, value reflect.Value, forceCheckType bool) {
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
				typedmemmove(field.Type(), fieldPtr, valuePtr)
			}
		} else {
			// 类型不同，回退到使用标准反射
			field.Set(value)
		}
	} else if !forceCheckType {
		// 类型不同，回退到使用标准反射
		field.Set(value)
	}
}

// 通过unsafe实现类似reflect.typedmemmove的功能
// 这个函数假设src和dst都有效，且大小匹配
func typedmemmove(typ reflect.Type, dst, src unsafe.Pointer) {
	// 使用非导出的内存复制
	// 实际实现中，我们应该参考runtime.typedmemmove的实现
	// 这里为了简化，直接使用标准库的copy
	size := typ.Size()
	memmove(dst, src, size)
}

// 封装底层内存复制操作
func memmove(dst, src unsafe.Pointer, size uintptr) {
	// 转为切片进行复制
	dstSlice := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(dst),
		Len:  int(size),
		Cap:  int(size),
	}))

	srcSlice := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(src),
		Len:  int(size),
		Cap:  int(size),
	}))

	copy(dstSlice, srcSlice)
}
