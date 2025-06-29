package objx

import (
	"fmt"
	"reflect"
	"unsafe"
)

// 示例：基本类型的深拷贝
func ExampleDeepClone_basicTypes() {
	// 基本类型
	original := 42
	cloned, _ := DeepClone(original)
	fmt.Printf("Original: %d, Cloned: %d\n", original, cloned)

	// 字符串
	str := "hello world"
	clonedStr, _ := DeepClone(str)
	fmt.Printf("Original: %s, Cloned: %s\n", str, clonedStr)

	// Output:
	// Original: 42, Cloned: 42
	// Original: hello world, Cloned: hello world
}

// 示例：切片的深拷贝
func ExampleDeepClone_slice() {
	original := []int{1, 2, 3, 4, 5}
	cloned, _ := DeepClone(original)

	fmt.Printf("Original: %v\n", original)
	fmt.Printf("Cloned: %v\n", cloned)

	// 修改原始切片
	original[0] = 999
	fmt.Printf("After modifying original[0] to 999:\n")
	fmt.Printf("Original: %v\n", original)
	fmt.Printf("Cloned: %v (unchanged)\n", cloned)

	// Output:
	// Original: [1 2 3 4 5]
	// Cloned: [1 2 3 4 5]
	// After modifying original[0] to 999:
	// Original: [999 2 3 4 5]
	// Cloned: [1 2 3 4 5] (unchanged)
}

// 示例：包含私有字段的结构体深拷贝
func ExampleDeepClone_structWithPrivateFields() {
	type Person struct {
		Name   string  // 公共字段
		age    int     // 私有字段
		salary float64 // 私有字段
	}

	original := Person{
		Name:   "张三",
		age:    30,
		salary: 50000.0,
	}

	cloned, _ := DeepClone(original)

	fmt.Printf("Original: Name=%s\n", original.Name)
	fmt.Printf("Cloned: Name=%s\n", cloned.Name)

	// 使用unsafe访问私有字段验证拷贝
	originalVal := reflect.ValueOf(&original).Elem()
	clonedVal := reflect.ValueOf(&cloned).Elem()

	ageField, _ := originalVal.Type().FieldByName("age")
	originalAge := *(*int)(unsafe.Pointer(originalVal.UnsafeAddr() + ageField.Offset))
	clonedAge := *(*int)(unsafe.Pointer(clonedVal.UnsafeAddr() + ageField.Offset))

	fmt.Printf("Original age (private): %d\n", originalAge)
	fmt.Printf("Cloned age (private): %d\n", clonedAge)

	// Output:
	// Original: Name=张三
	// Cloned: Name=张三
	// Original age (private): 30
	// Cloned age (private): 30
}

// 示例：复杂嵌套结构的深拷贝
func ExampleDeepClone_complexStruct() {
	type Address struct {
		City    string
		zipCode int
	}

	type Person struct {
		Name     string
		Age      int
		Address  *Address
		Hobbies  []string
		Contacts map[string]string // 改为公共字段
	}

	original := Person{
		Name: "李四",
		Age:  25,
		Address: &Address{
			City:    "北京",
			zipCode: 100000,
		},
		Hobbies:  []string{"读书", "游泳", "编程"},
		Contacts: map[string]string{"email": "lisi@example.com"},
	}

	cloned, _ := DeepClone(original)

	fmt.Printf("Original: %s, %d岁, 住在%s\n",
		original.Name, original.Age, original.Address.City)
	fmt.Printf("Cloned: %s, %d岁, 住在%s\n",
		cloned.Name, cloned.Age, cloned.Address.City)

	// 修改原始数据
	original.Address.City = "上海"
	original.Hobbies[0] = "电影"

	fmt.Printf("After modification:\n")
	fmt.Printf("Original address: %s\n", original.Address.City)
	fmt.Printf("Cloned address: %s (unchanged)\n", cloned.Address.City)
	fmt.Printf("Original hobbies: %v\n", original.Hobbies)
	fmt.Printf("Cloned hobbies: %v (unchanged)\n", cloned.Hobbies)

	// Output:
	// Original: 李四, 25岁, 住在北京
	// Cloned: 李四, 25岁, 住在北京
	// After modification:
	// Original address: 上海
	// Cloned address: 北京 (unchanged)
	// Original hobbies: [电影 游泳 编程]
	// Cloned hobbies: [读书 游泳 编程] (unchanged)
}

// 示例：使用 MustDeepClone
func ExampleMustDeepClone() {
	data := map[string][]int{
		"group1": {1, 2, 3},
		"group2": {4, 5, 6},
	}

	// MustDeepClone 在出错时会panic，适合确定不会出错的场景
	cloned := MustDeepClone(data)

	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cloned: %v\n", cloned)

	// 修改原始数据
	data["group1"][0] = 999
	delete(data, "group2")

	fmt.Printf("After modification:\n")
	fmt.Printf("Original: %v\n", data)
	fmt.Printf("Cloned: %v (unchanged)\n", cloned)

	// Output:
	// Original: map[group1:[1 2 3] group2:[4 5 6]]
	// Cloned: map[group1:[1 2 3] group2:[4 5 6]]
	// After modification:
	// Original: map[group1:[999 2 3]]
	// Cloned: map[group1:[1 2 3] group2:[4 5 6]] (unchanged)
}
