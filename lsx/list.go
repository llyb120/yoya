package lsx

import (
	// "reflect"
	reflect "github.com/goccy/go-reflect"
	"github.com/llyb120/yoya/objx"
)

type iterable[T any] interface {
	[]T | *[]T
}

func Map[T any, R any](arr []T, fn func(T) (R, bool)) []R {
	result := make([]R, 0, len(arr))
	for _, v := range arr {
		r, ok := fn(v)
		if !ok {
			continue
		}
		result = append(result, r)
	}
	return result
}

func Filter[T any, K iterable[T], M func(T) bool](arr K, fn M) K {
	var result []T
	var _arr []T
	var isPtr bool
	switch any(arr).(type) {
	case []T:
		_arr = any(arr).([]T)
	case *[]T:
		_arr = *any(arr).(*[]T)
		isPtr = true
	}
	result = make([]T, 0, len(_arr))
	for i, v := range _arr {
		switch fn := any(fn).(type) {
		case func(T) bool:
			if fn(v) {
				result = append(result, v)
			}
		case func(T, int) bool:
			if fn(v, i) {
				result = append(result, v)
			}
		}
	}
	if isPtr {
		*any(arr).(*[]T) = result
		return arr
	}
	return any(result).(K)
}

func Find[T any, K iterable[T]](arr K, fn func(T, int) bool) (T, bool) {
	index := FindIndex(arr, fn)
	if index == -1 {
		var zero T
		return zero, false
	}
	return any(arr).([]T)[index], true
}

func FindIndex[T any, K iterable[T]](arr K, fn func(T, int) bool) int {
	switch any(arr).(type) {
	case []T:
		for i, v := range any(arr).([]T) {
			if fn(v, i) {
				return i
			}
		}
	case *[]T:
		for i, v := range *any(arr).(*[]T) {
			if fn(v, i) {
				return i
			}
		}
	}
	return -1
}

func Reduce[T any, K iterable[T], R any](arr K, fn func(R, T) R, initial R) R {
	var source []T
	result := initial
	switch any(arr).(type) {
	case []T:
		source = any(arr).([]T)
	case *[]T:
		source = *any(arr).(*[]T)
	}
	for _, v := range source {
		result = fn(result, v)
	}
	return result
}

func For[T any, K iterable[T]](arr K, fn func(T, int) bool) {
	switch any(arr).(type) {
	case []T:
		for i, v := range any(arr).([]T) {
			v := v
			i := i
			c := fn(v, i)
			if !c {
				return
			}
		}
	case *[]T:
		for i, v := range *any(arr).(*[]T) {
			v := v
			i := i
			c := fn(v, i)
			if !c {
				return
			}
		}
	}
}

func Distinct[T any](arr T) T {
	// var source []T
	var isPtr bool
	tp := reflect.TypeOf(arr)
	val := reflect.ValueOf(arr)
	if tp.Kind() == reflect.Ptr {
		isPtr = true
		val = val.Elem()
	}
	if val.Kind() != reflect.Slice {
		var zero T
		return zero
	}
	var source = reflect.MakeSlice(val.Type(), 0, val.Len())
	var mp = make(map[any]bool)
	for i := 0; i < val.Len(); i++ {
		refV := val.Index(i)
		v := refV.Interface()
		if mp[v] {
			continue
		}
		source = reflect.Append(source, refV)
		mp[v] = true
	}
	if isPtr {
		val.Set(source)
		return arr
	} else {
		return source.Interface().(T)
	}
}

func Sort[T any, K iterable[T]](arr K, less func(T, T) bool) K {
	var cp []T
	var isPtr bool
	switch any(arr).(type) {
	case []T:
		cp = make([]T, len(cp))
		copy(cp, any(arr).([]T))
	case *[]T:
		cp = make([]T, len(*any(arr).(*[]T)))
		copy(cp, *any(arr).(*[]T))
		isPtr = true
	}
	cp = timSort(cp, less)
	if isPtr {
		*any(arr).(*[]T) = cp
		return arr
	}
	return any(cp).(K)
}

func Keys[K comparable, V any](mp map[K]V) []K {
	result := make([]K, 0, len(mp))
	for k := range mp {
		result = append(result, k)
	}
	return result
}

func Vals[K comparable, V any](mp map[K]V) []V {
	result := make([]V, 0, len(mp))
	for _, v := range mp {
		result = append(result, v)
	}
	return result
}

func Mock[K any, T any](arr *[]K, fn func(*[]T)) error {
	var mock []T
	err := objx.Cast(arr, &mock)
	if err != nil {
		return err
	}
	fn(&mock)
	// 还原
	var result []K
	err = objx.Cast(mock, &result)
	if err != nil {
		return err
	}
	*arr = result
	return nil
}

func timSort[T any](arr []T, less func(T, T) bool) []T {
	if len(arr) <= 1 {
		return arr
	}

	const minRun = 32
	const maxRun = 64

	// 计算最小运行长度
	runLength := minRun
	if len(arr) < minRun {
		runLength = len(arr)
	}

	// 插入排序辅助函数
	insertionSort := func(arr []T, left, right int) {
		for i := left + 1; i <= right; i++ {
			key := arr[i]
			j := i - 1
			for j >= left && less(key, arr[j]) {
				arr[j+1] = arr[j]
				j--
			}
			arr[j+1] = key
		}
	}

	// 合并两个有序子数组
	merge := func(arr []T, l, m, r int) {
		left := make([]T, m-l+1)
		right := make([]T, r-m)
		copy(left, arr[l:m+1])
		copy(right, arr[m+1:r+1])

		i, j, k := 0, 0, l
		for i < len(left) && j < len(right) {
			if !less(right[j], left[i]) {
				arr[k] = left[i]
				i++
			} else {
				arr[k] = right[j]
				j++
			}
			k++
		}

		for i < len(left) {
			arr[k] = left[i]
			i++
			k++
		}

		for j < len(right) {
			arr[k] = right[j]
			j++
			k++
		}
	}

	// 对小块使用插入排序
	for i := 0; i < len(arr); i += runLength {
		end := i + runLength - 1
		if end >= len(arr) {
			end = len(arr) - 1
		}
		insertionSort(arr, i, end)
	}

	// 合并排序后的小块
	for size := runLength; size < len(arr); size *= 2 {
		for start := 0; start < len(arr); start += size * 2 {
			mid := start + size - 1
			end := start + size*2 - 1
			if mid >= len(arr) {
				mid = len(arr) - 1
			}
			if end >= len(arr) {
				end = len(arr) - 1
			}
			if mid > start {
				merge(arr, start, mid, end)
			}
		}
	}

	return arr
}
