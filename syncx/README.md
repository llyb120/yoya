# syncx - 并发控制和异步编程

syncx包提供了强大的并发控制和异步编程工具，包括Future模式、协程组管理、线程本地存储等功能。

## 核心特性

- **Async2系列**: 新一代异步函数包装器，支持多种参数和返回值组合
- **Future模式**: 手动控制的Future实现，支持异步结果获取
- **协程组管理**: Group类型，支持错误收集和超时控制
- **线程本地存储**: Holder类型，支持协程级别的数据存储
- **对象池**: 增强的Pool实现，支持自定义回收逻辑

## 主要类型和函数

### 1. Async2系列函数（推荐使用）

**旧版Async函数已废弃，请使用Async2系列函数**

#### 无参数函数异步化
```go
// Async2_0_0: func() -> func() Future[any]
func Async2_0_0(fn func()) func() Future[any]

// Async2_0_1: func() T -> func() Future[T]  
func Async2_0_1[T any](fn func() T) func() Future[T]

// Async2_0_2: func() (R0, R1) -> func() (Future[R0], Future[R1])
func Async2_0_2[R0 any, R1 any](fn func() (R0, R1)) func() (Future[R0], Future[R1])

// Async2_0_3: func() (R0, R1, R2) -> func() (Future[R0], Future[R1], Future[R2])
func Async2_0_3[R0 any, R1 any, R2 any](fn func() (R0, R1, R2)) func() (Future[R0], Future[R1], Future[R2])

// Async2_0_4: func() (R0, R1, R2, R3) -> func() (Future[R0], Future[R1], Future[R2], Future[R3])
func Async2_0_4[R0 any, R1 any, R2 any, R3 any](fn func() (R0, R1, R2, R3)) func() (Future[R0], Future[R1], Future[R2], Future[R3])
```

#### 单参数函数异步化
```go
// Async2_1_0: func(T) -> func(T) Future[any]
func Async2_1_0[T any](fn func(T)) func(T) Future[any]

// Async2_1_1: func(P0) R0 -> func(P0) Future[R0]
func Async2_1_1[P0, R0 any](fn func(P0) R0) func(P0) Future[R0]

// Async2_1_2: func(P0, P1) (R0, R1) -> func(P0, P1) (Future[R0], Future[R1])
func Async2_1_2[P0, P1, R0 any, R1 any](fn func(P0, P1) (R0, R1)) func(P0, P1) (Future[R0], Future[R1])

// ... 更多参数组合
```

#### 双参数函数异步化
```go
// Async2_2_0: func(P0, P1) -> func(P0, P1) Future[any]
func Async2_2_0[P0, P1 any](fn func(P0, P1)) func(P0, P1) Future[any]

// Async2_2_1: func(P0, P1) R0 -> func(P0, P1) Future[R0]
func Async2_2_1[P0, P1, R0 any](fn func(P0, P1) R0) func(P0, P1) Future[R0]

// Async2_2_2: func(P0, P1) (R0, R1) -> func(P0, P1) (Future[R0], Future[R1])
func Async2_2_2[P0, P1, R0 any, R1 any](fn func(P0, P1) (R0, R1)) func(P0, P1) (Future[R0], Future[R1])

// ... 更多返回值组合
```

#### 三参数、四参数、五参数函数异步化
```go
// 三参数系列: Async2_3_0 到 Async2_3_4
func Async2_3_0[P0, P1, P2 any](fn func(P0, P1, P2)) func(P0, P1, P2) Future[any]
func Async2_3_1[P0, P1, P2, R0 any](fn func(P0, P1, P2) R0) func(P0, P1, P2) Future[R0]
// ... 

// 四参数系列: Async2_4_0 到 Async2_4_4
func Async2_4_0[P0, P1, P2, P3 any](fn func(P0, P1, P2, P3)) func(P0, P1, P2, P3) Future[any]
func Async2_4_1[P0, P1, P2, P3, R0 any](fn func(P0, P1, P2, P3) R0) func(P0, P1, P2, P3) Future[R0]
// ...

// 五参数系列: Async2_5_0 到 Async2_5_4  
func Async2_5_0[P0, P1, P2, P3, P4 any](fn func(P0, P1, P2, P3, P4)) func(P0, P1, P2, P3, P4) Future[any]
func Async2_5_1[P0, P1, P2, P3, P4, R0 any](fn func(P0, P1, P2, P3, P4) R0) func(P0, P1, P2, P3, P4) Future[R0]
// ...
```

### 2. Future类型

```go
type Future[T any] interface {
    Get(timeout ...time.Duration) (T, error)
    GetType() reflect.Type
    MarshalJSON() ([]byte, error)
}

// 创建手动控制的Future
func Mirai[T any]() Future[T]
```

**Future方法说明:**
- `Get(timeout ...time.Duration) (T, error)`: 获取异步结果，可选超时时间
- `GetType() reflect.Type`: 获取结果类型信息
- `MarshalJSON() ([]byte, error)`: JSON序列化支持

### 3. 旧版Async函数（已废弃）

```go
// 已废弃: 使用 Async2 系列替代
func Async[T any](fn any) func(...any) *T

// 等待异步结果
func Await(objs ...any) error
```

### 4. Group - 协程组管理

```go
type Group struct {
    // 私有字段
}

// Group方法
func (g *Group) Go(fn func() error)                    // 启动协程
func (g *Group) Wait(timeout ...time.Duration) error  // 等待所有协程结束
func (g *Group) SetLimit(limit int)                   // 设置并发限制
```

**Group使用说明:**
- `Go(fn func() error)`: 在协程组中启动新协程，自动处理panic
- `Wait(timeout ...time.Duration) error`: 等待所有协程完成，支持超时
- `SetLimit(limit int)`: 设置同时运行的协程数量限制

### 5. Holder - 线程本地存储

```go
type Holder[V any] struct {
    InitFunc func() V  // 初始化函数
}

// Holder方法
func (h *Holder[V]) Get() V      // 获取当前协程的值
func (h *Holder[V]) Set(value V) // 设置当前协程的值  
func (h *Holder[V]) Del() V      // 删除当前协程的值
```

**Holder特性:**
- 支持协程级别的数据隔离
- 支持父子协程间的数据继承
- 支持自定义初始化函数

### 6. Pool - 对象池

```go
type Pool[T any] interface {
    Get() (T, func())  // 获取对象和回收函数
}

type PoolOption[T any] struct {
    Finalizer func(*T)  // 回收时的清理函数
    New       func() T  // 创建新对象的函数
}

// 构建对象池
func (opt PoolOption[T]) Build() *pool[T]
```

## 使用示例

### Async2系列使用示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/llyb120/yoya/syncx"
)

func main() {
    // 1. 无参数函数异步化
    asyncPrint := syncx.Async2_0_0(func() {
        fmt.Println("Hello from async!")
    })
    
    future := asyncPrint()
    result, err := future.Get(time.Second * 5)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    }
    
    // 2. 带返回值的函数异步化
    asyncCalc := syncx.Async2_0_1(func() int {
        time.Sleep(time.Second)
        return 42
    })
    
    calcFuture := asyncCalc()
    value, err := calcFuture.Get()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Result: %d\n", value)
    }
    
    // 3. 多返回值函数异步化
    asyncDivide := syncx.Async2_0_2(func() (int, error) {
        return 10 / 2, nil
    })
    
    resultFuture, errFuture := asyncDivide()
    divResult, _ := resultFuture.Get()
    divErr, _ := errFuture.Get()
    
    fmt.Printf("Division result: %v, error: %v\n", divResult, divErr)
    
    // 4. 带参数的函数异步化
    asyncAdd := syncx.Async2_2_1(func(a, b int) int {
        return a + b
    })
    
    addFuture := asyncAdd(10, 20)
    sum, err := addFuture.Get()
    if err != nil {
        fmt.Printf("Error: %v\n", err)
    } else {
        fmt.Printf("Sum: %d\n", sum)
    }
}
```

### Group使用示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/llyb120/yoya/syncx"
)

func main() {
    var g syncx.Group
    
    // 启动多个协程
    for i := 0; i < 5; i++ {
        num := i
        g.Go(func() error {
            time.Sleep(time.Duration(num) * time.Millisecond * 100)
            fmt.Printf("Task %d completed\n", num)
            return nil
        })
    }
    
    // 等待所有协程完成，设置5秒超时
    err := g.Wait(time.Second * 5)
    if err != nil {
        fmt.Printf("Group error: %v\n", err)
    } else {
        fmt.Println("All tasks completed successfully")
    }
}
```

### Holder使用示例

```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/syncx"
)

func main() {
    // 创建一个整数类型的Holder
    holder := &syncx.Holder[int]{
        InitFunc: func() int {
            return 100 // 默认值
        },
    }
    
    // 在主协程中设置值
    holder.Set(42)
    fmt.Printf("Main goroutine value: %d\n", holder.Get())
    
    // 在子协程中访问
    go func() {
        // 子协程会继承父协程的值
        fmt.Printf("Child goroutine inherited value: %d\n", holder.Get())
        
        // 子协程设置自己的值
        holder.Set(999)
        fmt.Printf("Child goroutine new value: %d\n", holder.Get())
    }()
    
    time.Sleep(time.Millisecond * 100)
    
    // 主协程的值不受影响
    fmt.Printf("Main goroutine value after child: %d\n", holder.Get())
}
```

### Pool使用示例

```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/syncx"
)

type Connection struct {
    ID     int
    Active bool
}

func main() {
    // 创建连接池
    poolOpt := syncx.PoolOption[Connection]{
        New: func() Connection {
            return Connection{
                ID:     rand.Intn(1000),
                Active: true,
            }
        },
        Finalizer: func(conn *Connection) {
            // 回收时清理连接
            conn.Active = false
            fmt.Printf("Connection %d cleaned\n", conn.ID)
        },
    }
    
    pool := poolOpt.Build()
    
    // 获取连接
    conn, release := pool.Get()
    fmt.Printf("Got connection: %+v\n", conn)
    
    // 使用连接...
    
    // 释放连接回池中
    release()
}
```

## 性能特性

- **零分配**: Async2系列函数在热路径上避免不必要的内存分配
- **类型安全**: 全泛型实现，编译时类型检查
- **并发安全**: 所有类型都支持并发访问
- **内存优化**: Holder支持协程级别的内存隔离

## 注意事项

1. **版本迁移**: 旧版`Async`函数已废弃，请迁移到`Async2`系列
2. **超时处理**: Future的Get方法支持超时，避免无限等待
3. **错误处理**: Group会收集所有协程的错误和panic
4. **内存泄漏**: 使用Pool时确保调用release函数回收对象
5. **协程限制**: Group的SetLimit可以控制并发数量，避免协程爆炸

## 适用场景

- **异步任务处理**: 将同步函数快速转换为异步执行
- **并发控制**: 需要限制并发数量和错误收集的场景  
- **线程本地存储**: 需要协程级别数据隔离的场景
- **对象池化**: 高频创建销毁对象的性能优化场景