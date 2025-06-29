package objx

import (
	"reflect"
	"testing"
	"unsafe"
)

// 测试结构体，包含私有字段
type TestStruct struct {
	PublicField   string
	privateField  int
	privatePtr    *string
	PublicSlice   []int
	privateMap    map[string]int
	PublicNested  *NestedStruct
	privateNested *NestedStruct
}

type NestedStruct struct {
	Value        string
	privateValue int
}

// 带循环引用的结构体
type CircularStruct struct {
	Name string
	Next *CircularStruct
}

func TestDeepClone_BasicTypes(t *testing.T) {
	// 测试基本类型
	tests := []any{
		42,
		3.14,
		"hello",
		true,
		'a',
	}

	for _, original := range tests {
		cloned, err := DeepCloneAny(original)
		if err != nil {
			t.Errorf("DeepCloneAny failed for %T: %v", original, err)
			continue
		}

		if !reflect.DeepEqual(original, cloned) {
			t.Errorf("Basic type clone failed for %T: expected %v, got %v", original, original, cloned)
		}
	}
}

func TestDeepClone_Slice(t *testing.T) {
	original := []int{1, 2, 3, 4, 5}
	cloned, err := DeepClone(original)
	if err != nil {
		t.Fatalf("DeepClone failed: %v", err)
	}

	// 验证值相等
	if !reflect.DeepEqual(original, cloned) {
		t.Errorf("Slice clone failed: expected %v, got %v", original, cloned)
	}

	// 验证是不同的对象
	if &original[0] == &cloned[0] {
		t.Error("Slice clone should create new slice, not reference")
	}

	// 修改原始切片，确保克隆不受影响
	original[0] = 999
	if cloned[0] == 999 {
		t.Error("Cloned slice should be independent of original")
	}
}

func TestDeepClone_Map(t *testing.T) {
	original := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	cloned, err := DeepClone(original)
	if err != nil {
		t.Fatalf("DeepClone failed: %v", err)
	}

	// 验证值相等
	if !reflect.DeepEqual(original, cloned) {
		t.Errorf("Map clone failed: expected %v, got %v", original, cloned)
	}

	// 修改原始map，确保克隆不受影响
	original["four"] = 4
	if _, exists := cloned["four"]; exists {
		t.Error("Cloned map should be independent of original")
	}
}

func TestDeepClone_Struct(t *testing.T) {
	original := TestStruct{
		PublicField:  "public",
		privateField: 42,
		PublicSlice:  []int{1, 2, 3},
		PublicNested: &NestedStruct{
			Value:        "nested public",
			privateValue: 456,
		},
	}

	cloned, err := DeepClone(original)
	if err != nil {
		t.Fatalf("DeepClone failed: %v", err)
	}

	// 验证公共字段
	if cloned.PublicField != original.PublicField {
		t.Errorf("Public field not cloned correctly: expected %s, got %s", original.PublicField, cloned.PublicField)
	}

	if !reflect.DeepEqual(cloned.PublicSlice, original.PublicSlice) {
		t.Errorf("Public slice not cloned correctly")
	}

	// 验证私有字段（使用反射访问）
	clonedVal := reflect.ValueOf(&cloned).Elem()

	privateField := clonedVal.FieldByName("privateField")
	if privateField.IsValid() && privateField.CanInterface() {
		if privateField.Int() != 42 {
			t.Errorf("Private field not cloned correctly: expected 42, got %d", privateField.Int())
		}
	} else {
		// 使用unsafe访问私有字段
		field, found := clonedVal.Type().FieldByName("privateField")
		if found && clonedVal.CanAddr() {
			privateFieldPtr := unsafe.Pointer(clonedVal.UnsafeAddr() + field.Offset)
			privateFieldValue := *(*int)(privateFieldPtr)
			if privateFieldValue != 42 {
				t.Errorf("Private field not cloned correctly: expected 42, got %d", privateFieldValue)
			}
		}
	}

	// 验证嵌套结构体
	if cloned.PublicNested == nil || cloned.PublicNested.Value != "nested public" {
		t.Error("Nested struct not cloned correctly")
	}

	// 验证独立性
	original.PublicSlice[0] = 999
	if cloned.PublicSlice[0] == 999 {
		t.Error("Cloned struct should be independent of original")
	}
}

func TestDeepClone_Pointer(t *testing.T) {
	value := 42
	original := &value

	cloned, err := DeepClone(original)
	if err != nil {
		t.Fatalf("DeepClone failed: %v", err)
	}

	// 验证值相等
	if *cloned != *original {
		t.Errorf("Pointer clone failed: expected %d, got %d", *original, *cloned)
	}

	// 验证是不同的指针
	if cloned == original {
		t.Error("Pointer clone should create new pointer, not reference")
	}

	// 修改原始值，确保克隆不受影响
	*original = 999
	if *cloned == 999 {
		t.Error("Cloned pointer should be independent of original")
	}
}

func TestDeepClone_CircularReference(t *testing.T) {
	// 简化循环引用测试：只测试简单的自引用
	node := &CircularStruct{Name: "self"}
	node.Next = node // 自引用

	// 深拷贝应该能处理循环引用而不会栈溢出
	cloned, err := DeepClone(node)
	if err != nil {
		t.Fatalf("DeepClone with circular reference failed: %v", err)
	}

	// 验证结构正确
	if cloned.Name != "self" {
		t.Errorf("Circular reference clone failed: expected name 'self', got '%s'", cloned.Name)
	}

	// 验证自引用
	if cloned.Next != cloned {
		t.Error("Circular reference clone failed: self reference not maintained")
	}

	// 验证独立性
	if cloned == node {
		t.Error("Circular reference clone should create new object")
	}
}

func TestDeepClone_Interface(t *testing.T) {
	var original interface{} = &TestStruct{
		PublicField:  "interface test",
		privateField: 100,
	}

	cloned, err := DeepCloneAny(original)
	if err != nil {
		t.Fatalf("DeepClone interface failed: %v", err)
	}

	// 验证类型
	clonedStruct, ok := cloned.(*TestStruct)
	if !ok {
		t.Fatalf("Interface clone failed: expected *TestStruct, got %T", cloned)
	}

	if clonedStruct.PublicField != "interface test" {
		t.Errorf("Interface clone failed: expected 'interface test', got '%s'", clonedStruct.PublicField)
	}
}

func TestDeepClone_NilValues(t *testing.T) {
	// 测试nil指针
	var nilPtr *int
	clonedPtr, err := DeepClone(nilPtr)
	if err != nil {
		t.Fatalf("DeepClone nil pointer failed: %v", err)
	}
	if clonedPtr != nil {
		t.Error("Nil pointer should remain nil after clone")
	}

	// 测试nil切片
	var nilSlice []int
	clonedSlice, err := DeepClone(nilSlice)
	if err != nil {
		t.Fatalf("DeepClone nil slice failed: %v", err)
	}
	if clonedSlice != nil {
		t.Error("Nil slice should remain nil after clone")
	}

	// 测试nil map
	var nilMap map[string]int
	clonedMap, err := DeepClone(nilMap)
	if err != nil {
		t.Fatalf("DeepClone nil map failed: %v", err)
	}
	if clonedMap != nil {
		t.Error("Nil map should remain nil after clone")
	}
}

func TestMustDeepClone(t *testing.T) {
	original := []int{1, 2, 3}
	cloned := MustDeepClone(original)

	if !reflect.DeepEqual(original, cloned) {
		t.Errorf("MustDeepClone failed: expected %v, got %v", original, cloned)
	}

	// 验证独立性
	if &original[0] == &cloned[0] {
		t.Error("MustDeepClone should create new slice")
	}
}

// 性能测试
func BenchmarkDeepClone_SimpleStruct(b *testing.B) {
	original := TestStruct{
		PublicField:  "benchmark",
		privateField: 123,
		PublicSlice:  []int{1, 2, 3, 4, 5},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DeepClone(original)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkDeepClone_LargeSlice(b *testing.B) {
	original := make([]int, 1000)
	for i := range original {
		original[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DeepClone(original)
		if err != nil {
			b.Fatal(err)
		}
	}
}
