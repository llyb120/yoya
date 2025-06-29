package supx

import (
    "runtime"
    "testing"
    "time"
)

type SmartTestKey struct {
    id   int
    data string
}

func TestSmartWeakMapBasicOperations(t *testing.T) {
    wm := NewSmartWeakMap[SmartTestKey, string](10, 5*time.Second)
    
    key := &SmartTestKey{id: 1, data: "test"}
    wm.Set(key, "value1")
    
    if wm.Len() != 1 {
        t.Fatalf("expected len 1, got %d", wm.Len())
    }
    
    if val, ok := wm.Get(key); !ok || val != "value1" {
        t.Fatalf("expected value1, got %v, %v", val, ok)
    }
}

func TestSmartWeakMapCapacityControl(t *testing.T) {
    // 设置容量为 3
    wm := NewSmartWeakMap[SmartTestKey, string](3, 0) // 不设置过期时间
    
    // 保持所有键的引用，防止被 GC
    var keys []*SmartTestKey
    
    // 添加 5 个条目，应该只保留最新的 3 个
    for i := 0; i < 5; i++ {
        key := &SmartTestKey{id: i, data: "test"}
        keys = append(keys, key)
        wm.Set(key, "value")
        
        time.Sleep(1 * time.Millisecond) // 确保时间不同
    }
    
    if wm.Len() != 3 {
        t.Fatalf("expected len 3 after capacity enforcement, got %d", wm.Len())
    }
    
    // 前两个键应该被清理了
    if _, ok := wm.Get(keys[0]); ok {
        t.Fatal("expected first key to be evicted")
    }
    if _, ok := wm.Get(keys[1]); ok {
        t.Fatal("expected second key to be evicted")
    }
    
    // 后三个键应该还在
    for i := 2; i < 5; i++ {
        if _, ok := wm.Get(keys[i]); !ok {
            t.Fatalf("expected key %d to still exist", i)
        }
    }
    
    runtime.KeepAlive(keys)
}

func TestSmartWeakMapLRUBehavior(t *testing.T) {
    wm := NewSmartWeakMap[SmartTestKey, string](3, 0)
    
    var keys []*SmartTestKey
    for i := 0; i < 3; i++ {
        key := &SmartTestKey{id: i, data: "test"}
        keys = append(keys, key)
        wm.Set(key, "value")
        time.Sleep(1 * time.Millisecond)
    }
    
    // 访问第一个键，使其成为最近使用的
    wm.Get(keys[0])
    time.Sleep(1 * time.Millisecond)
    
    // 添加新键，应该淘汰第二个键（最久未使用的）
    newKey := &SmartTestKey{id: 99, data: "new"}
    wm.Set(newKey, "new value")
    
    if wm.Len() != 3 {
        t.Fatalf("expected len 3, got %d", wm.Len())
    }
    
    // 第一个和第三个键应该还在
    if _, ok := wm.Get(keys[0]); !ok {
        t.Fatal("expected first key to still exist (was recently accessed)")
    }
    if _, ok := wm.Get(keys[2]); !ok {
        t.Fatal("expected third key to still exist")
    }
    
    // 第二个键应该被淘汰了
    if _, ok := wm.Get(keys[1]); ok {
        t.Fatal("expected second key to be evicted (least recently used)")
    }
    
    // 新键应该存在
    if _, ok := wm.Get(newKey); !ok {
        t.Fatal("expected new key to exist")
    }
    
    runtime.KeepAlive(keys)
    runtime.KeepAlive(newKey)
}

func TestSmartWeakMapAutoCleanup(t *testing.T) {
    wm := NewSmartWeakMap[SmartTestKey, string](10, 0) // 不设置过期时间，只测试 GC 清理
    
    // 使用函数作用域确保键可以被 GC
    func() {
        key := &SmartTestKey{id: 1, data: "test"}
        wm.Set(key, "should be cleaned")
        
        if wm.Len() != 1 {
            t.Fatalf("expected len 1 after set, got %d", wm.Len())
        }
    }()
    
    // 强制 GC，等待 finalizer 执行
    for i := 0; i < 10; i++ {
        runtime.GC()
        time.Sleep(10 * time.Millisecond)
        
        if wm.Len() == 0 {
            t.Logf("SmartWeakMap cleaned up after %d GC cycles", i+1)
            return
        }
    }
    
    t.Fatalf("expected SmartWeakMap to be cleaned up, but still has %d items", wm.Len())
}

func TestSmartWeakMapExpireCleanup(t *testing.T) {
    // 设置很短的过期时间
    wm := NewSmartWeakMap[SmartTestKey, string](10, 50*time.Millisecond)
    
    var keepAlive []*SmartTestKey
    
    // 添加一些条目
    for i := 0; i < 3; i++ {
        key := &SmartTestKey{id: i, data: "test"}
        keepAlive = append(keepAlive, key)
        wm.Set(key, "value")
    }
    
    if wm.Len() != 3 {
        t.Fatalf("expected len 3 initially, got %d", wm.Len())
    }
    
    // 等待条目过期
    time.Sleep(100 * time.Millisecond)
    
    // 触发一个 finalizer 来启动过期清理
    // 通过添加一个临时键然后让它被 GC 来触发
    func() {
        triggerKey := &SmartTestKey{id: 999, data: "trigger"}
        wm.Set(triggerKey, "trigger")
        // triggerKey 离开作用域
    }()
    
    // 强制 GC 来触发 finalizer
    for i := 0; i < 5; i++ {
        runtime.GC()
        time.Sleep(20 * time.Millisecond)
        
        currentLen := wm.Len()
        if currentLen < 3 {
            t.Logf("Expire cleanup triggered, len reduced to %d", currentLen)
            break
        }
    }
    
    runtime.KeepAlive(keepAlive)
}

func TestSmartWeakMapUpdate(t *testing.T) {
    wm := NewSmartWeakMap[SmartTestKey, string](10, 5*time.Second)
    
    key := &SmartTestKey{id: 1, data: "test"}
    
    // 设置初始值
    wm.Set(key, "value1")
    if val, ok := wm.Get(key); !ok || val != "value1" {
        t.Fatalf("expected value1, got %v, %v", val, ok)
    }
    
    // 更新值
    wm.Set(key, "value2")
    if val, ok := wm.Get(key); !ok || val != "value2" {
        t.Fatalf("expected value2 after update, got %v, %v", val, ok)
    }
    
    // 长度应该仍然是 1
    if wm.Len() != 1 {
        t.Fatalf("expected len 1 after update, got %d", wm.Len())
    }
    
    runtime.KeepAlive(key)
}

func TestSmartWeakMapDelete(t *testing.T) {
    wm := NewSmartWeakMap[SmartTestKey, string](10, 5*time.Second)
    
    key := &SmartTestKey{id: 1, data: "test"}
    wm.Set(key, "value1")
    
    if wm.Len() != 1 {
        t.Fatalf("expected len 1, got %d", wm.Len())
    }
    
    // 手动删除
    if !wm.Delete(key) {
        t.Fatal("expected Delete to return true")
    }
    
    if wm.Len() != 0 {
        t.Fatalf("expected len 0 after delete, got %d", wm.Len())
    }
    
    // 再次删除应该返回 false
    if wm.Delete(key) {
        t.Fatal("expected Delete to return false for non-existent key")
    }
}

func TestSmartWeakMapConcurrency(t *testing.T) {
    wm := NewSmartWeakMap[SmartTestKey, string](50, 100*time.Millisecond)
    
    done := make(chan bool, 3)
    
    // 并发写入
    go func() {
        for i := 0; i < 50; i++ {
            key := &SmartTestKey{id: i, data: "writer"}
            wm.Set(key, "value")
            time.Sleep(1 * time.Millisecond)
        }
        done <- true
    }()
    
    // 并发读取
    go func() {
        for i := 0; i < 50; i++ {
            key := &SmartTestKey{id: i % 10, data: "reader"}
            wm.Get(key)
            time.Sleep(1 * time.Millisecond)
        }
        done <- true
    }()
    
    // 并发 GC
    go func() {
        for i := 0; i < 10; i++ {
            runtime.GC()
            time.Sleep(10 * time.Millisecond)
        }
        done <- true
    }()
    
    // 等待所有操作完成
    for i := 0; i < 3; i++ {
        <-done
    }
    
    finalLen := wm.Len()
    t.Logf("Final SmartWeakMap length: %d", finalLen)
    
    // 应该不会超过最大容量
    if finalLen > 50 {
        t.Fatalf("expected len <= 50, got %d", finalLen)
    }
} 