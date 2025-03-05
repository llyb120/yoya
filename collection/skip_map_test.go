package collection

import (
	"testing"
)

func TestSkipMap(t *testing.T) {
	// 创建一个整数跳表，使用默认的比较函数
	sl := NewSkipMap[int, string](func(a, b int) bool {
		return a < b
	})

	// 测试 Set 和 Get
	sl.Set(1, "one")
	sl.Set(2, "two")
	sl.Set(3, "three")

	// 测试 Len
	if sl.Len() != 3 {
		t.Errorf("Expected length 3, got %d", sl.Len())
	}

	// 测试 Get
	if val, ok := sl.Get(2); !ok || val != "two" {
		t.Errorf("Expected 'two', got '%s', ok: %v", val, ok)
	}

	// 测试不存在的键
	if _, ok := sl.Get(4); ok {
		t.Errorf("Expected not found for key 4")
	}

	// 测试更新值
	sl.Set(2, "TWO")
	if val, ok := sl.Get(2); !ok || val != "TWO" {
		t.Errorf("Expected 'TWO', got '%s', ok: %v", val, ok)
	}

	// 测试 Keys
	keys := sl.Keys()
	expectedKeys := []int{1, 2, 3}
	if len(keys) != len(expectedKeys) {
		t.Errorf("Expected keys %v, got %v", expectedKeys, keys)
	}
	for i, k := range keys {
		if k != expectedKeys[i] {
			t.Errorf("Expected key %d at position %d, got %d", expectedKeys[i], i, k)
		}
	}

	// 测试 Vals
	vals := sl.Vals()
	expectedVals := []string{"one", "TWO", "three"}
	if len(vals) != len(expectedVals) {
		t.Errorf("Expected values %v, got %v", expectedVals, vals)
	}
	for i, v := range vals {
		if v != expectedVals[i] {
			t.Errorf("Expected value %s at position %d, got %s", expectedVals[i], i, v)
		}
	}

	// 测试 Del
	val := sl.Del(2)
	if val != "TWO" {
		t.Errorf("Expected 'TWO', got '%s'", val)
	}
	if sl.Len() != 2 {
		t.Errorf("Expected length 2, got %d", sl.Len())
	}
	if _, ok := sl.Get(2); ok {
		t.Errorf("Expected not found for key 2 after deletion")
	}

	// 测试 Each
	visited := make(map[int]string)
	sl.Each(func(key int, value string) bool {
		visited[key] = value
		return true
	})
	if len(visited) != 2 {
		t.Errorf("Expected 2 items visited, got %d", len(visited))
	}
	if v, ok := visited[1]; !ok || v != "one" {
		t.Errorf("Expected 'one' for key 1, got '%s', ok: %v", v, ok)
	}
	if v, ok := visited[3]; !ok || v != "three" {
		t.Errorf("Expected 'three' for key 3, got '%s', ok: %v", v, ok)
	}

	// 测试 Each 提前终止
	count := 0
	sl.Each(func(key int, value string) bool {
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
	if len(sl.Keys()) != 0 {
		t.Errorf("Expected empty keys after Clear, got %v", sl.Keys())
	}
}

func TestSkipMapWithStrings(t *testing.T) {
	// 创建一个字符串跳表
	sl := NewSkipMap[string, int](func(a, b string) bool {
		return a < b
	})

	// 添加一些数据
	sl.Set("apple", 1)
	sl.Set("banana", 2)
	sl.Set("cherry", 3)
	sl.Set("date", 4)
	sl.Set("elderberry", 5)

	// 测试按顺序遍历
	expected := []struct {
		key   string
		value int
	}{
		{"apple", 1},
		{"banana", 2},
		{"cherry", 3},
		{"date", 4},
		{"elderberry", 5},
	}

	i := 0
	sl.Each(func(key string, value int) bool {
		if i >= len(expected) {
			t.Errorf("Too many items in skiplist")
			return false
		}
		if key != expected[i].key || value != expected[i].value {
			t.Errorf("Expected (%s, %d), got (%s, %d)", expected[i].key, expected[i].value, key, value)
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
		key   string
		value int
	}{
		{"apple", 1},
		{"banana", 2},
		{"date", 4},
		{"elderberry", 5},
	}

	i = 0
	sl.Each(func(key string, value int) bool {
		if i >= len(expected) {
			t.Errorf("Too many items in skiplist")
			return false
		}
		if key != expected[i].key || value != expected[i].value {
			t.Errorf("Expected (%s, %d), got (%s, %d)", expected[i].key, expected[i].value, key, value)
		}
		i++
		return true
	})

	if i != len(expected) {
		t.Errorf("Expected %d items, got %d", len(expected), i)
	}
}

// 测试跳表是否正确实现了 Map 接口
func TestSkipMapAsMap(t *testing.T) {
	// 创建一个跳表，并将其作为 Map 接口使用
	var m Map[string, int] = NewSkipMap[string, int](func(a, b string) bool {
		return a < b
	})

	// 测试 Map 接口的方法
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	// 测试 Len
	if m.Len() != 3 {
		t.Errorf("Expected length 3, got %d", m.Len())
	}

	// 测试 Get
	if val, ok := m.Get("b"); !ok || val != 2 {
		t.Errorf("Expected 2, got %d, ok: %v", val, ok)
	}

	// 测试 Keys
	keys := m.Keys()
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
	val := m.Del("b")
	if val != 2 {
		t.Errorf("Expected 2, got %d", val)
	}
	if m.Len() != 2 {
		t.Errorf("Expected length 2, got %d", m.Len())
	}

	// 测试 Each
	visited := make(map[string]int)
	m.Each(func(key string, value int) bool {
		visited[key] = value
		return true
	})
	if len(visited) != 2 {
		t.Errorf("Expected 2 items visited, got %d", len(visited))
	}
	if v, ok := visited["a"]; !ok || v != 1 {
		t.Errorf("Expected 1 for key 'a', got %d, ok: %v", v, ok)
	}
	if v, ok := visited["c"]; !ok || v != 3 {
		t.Errorf("Expected 3 for key 'c', got %d, ok: %v", v, ok)
	}

	// 测试 Clear
	m.Clear()
	if m.Len() != 0 {
		t.Errorf("Expected length 0 after Clear, got %d", m.Len())
	}
}

// 测试跳表的性能
func BenchmarkSkipMapSet(b *testing.B) {
	sl := NewSkipMap[int, int](func(a, b int) bool {
		return a < b
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Set(i, i)
	}
}

func BenchmarkSkipMapGet(b *testing.B) {
	sl := NewSkipMap[int, int](func(a, b int) bool {
		return a < b
	})

	// 预先填充数据
	for i := 0; i < 10000; i++ {
		sl.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Get(i % 10000)
	}
}

func BenchmarkSkipMapDel(b *testing.B) {
	sl := NewSkipMap[int, int](func(a, b int) bool {
		return a < b
	})

	// 预先填充数据
	for i := 0; i < b.N; i++ {
		sl.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Del(i)
	}
}
