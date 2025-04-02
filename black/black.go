package black

import (
	"bytes"
	"encoding/gob"
	"errors"
	"reflect"
	"unsafe"
)

func Byte2Str(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func ToBytes(obj any) ([]byte, error) {
	// 获取对象的值和类型
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Struct {
		return nil, errors.New("struct must be pointer")
	}

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice:
		return sliceToBytes(obj)
	case reflect.Struct:
		return structToBytes(v)
	case reflect.Map:
		return mapToBytes(obj)
	}

	return nil, errors.New("invalid data size")
}

func FromBytes[T any](bs []byte) (T, error) {
	v := reflect.ValueOf(new(T))
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice:
		return bytesToSlice[T](bs)
	case reflect.Struct:
		return bytesToStruct[T](bs, v)
	case reflect.Map:
		return bytesToMap[T](bs)
	}

	var zero T
	return zero, errors.New("invalid data size")
}

func structToBytes(v reflect.Value) ([]byte, error) {
	// 计算内存大小
	size := int(v.Type().Size())

	// 构造字节切片
	var buf []byte
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&buf))
	sh.Data = uintptr(unsafe.Pointer(v.UnsafeAddr()))
	sh.Len = size
	sh.Cap = size

	// 复制内存内容到新切片（避免返回指向原对象的引用）
	result := make([]byte, size)
	copy(result, buf)
	return result, nil
}

func mapToBytes(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func bytesToMap[T any](data []byte) (T, error) {
	var s T
	dec := gob.NewDecoder(bytes.NewReader(data))
	return s, dec.Decode(&s)
}

func sliceToBytes(v any) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func bytesToSlice[T any](data []byte) (T, error) {
	var s T
	dec := gob.NewDecoder(bytes.NewReader(data))
	return s, dec.Decode(&s)
}

func bytesToStruct[T any](data []byte, v reflect.Value) (T, error) {
	sh := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	sh.Data = uintptr(unsafe.Pointer(v.UnsafeAddr()))
	sh.Len = len(data)
	sh.Cap = len(data)
	return *(*T)(unsafe.Pointer(sh.Data)), nil
}
