# lsx - 函数式编程和列表处理

lsx包提供了丰富的函数式编程工具，包括映射、过滤、归约等操作，支持异步处理和并发执行，让列表处理变得更加优雅和高效。

## 核心特性

- **函数式编程**: 提供Map、Filter、Reduce等经典函数式操作
- **异步支持**: 支持异步执行和并发处理
- **类型安全**: 全泛型实现，编译时类型检查
- **选项配置**: 支持多种处理选项，如忽略nil、去重等
- **高性能**: 优化的算法实现，支持并发执行

## 选项常量

```go
type lsxOption int

const (
    IgnoreNil   lsxOption = iota  // 忽略nil值
    IgnoreEmpty                   // 忽略空值
    Async                         // 异步执行
    DoDistinct                    // 执行去重
)
```

## 主要函数

### 1. 核心转换函数

#### Map - 映射转换
```go
func Map[T any, R any](arr []T, fn func(T, int) R, opts ...lsxOption) []R
```
将数组中的每个元素通过函数转换为新类型的元素。

**参数:**
- `arr []T`: 源数组
- `fn func(T, int) R`: 转换函数，接收元素和索引，返回新元素
- `opts ...lsxOption`: 可选配置

**支持的选项:**
- `IgnoreNil`: 过滤掉转换结果中的nil值
- `IgnoreEmpty`: 过滤掉转换结果中的空值
- `Async`: 异步并发执行转换
- `DoDistinct`: 对结果进行去重

#### FlatMap - 扁平化映射
```go
func FlatMap[T any, R any](arr []T, fn func(T, int) []R) []R
```
将数组中的每个元素转换为数组，然后将所有结果数组合并为一个数组。

### 2. 过滤和查找函数

#### Filter - 过滤
```go
func Filter[T any](arr *[]T, fn func(T, int) bool)
```
根据条件函数过滤数组元素，直接修改原数组。

#### Find - 查找元素
```go
func Find[T any](arr []T, fn func(T, int) bool) (T, bool)
```
查找第一个满足条件的元素。

#### Pos - 查找位置
```go
func Pos[T any](arr []T, fn func(T, int) bool) int
```
查找第一个满足条件的元素的位置，未找到返回-1。

#### Has - 检查存在
```go
func Has[T comparable](arr []T, target T) bool
```
检查数组中是否包含指定元素。

### 3. 数组操作函数

#### Del - 删除元素
```go
func Del[T comparable](arr *[]T, pos int)
```
删除指定位置的元素。

#### Distinct - 去重
```go
func Distinct[T any](arr *[]T, fn ...func(T, int) any)
```
对数组进行去重，可选择自定义键函数。

#### Sort - 排序
```go
func Sort[T any](arr *[]T, less func(T, T) bool)
```
使用自定义比较函数对数组进行排序。

#### Chunk - 分块
```go
func Chunk[T any](arr []T, size int) [][]T
```
将数组分割成指定大小的块。

### 4. 归约和聚合函数

#### Reduce - 归约
```go
func Reduce[T any, R any](arr []T, fn func(R, T) R, initial R) R
```
将数组归约为单个值。

#### Group - 分组
```go
func Group[T any](arr []T, fn func(T, int) any) [][]T
```
根据键函数将数组元素分组。

#### GroupMap - 分组映射
```go
func GroupMap[K comparable, V any](arr []V, fn func(V, int) K) map[K][]V
```
根据键函数将数组元素分组为映射。

#### ToMap - 转换为映射
```go
func ToMap[K comparable, V any](arr []V, fn func(V, int) K) map[K]V
```
将数组转换为映射，键由函数生成。

### 5. 遍历函数

#### For - 遍历
```go
func For[T any](arr []T, fn func(T, int) bool)
```
遍历数组，当函数返回false时停止。

### 6. 映射操作函数

#### Keys - 获取键
```go
func Keys[K comparable, V any](mp map[K]V) []K
```
获取映射的所有键。

#### Vals - 获取值
```go
func Vals[K comparable, V any](mp map[K]V) []V
```
获取映射的所有值。

### 7. 高级函数

#### Mock - 模拟操作
```go
func Mock[K any, T any](arr *[]K, fn func(*[]T)) error
```
在临时数组上执行操作，用于测试和模拟。

## 使用示例

### 基础映射和过滤

```go
package main

import (
    "fmt"
    "strings"
    "github.com/llyb120/yoya/lsx"
)

func main() {
    // 基本映射
    numbers := []int{1, 2, 3, 4, 5}
    doubled := lsx.Map(numbers, func(n, i int) int {
        return n * 2
    })
    fmt.Printf("Doubled: %v\n", doubled) // [2, 4, 6, 8, 10]
    
    // 类型转换映射
    strings := lsx.Map(numbers, func(n, i int) string {
        return fmt.Sprintf("num_%d", n)
    })
    fmt.Printf("Strings: %v\n", strings)
    
    // 过滤偶数
    lsx.Filter(&numbers, func(n, i int) bool {
        return n%2 == 0
    })
    fmt.Printf("Even numbers: %v\n", numbers)
    
    // 查找元素
    words := []string{"apple", "banana", "cherry", "date"}
    found, exists := lsx.Find(words, func(s string, i int) bool {
        return strings.HasPrefix(s, "c")
    })
    if exists {
        fmt.Printf("Found word starting with 'c': %s\n", found)
    }
}
```

### 异步处理示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/llyb120/yoya/lsx"
)

func expensiveOperation(n int) int {
    time.Sleep(time.Millisecond * 100) // 模拟耗时操作
    return n * n
}

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    
    // 同步处理
    start := time.Now()
    syncResult := lsx.Map(numbers, func(n, i int) int {
        return expensiveOperation(n)
    })
    syncDuration := time.Since(start)
    
    // 异步处理
    start = time.Now()
    asyncResult := lsx.Map(numbers, func(n, i int) int {
        return expensiveOperation(n)
    }, lsx.Async)
    asyncDuration := time.Since(start)
    
    fmt.Printf("Sync result: %v (took %v)\n", syncResult, syncDuration)
    fmt.Printf("Async result: %v (took %v)\n", asyncResult, asyncDuration)
    fmt.Printf("Speedup: %.2fx\n", float64(syncDuration)/float64(asyncDuration))
}
```

### 复杂数据处理示例

```go
package main

import (
    "fmt"
    "strings"
    "github.com/llyb120/yoya/lsx"
)

type Product struct {
    ID       int
    Name     string
    Price    float64
    Category string
    InStock  bool
}

func main() {
    products := []Product{
        {1, "iPhone 15", 999.99, "Electronics", true},
        {2, "MacBook Pro", 1999.99, "Electronics", true},
        {3, "Coffee Mug", 15.99, "Home", true},
        {4, "Wireless Mouse", 29.99, "Electronics", false},
        {5, "Desk Lamp", 45.99, "Home", true},
        {6, "Keyboard", 79.99, "Electronics", true},
    }
    
    // 1. 筛选有库存的电子产品
    electronics := lsx.Map(products, func(p Product, i int) Product {
        return p
    })
    lsx.Filter(&electronics, func(p Product, i int) bool {
        return p.Category == "Electronics" && p.InStock
    })
    
    fmt.Println("Available Electronics:")
    lsx.For(electronics, func(p Product, i int) bool {
        fmt.Printf("- %s: $%.2f\n", p.Name, p.Price)
        return true
    })
    
    // 2. 按类别分组
    categoryGroups := lsx.GroupMap(products, func(p Product, i int) string {
        return p.Category
    })
    
    fmt.Println("\nProducts by Category:")
    for category, items := range categoryGroups {
        fmt.Printf("%s (%d items):\n", category, len(items))
        for _, item := range items {
            fmt.Printf("  - %s\n", item.Name)
        }
    }
    
    // 3. 计算总价值
    totalValue := lsx.Reduce(products, func(sum float64, p Product) float64 {
        if p.InStock {
            return sum + p.Price
        }
        return sum
    }, 0.0)
    
    fmt.Printf("\nTotal inventory value: $%.2f\n", totalValue)
    
    // 4. 创建产品名称映射
    nameMap := lsx.ToMap(products, func(p Product, i int) int {
        return p.ID
    })
    
    fmt.Println("\nProduct ID to Name mapping:")
    for id, product := range nameMap {
        fmt.Printf("ID %d: %s\n", id, product.Name)
    }
    
    // 5. 扁平化处理 - 获取所有标签
    productTags := lsx.FlatMap(products, func(p Product, i int) []string {
        tags := []string{p.Category}
        if p.InStock {
            tags = append(tags, "Available")
        } else {
            tags = append(tags, "Out of Stock")
        }
        if p.Price > 100 {
            tags = append(tags, "Premium")
        }
        return tags
    })
    
    // 去重标签
    lsx.Distinct(&productTags)
    fmt.Printf("\nAll product tags: %v\n", productTags)
}
```

### 高级选项使用示例

```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/lsx"
)

func main() {
    // 测试数据，包含各种边界情况
    data := []interface{}{
        "hello",
        "",
        nil,
        42,
        0,
        []int{1, 2, 3},
        []int{},
        map[string]int{"a": 1},
        map[string]int{},
    }
    
    // 1. 转换为字符串，不过滤
    strings1 := lsx.Map(data, func(v interface{}, i int) string {
        if v == nil {
            return ""
        }
        return fmt.Sprintf("%v", v)
    })
    fmt.Printf("All strings: %v\n", strings1)
    
    // 2. 转换为字符串，忽略nil
    strings2 := lsx.Map(data, func(v interface{}, i int) string {
        if v == nil {
            return ""
        }
        return fmt.Sprintf("%v", v)
    }, lsx.IgnoreNil)
    fmt.Printf("Ignore nil: %v\n", strings2)
    
    // 3. 转换为字符串，忽略空值
    strings3 := lsx.Map(data, func(v interface{}, i int) string {
        if v == nil {
            return ""
        }
        return fmt.Sprintf("%v", v)
    }, lsx.IgnoreEmpty)
    fmt.Printf("Ignore empty: %v\n", strings3)
    
    // 4. 组合选项：忽略nil和空值，并去重
    strings4 := lsx.Map(data, func(v interface{}, i int) string {
        if v == nil {
            return ""
        }
        return fmt.Sprintf("%v", v)
    }, lsx.IgnoreNil, lsx.IgnoreEmpty, lsx.DoDistinct)
    fmt.Printf("Ignore nil & empty, distinct: %v\n", strings4)
    
    // 5. 数组分块
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    chunks := lsx.Chunk(numbers, 3)
    fmt.Printf("Chunks of 3: %v\n", chunks)
    
    // 6. 自定义分组
    words := []string{"apple", "apricot", "banana", "blueberry", "cherry", "coconut"}
    grouped := lsx.Group(words, func(word string, i int) any {
        return word[0] // 按首字母分组
    })
    fmt.Printf("Grouped by first letter: %v\n", grouped)
}
```

## 性能特性

- **并发优化**: 异步选项使用协程池，避免协程泄漏
- **内存效率**: 就地修改数组，减少内存分配
- **类型安全**: 泛型实现，零运行时类型转换开销
- **算法优化**: 使用高效的排序和去重算法

## 注意事项

1. **数组修改**: Filter、Del、Distinct等函数会直接修改原数组
2. **异步执行**: 使用Async选项时，函数执行顺序不确定
3. **nil处理**: IgnoreNil选项会检查多种nil情况（指针、切片、映射等）
4. **空值判断**: IgnoreEmpty选项使用反射检查零值，有一定性能开销
5. **并发安全**: 异步执行时，确保传入的函数是线程安全的

## 适用场景

- **数据转换**: 批量数据类型转换和格式化
- **数据清洗**: 过滤、去重、排序等数据预处理
- **并发处理**: 需要并发处理大量数据的场景
- **函数式编程**: 需要链式调用和函数组合的场景
- **业务逻辑**: 复杂的数据聚合和分析需求 