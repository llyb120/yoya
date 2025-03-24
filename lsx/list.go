package lsx

import (
	// "reflect"

	"github.com/llyb120/yoya/objx"
	"github.com/llyb120/yoya/stlx"
)

func Map[T any, R any](arr []T, fn func(T, int) (R, bool)) []R {
	result := make([]R, 0, len(arr))
	for i, v := range arr {
		r, ok := fn(v, i)
		if !ok {
			continue
		}
		result = append(result, r)
	}
	return result
}

func Filter[T any](arr *[]T, fn func(T, int) bool) {
	var result []T
	result = make([]T, 0, len(*arr))
	for i, v := range *arr {
		if fn(v, i) {
			result = append(result, v)
		}
	}
	*arr = result
}

func Find[T any](arr []T, fn func(T, int) bool) (T, bool) {
	index := FindIndex(arr, fn)
	if index == -1 {
		var zero T
		return zero, false
	}
	return arr[index], true
}

func FindIndex[T any](arr []T, fn func(T, int) bool) int {
	for i, v := range arr {
		if fn(v, i) {
			return i
		}
	}
	return -1
}

func Reduce[T any, R any](arr []T, fn func(R, T) R, initial R) R {
	var source []T
	result := initial
	source = arr
	for _, v := range source {
		result = fn(result, v)
	}
	return result
}

func For[T any](arr []T, fn func(T, int) bool) {
	var source []T
	source = arr
	for i, v := range source {
		c := fn(v, i)
		if !c {
			return
		}
	}
}

func Distinct[T any](arr *[]T, fn ...func(T, int) any) {
	var mp = make(map[any]bool)
	var result []T
	for i, v := range *arr {
		var k any
		if len(fn) > 0 {
			k = fn[0](v, i)
		} else {
			k = v
		}
		if mp[k] {
			continue
		}
		result = append(result, v)
		mp[k] = true
	}
	*arr = result
}

func Sort[T any](arr *[]T, less func(T, T) bool) {
	var cp []T
	cp = make([]T, len(*arr))
	copy(cp, *arr)
	cp = timSort(cp, less)
	*arr = cp
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

func Group[T any](arr []T, fn func(T) any) [][]T {
	var result = stlx.NewMultiMap[any, T]()
	for _, v := range arr {
		k := fn(v)
		if k == nil {
			continue
		}
		result.Set(k, v)
	}
	return result.Vals()
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
