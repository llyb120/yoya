package lsx

type iterable[T any] interface {
	[]T | *[]T
}

func Map[T any, R any](arr []T, fn func(T, int) R) []R {
	result := make([]R, len(arr))
	for i, v := range arr {
		result[i] = fn(v, i)
	}
	return result
}

func Filter[T any, K iterable[T]](arr K, fn func(T, int) bool) K {
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
	for i, v := range _arr {
		if fn(v, i) {
			result = append(result, v)
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

func For[T any, K iterable[T]](arr K, fn func(T, int)) {
	switch any(arr).(type) {
	case []T:
		for i, v := range any(arr).([]T) {
			fn(v, i)
		}
	case *[]T:
		for i, v := range *any(arr).(*[]T) {
			fn(v, i)
		}
	}
}

func Distinct[T comparable, K iterable[T]](arr K) K {
	var source []T
	var isPtr bool
	switch any(arr).(type) {
	case []T:
		source = any(arr).([]T)
	case *[]T:
		source = *any(arr).(*[]T)
		isPtr = true
	}
	var mp = make(map[T]bool)
	result := make([]T, 0)
	for _, v := range source {
		if !mp[v] {
			mp[v] = true
			result = append(result, v)
		}
	}
	if isPtr {
		*any(arr).(*[]T) = result
		return arr
	}
	return any(result).(K)
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
