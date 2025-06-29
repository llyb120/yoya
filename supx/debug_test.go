package supx

import (
	"encoding/json"
	"fmt"
	"runtime"
	"sync/atomic"
	"testing"
	"time"
)

// 简单的 finalizer 测试
func TestFinalizerBasic(t *testing.T) {
	var finalizerCalled int64

	func() {
		obj := &SmartTestKey{id: 1, data: "test"}
		runtime.SetFinalizer(obj, func(*SmartTestKey) {
			atomic.AddInt64(&finalizerCalled, 1)
		})
		// obj 在这里离开作用域
	}()

	// 强制 GC
	for i := 0; i < 10; i++ {
		runtime.GC()
		time.Sleep(10 * time.Millisecond)

		if atomic.LoadInt64(&finalizerCalled) > 0 {
			t.Logf("Finalizer called after %d GC cycles", i+1)
			return
		}
	}

	t.Fatalf("Finalizer was not called")
}

// 测试我们的智能 WeakMap 是否有隐藏的引用
func TestSmartWeakMapHiddenReferences(t *testing.T) {
	wm := NewSmartWeakMap[SmartTestKey, string](10, 0)
	var finalizerCalled int64

	func() {
		key := &SmartTestKey{id: 1, data: "test"}

		// 先测试没有 WeakMap 时 finalizer 是否工作
		runtime.SetFinalizer(key, func(*SmartTestKey) {
			atomic.AddInt64(&finalizerCalled, 1)
		})

		// 清除 finalizer，然后设置到 WeakMap
		runtime.SetFinalizer(key, nil)
		wm.Set(key, "test value")

		// key 在这里离开作用域
	}()

	t.Logf("Initial map length: %d", wm.Len())

	// 强制 GC
	for i := 0; i < 15; i++ {
		runtime.GC()
		time.Sleep(20 * time.Millisecond)

		currentLen := wm.Len()
		t.Logf("After GC %d: map length = %d", i+1, currentLen)

		if currentLen == 0 {
			t.Logf("SmartWeakMap cleaned up after %d GC cycles", i+1)
			return
		}
	}

	t.Fatalf("SmartWeakMap was not cleaned up, still has %d items", wm.Len())
}

func TestDebugJSON(t *testing.T) {
	// 使用最简单的测试用例
	type SimpleUser struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}

	data := NewData[SimpleUser]()
	user := SimpleUser{
		ID:   1,
		Name: "测试",
	}
	data.Set(user)

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}

	fmt.Printf("JSON输出: %s\n", string(jsonBytes))

	// 查看解析后的结果
	var result map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &result); err != nil {
		t.Fatalf("解析JSON失败: %v", err)
	}

	fmt.Printf("解析结果: %+v\n", result)
	for key, value := range result {
		fmt.Printf("键: %s, 值: %v, 类型: %T\n", key, value, value)
	}
}
