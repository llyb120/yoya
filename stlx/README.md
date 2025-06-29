## stlx

`stlx`（structure-extension）包提供了丰富的高性能数据结构实现，弥补 Go 标准库在集合类型方面的不足。包含有序容器、跳表、双向映射等多种数据结构。

---

## 核心数据结构

### OrderedMap - 有序映射

```go
type OrderedMap[K comparable, V any] struct {
    // 内部实现...
}

func NewOrderedMap[K comparable, V any]() *OrderedMap[K, V]
```

**功能**：保持插入顺序的映射，支持按索引访问。

**主要方法**：
- `Set(key K, value V)` - 设置键值对
- `Get(key K) (V, bool)` - 获取值
- `GetByIndex(index int) (V, bool)` - 按索引获取值
- `Keys() []K` - 获取所有键（有序）
- `Values() []V` - 获取所有值（有序）
- `Delete(key K) bool` - 删除键值对
- `Len() int` - 获取元素数量

**示例**：
```go
om := stlx.NewOrderedMap[string, int]()
om.Set("first", 1)
om.Set("second", 2)
om.Set("third", 3)

// 保持插入顺序
keys := om.Keys() // ["first", "second", "third"]

// 按索引访问
value, ok := om.GetByIndex(1) // 2, true

// 遍历保持顺序
for i := 0; i < om.Len(); i++ {
    value, _ := om.GetByIndex(i)
    fmt.Printf("Index %d: %v\n", i, value)
}
```

---

### OrderedSet - 有序集合

```go
type OrderedSet[T comparable] struct {
    // 内部实现...
}

func NewOrderedSet[T comparable]() *OrderedSet[T]
```

**功能**：保持插入顺序的集合，自动去重。

**主要方法**：
- `Add(item T) bool` - 添加元素
- `Contains(item T) bool` - 检查是否包含
- `Remove(item T) bool` - 移除元素
- `ToSlice() []T` - 转为切片
- `Len() int` - 获取元素数量

**示例**：
```go
os := stlx.NewOrderedSet[string]()
os.Add("apple")
os.Add("banana")
os.Add("apple") // 重复添加，被忽略

fmt.Println(os.ToSlice()) // ["apple", "banana"]
fmt.Println(os.Len())     // 2
```

---

### SkipList - 跳表

```go
type SkipList[T any] struct {
    // 内部实现...
}

func NewSkipList[T any]() *SkipList[T]
```

**功能**：基于跳表的有序数据结构，提供 O(log n) 的查找、插入、删除性能。

**主要方法**：
- `Insert(value T)` - 插入元素
- `Contains(value T) bool` - 检查是否包含
- `Remove(value T) bool` - 移除元素
- `Range(start, end T) []T` - 范围查询
- `Min() (T, bool)` - 最小值
- `Max() (T, bool)` - 最大值

**示例**：
```go
sl := stlx.NewSkipList[int]()
sl.Insert(10)
sl.Insert(5)
sl.Insert(15)
sl.Insert(3)

// 范围查询
values := sl.Range(5, 12) // [5, 10]

// 最值查询
min, _ := sl.Min()
max, _ := sl.Max()
```

---

### SkipMap - 跳表映射

```go
type SkipMap[K, V any] struct {
    // 内部实现...
}

func NewSkipMap[K, V any]() *SkipMap[K, V]
```

**功能**：基于跳表的有序映射，键按顺序存储。

**主要方法**：
- `Set(key K, value V)` - 设置键值对
- `Get(key K) (V, bool)` - 获取值
- `Remove(key K) bool` - 移除键值对
- `RangeKeys(start, end K) []K` - 键范围查询
- `ForEach(fn func(K, V) bool)` - 有序遍历

**示例**：
```go
sm := stlx.NewSkipMap[int, string]()
sm.Set(10, "ten")
sm.Set(5, "five")
sm.Set(15, "fifteen")

// 按键顺序遍历
sm.ForEach(func(k int, v string) bool {
    fmt.Printf("%d: %s\n", k, v)
    return true // 继续遍历
})
// 输出: 5: five, 10: ten, 15: fifteen
```

---

### BiMap - 双向映射

```go
type BiMap[K comparable, V comparable] struct {
    // 内部实现...
}

func NewBiMap[K comparable, V comparable]() *BiMap[K, V]
```

**功能**：支持双向查找的映射，键值互相唯一。

**主要方法**：
- `Set(key K, value V)` - 设置键值对
- `Get(key K) (V, bool)` - 通过键获取值
- `GetKey(value V) (K, bool)` - 通过值获取键
- `RemoveByKey(key K) bool` - 通过键删除
- `RemoveByValue(value V) bool` - 通过值删除

**示例**：
```go
bm := stlx.NewBiMap[string, int]()
bm.Set("one", 1)
bm.Set("two", 2)

// 正向查找
value, ok := bm.Get("one") // 1, true

// 反向查找
key, ok := bm.GetKey(2) // "two", true
```

---

### MultiMap - 多值映射

```go
type MultiMap[K comparable, V any] struct {
    // 内部实现...
}

func NewMultiMap[K comparable, V any]() *MultiMap[K, V]
```

**功能**：一个键可以对应多个值的映射。

**主要方法**：
- `Set(key K, value V)` - 添加键值对
- `Get(key K) []V` - 获取键对应的所有值
- `Remove(key K, value V) bool` - 移除特定键值对
- `RemoveAll(key K) bool` - 移除键的所有值
- `Keys() []K` - 获取所有键

**示例**：
```go
mm := stlx.NewMultiMap[string, string]()
mm.Set("fruits", "apple")
mm.Set("fruits", "banana")
mm.Set("fruits", "orange")

fruits := mm.Get("fruits") // ["apple", "banana", "orange"]
```

---

## 完整使用示例

```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/stlx"
)

func main() {
    // 1. OrderedMap 示例
    fmt.Println("=== OrderedMap 示例 ===")
    orderedMapExample()
    
    // 2. SkipList 示例  
    fmt.Println("\n=== SkipList 示例 ===")
    skipListExample()
    
    // 3. BiMap 示例
    fmt.Println("\n=== BiMap 示例 ===")
    biMapExample()
    
    // 4. 实际应用场景
    fmt.Println("\n=== 实际应用场景 ===")
    realWorldExample()
}

func orderedMapExample() {
    // 创建有序映射
    config := stlx.NewOrderedMap[string, interface{}]()
    
    // 按顺序添加配置项
    config.Set("database.host", "localhost")
    config.Set("database.port", 5432)
    config.Set("database.name", "myapp")
    config.Set("server.port", 8080)
    config.Set("server.debug", true)
    
    // 保持插入顺序输出
    fmt.Println("配置项（按插入顺序）:")
    for i := 0; i < config.Len(); i++ {
        key := config.Keys()[i]
        value, _ := config.Get(key)
        fmt.Printf("  %s: %v\n", key, value)
    }
    
    // 按索引访问
    firstValue, _ := config.GetByIndex(0)
    fmt.Printf("第一个配置项的值: %v\n", firstValue)
}

func skipListExample() {
    // 创建跳表存储分数
    scores := stlx.NewSkipList[int]()
    
    // 插入分数
    testScores := []int{85, 92, 78, 96, 88, 73, 91}
    for _, score := range testScores {
        scores.Insert(score)
    }
    
    // 查找特定分数
    fmt.Printf("是否包含90分: %v\n", scores.Contains(90))
    fmt.Printf("是否包含85分: %v\n", scores.Contains(85))
    
    // 范围查询：80-90分的成绩
    goodScores := scores.Range(80, 90)
    fmt.Printf("80-90分的成绩: %v\n", goodScores)
    
    // 最值查询
    min, _ := scores.Min()
    max, _ := scores.Max()
    fmt.Printf("最低分: %d, 最高分: %d\n", min, max)
}

func biMapExample() {
    // 创建用户ID和用户名的双向映射
    userMap := stlx.NewBiMap[int, string]()
    
    userMap.Set(1001, "alice")
    userMap.Set(1002, "bob")
    userMap.Set(1003, "charlie")
    
    // 通过ID查找用户名
    if name, ok := userMap.Get(1002); ok {
        fmt.Printf("用户ID 1002 对应的用户名: %s\n", name)
    }
    
    // 通过用户名查找ID
    if id, ok := userMap.GetKey("alice"); ok {
        fmt.Printf("用户名 alice 对应的ID: %d\n", id)
    }
    
    // 删除映射
    userMap.RemoveByKey(1003)
    fmt.Printf("删除后的用户数量: %d\n", userMap.Len())
}

func realWorldExample() {
    // 模拟一个简单的任务调度系统
    
    // 使用 SkipList 管理任务优先级
    taskQueue := stlx.NewSkipList[Task]()
    
    // 使用 MultiMap 管理用户的多个任务
    userTasks := stlx.NewMultiMap[string, int]()
    
    // 使用 OrderedMap 记录任务执行历史
    taskHistory := stlx.NewOrderedMap[int, string]()
    
    // 添加任务
    tasks := []Task{
        {ID: 1, Priority: 5, Name: "发送邮件"},
        {ID: 2, Priority: 1, Name: "备份数据"},
        {ID: 3, Priority: 8, Name: "处理订单"},
        {ID: 4, Priority: 3, Name: "生成报告"},
    }
    
    for _, task := range tasks {
        taskQueue.Insert(task)
        userTasks.Set("admin", task.ID)
        taskHistory.Set(task.ID, fmt.Sprintf("任务创建: %s", task.Name))
    }
    
    // 按优先级处理任务
    fmt.Println("按优先级处理任务:")
    min, ok := taskQueue.Min()
    for ok {
        fmt.Printf("  执行任务: %s (优先级: %d)\n", min.Name, min.Priority)
        taskQueue.Remove(min)
        
        // 更新历史
        taskHistory.Set(min.ID, fmt.Sprintf("任务完成: %s", min.Name))
        
        min, ok = taskQueue.Min()
    }
    
    // 查看用户的所有任务
    adminTasks := userTasks.Get("admin")
    fmt.Printf("\n管理员的任务ID: %v\n", adminTasks)
    
    // 查看任务历史（按时间顺序）
    fmt.Println("\n任务历史记录:")
    for i := 0; i < taskHistory.Len(); i++ {
        taskID := taskHistory.Keys()[i]
        history, _ := taskHistory.Get(taskID)
        fmt.Printf("  任务%d: %s\n", taskID, history)
    }
}

type Task struct {
    ID       int
    Priority int
    Name     string
}

// 实现比较接口，用于 SkipList 排序
func (t Task) Less(other Task) bool {
    return t.Priority < other.Priority
}
```

---

## 性能特性

### 时间复杂度对比

| 数据结构 | 插入 | 查找 | 删除 | 空间复杂度 |
|----------|------|------|------|------------|
| OrderedMap | O(1) | O(1) | O(1) | O(n) |
| OrderedSet | O(1) | O(1) | O(1) | O(n) |
| SkipList | O(log n) | O(log n) | O(log n) | O(n) |
| SkipMap | O(log n) | O(log n) | O(log n) | O(n) |
| BiMap | O(1) | O(1) | O(1) | O(n) |
| MultiMap | O(1) | O(1) | O(k) | O(n) |

### 内存优化
- **对象池**：频繁创建的节点使用对象池减少 GC 压力
- **紧凑存储**：OrderedMap 使用数组+哈希表的混合结构
- **延迟删除**：SkipList 使用标记删除，批量清理

---

## 选择指南

### 何时使用 OrderedMap
- 需要保持插入顺序的映射
- 需要按索引访问元素
- 配置管理、属性列表等场景

### 何时使用 SkipList
- 需要有序存储且频繁查找
- 范围查询需求
- 替代平衡二叉树的轻量级方案

### 何时使用 BiMap
- 需要双向查找的映射关系
- 用户ID与用户名映射
- 缓存键值对互查

### 何时使用 MultiMap
- 一个键对应多个值
- 分组数据存储
- 索引构建

---

## 注意事项

1. **线程安全**：所有数据结构都不是线程安全的，并发使用需要外部同步
2. **比较函数**：SkipList 和 SkipMap 需要元素实现比较接口
3. **内存占用**：有序结构会有额外的内存开销
4. **性能权衡**：根据具体使用场景选择合适的数据结构

---

## 安装
```bash
go get github.com/llyb120/yoya/stlx
``` 