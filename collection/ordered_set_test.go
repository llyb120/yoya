package collection

import (
	"encoding/json"
	"testing"
)

func TestOrderedSet(t *testing.T) {
	set := NewSet[string]()

	// 测试添加元素
	set.Add("a")
	set.Add("b")
	set.Add("c")
	set.Add("a") // 重复添加

	if set.Len() != 3 {
		t.Errorf("Expected size 3, got %d", set.Len())
	}

	// 测试包含元素
	if !set.Has("b") {
		t.Error("Expected set to contain 'b'")
	}

	// 测试移除元素
	set.Del("b")
	if set.Has("b") {
		t.Error("Expected set to not contain 'b' after removal")
	}

	// 测试序列化
	jsonData, err := json.Marshal(set)
	if err != nil {
		t.Errorf("MarshalJSON failed: %v", err)
	}
	t.Logf("JSON data: %s", jsonData)

	// 测试反序列化
	set2 := NewSet[string]()
	err = json.Unmarshal(jsonData, set2)
	if err != nil {
		t.Errorf("UnmarshalJSON failed: %v", err)
	}

	// 验证元素顺序
	elements := set2.Vals()
	expected := []string{"a", "c"}
	if len(elements) != len(expected) {
		t.Errorf("Expected %v elements, got %v", expected, elements)
	}
	for i, e := range expected {
		if elements[i] != e {
			t.Errorf("Expected element %d to be %s, got %s", i, e, elements[i])
		}
	}

	// 测试遍历
	var result []string
	set2.Each(func(element string) bool {
		result = append(result, element)
		return true
	})

	if len(result) != len(expected) {
		t.Errorf("Each: Expected %v elements, got %v", expected, result)
	}
}
