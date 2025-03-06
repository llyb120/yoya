package stlx

import (
	"runtime"
	"testing"
	"time"
)

func TestWeakMap(t *testing.T) {
	wm := NewWeakMap[string, string]()
	wm.Set("key", "value")
	if val, ok := wm.Get("key"); !ok || val != "value" {
		t.Errorf("Expected value to be 'value'")
	}
}

func TestWeakMapGC(t *testing.T) {
	wm := NewWeakMap[*string, string]()

	// 创建一个作用域，让 key 在作用域结束后失去引用
	func() {
		key := new(string)
		*key = "test"
		wm.Set(key, "value")

		// 验证值是否正确存储
		if val, ok := wm.Get(key); !ok || val != "value" {
			t.Errorf("Expected value to be 'value', got %v, exists: %v", val, ok)
		}
	}()

	// 手动触发垃圾回收多次，确保弱引用被清理
	for i := 0; i < 10; i++ {
		runtime.GC()
		// 短暂等待让 GC 有机会工作
		time.Sleep(time.Millisecond * 10)
	}

	// 等待更长时间，确保终结器有机会运行
	time.Sleep(time.Millisecond * 500)

	// 再次触发 GC
	runtime.GC()
	time.Sleep(time.Millisecond * 500)

	// 验证 map 是否为空
	if size := wm.Len(); size != 0 {
		t.Errorf("Expected WeakMap to be empty after GC, but got size: %d", size)
	}
}
