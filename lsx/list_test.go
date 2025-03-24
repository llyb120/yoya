package lsx

import (
	"fmt"
	"regexp"
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
	result := Distinct(arr)
	if len(result) != 5 {
		t.Errorf("Distinct result length is not equal to input length")
	}
	Distinct(&arr)
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
	for _, v := range arr {
		fmt.Println(v)
	}
}

type TestStruct struct {
	Field0 string
	Field1 string
	Field2 string

	NormalField  string
	NormalField1 string
}

func TestMock2(t *testing.T) {
	arr := []TestStruct{
		{
			Field0:       "1",
			Field1:       "2",
			Field2:       "3",
			NormalField:  "4",
			NormalField1: "5",
		},
		{
			Field0:       "6",
			Field1:       "7",
			Field2:       "8",
			NormalField:  "9",
			NormalField1: "10",
		},
	}
	var re = regexp.MustCompile(`^Field\d+$`)
	if err := Mock(&arr, func(arr *[]map[string]any) {
		for i, _ := range *arr {
			for k, v := range (*arr)[i] {
				if re.MatchString(k) {
					(*arr)[i][k] = v.(string) + " hey!"
				}
			}
		}
	}); err != nil {
		t.Errorf("Mock error: %v", err)
	}
	for _, v := range arr {
		fmt.Println(v)
	}
}

func TestGroup(t *testing.T) {
	arr := []int{
		1, 2, 3, 4, 5, 1, 2, 3, 4, 5,
	}
	i := 0
	result := Group(arr, func(v int) any {
		defer func() { i++ }()
		return int(i / 3)
	})
	fmt.Println(result)
}
