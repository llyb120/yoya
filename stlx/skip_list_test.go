package stlx

import (
	"fmt"
	"reflect"
	"testing"
)

func TestSkipList(t *testing.T) {
	// 创建一个整数跳表，使用默认的比较函数
	sl := NewSkipList[int](func(a, b int) bool {
		return a < b
	})

	// 测试 Set 和 Get
	sl.Add(1)
	sl.Add(2)
	sl.Add(3)

	// 测试 Len
	if sl.Len() != 3 {
		t.Errorf("Expected length 3, got %d", sl.Len())
	}

	// 测试 Get
	if !sl.Has(2) {
		t.Errorf("Expected 'two'")
	}

	// 测试不存在的键
	if sl.Has(4) {
		t.Errorf("Expected not found for key 4")
	}

	// 测试更新值
	sl.Add(2)
	if !sl.Has(2) {
		t.Errorf("Expected 'TWO'")
	}

	sl.Del(2)

	// 测试 Keys
	keys := sl.Vals()
	expectedKeys := []int{1, 2, 3}
	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("Expected keys %v, got %v", expectedKeys, keys)
	}

	// 测试 Vals
	vals := sl.Vals()
	expectedVals := []int{1, 2, 3}
	if !reflect.DeepEqual(vals, expectedVals) {
		t.Errorf("Expected values %v, got %v", expectedVals, vals)
	}

	// 测试 Del
	sl.Del(2)
	if sl.Has(2) {
		t.Errorf("Expected 'TWO'")
	}
	if sl.Len() != 2 {
		t.Errorf("Expected length 2, got %d", sl.Len())
	}
	if sl.Has(2) {
		t.Errorf("Expected not found for key 2 after deletion")
	}

	// 测试 For
	sl.For(func(key int) bool {
		t.Log(key)
		return true
	})

	// 测试 For 提前终止
	count := 0
	sl.For(func(key int) bool {
		count++
		return false // 只访问第一个元素
	})
	if count != 1 {
		t.Errorf("Expected 1 item visited, got %d", count)
	}

	// 测试 Clear
	sl.Clear()
	if sl.Len() != 0 {
		t.Errorf("Expected length 0 after Clear, got %d", sl.Len())
	}
	if len(sl.Vals()) != 0 {
		t.Errorf("Expected empty keys after Clear, got %v", sl.Vals())
	}
}

func TestSkipListWithStrings(t *testing.T) {
	// 创建一个字符串跳表
	sl := NewSkipList[string](func(a, b string) bool {
		return a < b
	})

	// 添加一些数据
	sl.Add("apple")
	sl.Add("banana")
	sl.Add("cherry")
	sl.Add("date")
	sl.Add("elderberry")

	// 测试按顺序遍历
	expected := []struct {
		key string
	}{
		{"apple"},
		{"banana"},
		{"cherry"},
		{"date"},
		{"elderberry"},
	}

	i := 0
	sl.For(func(key string) bool {
		if i >= len(expected) {
			t.Errorf("Too many items in skiplist")
			return false
		}
		if key != expected[i].key {
			t.Errorf("Expected (%s), got (%s)", expected[i].key, key)
		}
		i++
		return true
	})

	if i != len(expected) {
		t.Errorf("Expected %d items, got %d", len(expected), i)
	}

	// 测试删除中间元素
	sl.Del("cherry")

	// 验证删除后的顺序
	expected = []struct {
		key string
	}{
		{"apple"},
		{"banana"},
		{"date"},
		{"elderberry"},
	}

	i = 0
	sl.For(func(key string) bool {
		if i >= len(expected) {
			t.Errorf("Too many items in skiplist")
			return false
		}
		if key != expected[i].key {
			t.Errorf("Expected (%s), got (%s)", expected[i].key, key)
		}
		i++
		return true
	})

	if i != len(expected) {
		t.Errorf("Expected %d items, got %d", len(expected), i)
	}
}

// 测试跳表是否正确实现了 Map 接口
func TestSkipListAsMap(t *testing.T) {
	// 创建一个跳表，并将其作为 Map 接口使用
	var m = NewSkipList[string](func(a, b string) bool {
		return a < b
	})

	// 测试 Map 接口的方法
	m.Add("a")
	m.Add("b")
	m.Add("c")

	// 测试 Len
	if m.Len() != 3 {
		t.Errorf("Expected length 3, got %d", m.Len())
	}

	// 测试 Get
	if !m.Has("b") {
		t.Errorf("Expected 2")
	}

	// 测试 Keys
	keys := m.Vals()
	if len(keys) != 3 {
		t.Errorf("Expected 3 keys, got %d", len(keys))
	}
	expectedKeys := []string{"a", "b", "c"}
	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("Expected key %s at position %d, got %s", expectedKeys[i], i, k)
		}
	}

	// 测试 Del
	m.Del("b")
	if m.Has("b") {
		t.Errorf("Expected 2")
	}
	if m.Len() != 2 {
		t.Errorf("Expected length 2, got %d", m.Len())
	}

	// 测试 For
	m.For(func(key string) bool {
		t.Log(key)
		return true
	})

	// 测试 Clear
	m.Clear()
	if m.Len() != 0 {
		t.Errorf("Expected length 0 after Clear, got %d", m.Len())
	}
}

// 测试跳表的性能
func BenchmarkSkipListSet(b *testing.B) {
	sl := NewSkipList[int](func(a, b int) bool {
		return a < b
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Add(i)
	}
}

func BenchmarkSkipListGet(b *testing.B) {
	sl := NewSkipList[int](func(a, b int) bool {
		return a < b
	})

	// 预先填充数据
	for i := 0; i < 10000; i++ {
		sl.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Get(i % 10000)
	}
}

func BenchmarkSkipListDel(b *testing.B) {
	sl := NewSkipList[int](func(a, b int) bool {
		return a < b
	})

	// 预先填充数据
	for i := 0; i < b.N; i++ {
		sl.Add(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Del(i)
	}
}

// 示例：如何使用跳表
func ExampleSkipList() {
	// 创建一个整数跳表
	sl := NewSkipList[int](func(a, b int) bool {
		return a < b
	})

	// 添加一些数据
	sl.Add(3)
	sl.Add(1)
	sl.Add(2)

	// 遍历跳表（将按键的顺序输出）
	sl.For(func(key int) bool {
		fmt.Printf("%d\n", key)
		return true
	})

	// Output:
	// 1
	// 2
	// 3
}
