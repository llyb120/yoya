# supx - 支持工具和延迟执行

supx包提供了一系列支持工具，包括动态数据容器、延迟执行记录、智能弱引用映射和时间跳跃等实用功能。

## 核心特性

- **Data容器**: 动态数据容器，支持JSON序列化和额外字段
- **Record记录**: 延迟执行记录，支持撤销和重做操作
- **SmartWeakMap**: 智能弱引用映射，自动内存管理
- **TimeLeap**: 时间跳跃工具，支持时间控制和回调

## 主要类型和函数

### 1. Data - 动态数据容器

```go
type Data[T any] struct {
    // 私有字段
}

// 创建新的Data实例
func NewData[T any]() Data[T]

// 工具函数
func Get[T any](d map[string]any, key string) T
func DeepCloneAny(src any) (any, error)
```

**Data方法:**
- `Set(data T)`: 设置主数据
- `Get() T`: 获取主数据
- `SetExtra(key string, value any)`: 设置额外字段
- `GetExtra(key string) any`: 获取额外字段
- `MarshalJSON() ([]byte, error)`: JSON序列化
- `UnmarshalJSON(data []byte) error`: JSON反序列化

**Data特性:**
- 支持任意类型的主数据
- 支持动态添加额外字段
- 自动JSON序列化，合并主数据和额外字段
- 类型安全的数据访问

### 2. Record - 延迟执行记录

```go
type Record[T any] struct {
    // 私有字段
}

// 创建新的Record实例
func NewRecord[T any](data T) *Record[T]

// JSON编码器设置
func SetJsonEncoder(encoder JSONEncoder)
func SetJsonDecoder(decoder JSONDecoder)
```

**Record方法:**
- `Set(data T)`: 设置数据
- `Get() T`: 获取当前数据
- `Commit()`: 提交当前状态
- `Rollback()`: 回滚到上一个提交点
- `History() []T`: 获取历史记录
- `Clear()`: 清空历史记录
- `MarshalJSON() ([]byte, error)`: JSON序列化
- `UnmarshalJSON(data []byte) error`: JSON反序列化

**Record特性:**
- 支持数据版本控制
- 支持撤销和重做操作
- 自动记录数据变更历史
- 支持自定义JSON编解码器

### 3. SmartWeakMap - 智能弱引用映射

```go
type SmartWeakMap[K any, V any] struct {
    // 私有字段
}

// 创建新的SmartWeakMap实例
func NewSmartWeakMap[K any, V any](maxSize int, expireDuration time.Duration) *SmartWeakMap[K, V]
```

**SmartWeakMap方法:**
- `Set(key K, value V)`: 设置键值对
- `Get(key K) (V, bool)`: 获取值
- `Delete(key K)`: 删除键值对
- `Has(key K) bool`: 检查键是否存在
- `Size() int`: 获取当前大小
- `Clear()`: 清空所有数据
- `Keys() []K`: 获取所有键
- `Values() []V`: 获取所有值

**SmartWeakMap特性:**
- 自动内存管理，当键不再被引用时自动清理
- 支持容量限制和LRU淘汰策略
- 支持过期时间控制
- 并发安全

### 4. TimeLeap - 时间跳跃工具

```go
type TimeLeapAble interface {
    Now() time.Time
    Sleep(duration time.Duration)
    After(duration time.Duration) <-chan time.Time
    // 更多时间相关方法...
}

// 创建时间跳跃实例
func TimeLeap(opts ...any) (TimeLeapAble, func())
```

**TimeLeap特性:**
- 支持时间控制和模拟
- 支持时间跳跃和回调
- 测试友好，可以模拟时间流逝
- 支持自定义时间源

## 使用示例

### Data容器使用示例

```go
package main

import (
    "fmt"
    "encoding/json"
    "github.com/llyb120/yoya/supx"
)

type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

func main() {
    // 创建Data容器
    data := supx.NewData[User]()
    
    // 设置主数据
    user := User{Name: "Alice", Age: 30}
    data.Set(user)
    
    // 设置额外字段
    data.SetExtra("created_at", "2024-01-01")
    data.SetExtra("updated_at", "2024-01-02")
    data.SetExtra("version", 1)
    
    // 获取数据
    retrievedUser := data.Get()
    fmt.Printf("User: %+v\n", retrievedUser)
    
    // 获取额外字段
    createdAt := data.GetExtra("created_at")
    fmt.Printf("Created at: %v\n", createdAt)
    
    // JSON序列化
    jsonData, err := json.Marshal(data)
    if err != nil {
        fmt.Printf("Marshal error: %v\n", err)
        return
    }
    fmt.Printf("JSON: %s\n", jsonData)
    
    // JSON反序列化
    var newData supx.Data[User]
    err = json.Unmarshal(jsonData, &newData)
    if err != nil {
        fmt.Printf("Unmarshal error: %v\n", err)
        return
    }
    
    fmt.Printf("Unmarshaled user: %+v\n", newData.Get())
    fmt.Printf("Unmarshaled version: %v\n", newData.GetExtra("version"))
}
```

### Record使用示例

```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/supx"
)

type Document struct {
    Title   string
    Content string
    Version int
}

func main() {
    // 创建Record
    doc := Document{
        Title:   "Initial Title",
        Content: "Initial content",
        Version: 1,
    }
    
    record := supx.NewRecord(doc)
    
    // 提交初始状态
    record.Commit()
    
    // 修改数据
    doc.Title = "Updated Title"
    doc.Version = 2
    record.Set(doc)
    
    // 再次修改
    doc.Content = "Updated content"
    doc.Version = 3
    record.Set(doc)
    
    // 提交当前状态
    record.Commit()
    
    // 查看当前数据
    current := record.Get()
    fmt.Printf("Current: %+v\n", current)
    
    // 回滚到上一个提交点
    record.Rollback()
    rolledBack := record.Get()
    fmt.Printf("After rollback: %+v\n", rolledBack)
    
    // 查看历史记录
    history := record.History()
    fmt.Printf("History count: %d\n", len(history))
    for i, h := range history {
        fmt.Printf("History[%d]: %+v\n", i, h)
    }
}
```

### SmartWeakMap使用示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/llyb120/yoya/supx"
)

type CacheKey struct {
    ID string
}

type CacheValue struct {
    Data      string
    Timestamp time.Time
}

func main() {
    // 创建SmartWeakMap，最大容量100，过期时间5分钟
    cache := supx.NewSmartWeakMap[*CacheKey, CacheValue](100, 5*time.Minute)
    
    // 设置数据
    key1 := &CacheKey{ID: "user:123"}
    value1 := CacheValue{
        Data:      "User data for 123",
        Timestamp: time.Now(),
    }
    cache.Set(key1, value1)
    
    key2 := &CacheKey{ID: "user:456"}
    value2 := CacheValue{
        Data:      "User data for 456",
        Timestamp: time.Now(),
    }
    cache.Set(key2, value2)
    
    // 获取数据
    if value, exists := cache.Get(key1); exists {
        fmt.Printf("Found: %+v\n", value)
    }
    
    // 检查存在性
    fmt.Printf("Key1 exists: %v\n", cache.Has(key1))
    fmt.Printf("Cache size: %d\n", cache.Size())
    
    // 获取所有键
    keys := cache.Keys()
    fmt.Printf("All keys: %v\n", len(keys))
    
    // 当key1不再被引用时，会自动从缓存中清理
    key1 = nil
    
    // 强制GC来演示自动清理
    runtime.GC()
    time.Sleep(time.Millisecond * 100)
    
    fmt.Printf("Cache size after GC: %d\n", cache.Size())
}
```

### TimeLeap使用示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/llyb120/yoya/supx"
)

func main() {
    // 创建时间跳跃实例
    timeLeap, restore := supx.TimeLeap()
    defer restore() // 确保恢复正常时间
    
    fmt.Printf("Current time: %v\n", timeLeap.Now())
    
    // 模拟时间流逝
    fmt.Println("Sleeping for 2 seconds...")
    timeLeap.Sleep(2 * time.Second)
    
    fmt.Printf("Time after sleep: %v\n", timeLeap.Now())
    
    // 使用After方法
    fmt.Println("Waiting for 1 second...")
    <-timeLeap.After(1 * time.Second)
    
    fmt.Printf("Final time: %v\n", timeLeap.Now())
}
```

## 工具函数使用示例

```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/supx"
)

func main() {
    // 使用Get函数从map中获取类型安全的值
    data := map[string]any{
        "name":    "Alice",
        "age":     30,
        "active":  true,
        "balance": 1234.56,
    }
    
    name := supx.Get[string](data, "name")
    age := supx.Get[int](data, "age")
    active := supx.Get[bool](data, "active")
    balance := supx.Get[float64](data, "balance")
    
    fmt.Printf("Name: %s, Age: %d, Active: %v, Balance: %.2f\n", 
        name, age, active, balance)
    
    // 使用DeepCloneAny进行深拷贝
    original := map[string]interface{}{
        "user": map[string]interface{}{
            "name": "Bob",
            "age":  25,
        },
        "items": []string{"item1", "item2"},
    }
    
    cloned, err := supx.DeepCloneAny(original)
    if err != nil {
        fmt.Printf("Clone error: %v\n", err)
        return
    }
    
    fmt.Printf("Original: %+v\n", original)
    fmt.Printf("Cloned: %+v\n", cloned)
    
    // 修改克隆的数据不会影响原始数据
    clonedMap := cloned.(map[string]interface{})
    userMap := clonedMap["user"].(map[string]interface{})
    userMap["name"] = "Modified Bob"
    
    fmt.Printf("After modification:\n")
    fmt.Printf("Original: %+v\n", original)
    fmt.Printf("Cloned: %+v\n", cloned)
}
```

## 性能特性

- **内存优化**: SmartWeakMap自动管理内存，避免内存泄漏
- **类型安全**: 全泛型实现，编译时类型检查
- **JSON优化**: Data和Record支持高效的JSON序列化
- **并发安全**: SmartWeakMap支持并发访问

## 注意事项

1. **内存管理**: SmartWeakMap依赖GC进行自动清理，不要长期持有键的引用
2. **JSON序列化**: 自定义JSON编解码器时需要注意类型兼容性
3. **时间控制**: TimeLeap主要用于测试场景，生产环境需谨慎使用
4. **历史记录**: Record会保存历史数据，注意内存使用量
5. **并发访问**: 除SmartWeakMap外，其他类型需要外部同步控制

## 适用场景

- **动态数据处理**: 需要在结构化数据基础上添加动态字段
- **版本控制**: 需要撤销重做功能的数据编辑场景
- **缓存管理**: 需要自动内存管理的缓存系统
- **测试工具**: 需要时间控制的单元测试和集成测试
- **数据克隆**: 需要深度复制复杂数据结构的场景

---

## 安装
```bash
go get github.com/llyb120/yoya/supx
```

---

## TimeLeap 示例
```go
leaper, done := supx.TimeLeap(supx.Async)

for i := 0; i < 3; i++ {
    idx := i
    leaper.Leap(func(){
        fmt.Println("job", idx)
    })
}

// 等待其他协程完成后调用
done() // 输出 job0 job1 job2
```

---

## 许可协议
本项目遵循 MIT License。 