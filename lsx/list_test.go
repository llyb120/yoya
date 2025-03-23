package lsx

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	result := Map(arr, func(v int) (int, bool) {
		return v * 2, true
	})
	if len(result) != len(arr) {
		t.Errorf("Map result length is not equal to input length")
	}
	for i, v := range result {
		if v != arr[i]*2 {
			t.Errorf("Map result is not equal to input")
		}
	}
}

func TestDistinct(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}
	// 居然不能自动推导？？
	result := Distinct[int](arr)
	if len(result) != 5 {
		t.Errorf("Distinct result length is not equal to input length")
	}
	Distinct[int](&arr)
	if len(arr) != 5 {
		t.Errorf("Distinct result length is not equal to input length")
	}
}

func TestFilter(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5, 1, 2, 3, 4, 5}
	result := Filter(arr, func(v int) bool {
		return v%2 == 0
	})
	if len(result) != 4 {
		t.Errorf("Filter result length is not equal to input length")
	}
}

func TestSort(t *testing.T) {
	arr := []int{5, 4, 3, 2, 1}
	Sort(&arr, func(a, b int) bool {
		return a < b
	})
	for _, v := range arr {
		fmt.Println(v)
	}
}

func ttt(arr []int) {
	for i, _ := range arr {
		arr[i] = arr[i] * 2
		fmt.Println(i, arr[i])
	}
}

func TestMock(t *testing.T) {
	arr := []string{"1", "2", "3", "4", "5"}
	err := Mock(&arr, func(arr *[]int) {
		ttt(*arr)
	})
	if err != nil {
		t.Errorf("Mock error: %v", err)
	}
}
