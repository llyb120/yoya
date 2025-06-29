## cachex

`cachex` 包提供了一套轻量级、类型安全的本地缓存实现，支持键值过期、定时清理、批量操作等功能。

---

## 核心接口

### Cache[K, V]

```go
type Cache[K comparable, V any] interface {
    Get(key K) (value V, ok bool)
    Gets(keys ...K) []V
    Set(key K, value V)
    SetExpire(key K, value V, expire time.Duration)
    Del(key ...K)
    Clear()
    Destroy()
    GetOrSetFunc(key K, fn func() V) V
}
```

**类型参数**：
- `K comparable` - 键类型，必须是可比较的类型
- `V any` - 值类型，可以是任意类型

---

## 配置选项

### CacheOption

```go
type CacheOption struct {
    Expire           time.Duration  // 整个缓存的生存时间
    DefaultKeyExpire time.Duration  // 单个键的默认过期时间
    CheckInterval    time.Duration  // 过期检查间隔
    Destroy          func()         // 缓存销毁时的回调函数
}
```

**字段说明**：
- `Expire`：整个缓存实例的生存时间，超时后自动销毁整个缓存
- `DefaultKeyExpire`：使用 `Set` 方法时的默认过期时间，0 表示永不过期
- `CheckInterval`：定时检查过期键的间隔，0 表示不进行定时清理
- `Destroy`：缓存销毁时执行的回调函数

---

## 构造函数

### NewBaseCache

```go
func NewBaseCache[K comparable, V any](opts CacheOption) *baseCache[K, V]
```

**功能**：创建一个新的缓存实例。

**参数**：
- `opts CacheOption` - 缓存配置选项

**返回值**：
- `*baseCache[K, V]` - 缓存实例

**示例**：
```go
cache := cachex.NewBaseCache[string, int](cachex.CacheOption{
    Expire:           5 * time.Minute,     // 5分钟后销毁
    DefaultKeyExpire: 1 * time.Minute,     // 键默认1分钟过期
    CheckInterval:    30 * time.Second,    // 30秒检查一次
    Destroy: func() {
        fmt.Println("Cache destroyed")
    },
})
```

---

## API 详细说明

### Get

```go
func (c *baseCache[K, V]) Get(key K) (V, bool)
```

**功能**：根据键获取值。

**参数**：
- `key K` - 要查找的键

**返回值**：
- `V` - 对应的值
- `bool` - 是否找到，true 表示找到，false 表示未找到或已过期

**示例**：
```go
value, found := cache.Get("user:123")
if found {
    fmt.Printf("Found: %v\n", value)
} else {
    fmt.Println("Not found or expired")
}
```

---

### Gets

```go
func (c *baseCache[K, V]) Gets(keys ...K) []V
```

**功能**：批量获取多个键的值。

**参数**：
- `keys ...K` - 要查找的键列表

**返回值**：
- `[]V` - 找到的值列表（不包含未找到的键）

**示例**：
```go
values := cache.Gets("key1", "key2", "key3")
fmt.Printf("Found %d values\n", len(values))
```

---

### Set

```go
func (c *baseCache[K, V]) Set(key K, value V)
```

**功能**：设置键值对，使用默认过期时间。

**参数**：
- `key K` - 键
- `value V` - 值

**示例**：
```go
cache.Set("user:123", "Alice")
```

---

### SetExpire

```go
func (c *baseCache[K, V]) SetExpire(key K, value V, expire time.Duration)
```

**功能**：设置键值对并指定过期时间。

**参数**：
- `key K` - 键
- `value V` - 值
- `expire time.Duration` - 过期时间，0 表示永不过期

**示例**：
```go
// 设置5分钟后过期
cache.SetExpire("session:abc", sessionData, 5*time.Minute)

// 设置永不过期
cache.SetExpire("config:db", dbConfig, 0)
```

---

### Del

```go
func (c *baseCache[K, V]) Del(key ...K)
```

**功能**：删除一个或多个键。

**参数**：
- `key ...K` - 要删除的键列表

**示例**：
```go
// 删除单个键
cache.Del("user:123")

// 删除多个键
cache.Del("key1", "key2", "key3")
```

---

### Clear

```go
func (c *baseCache[K, V]) Clear()
```

**功能**：清空缓存中的所有键值对。

**示例**：
```go
cache.Clear()
```

---

### Destroy

```go
func (c *baseCache[K, V]) Destroy()
```

**功能**：销毁缓存实例，停止后台协程并执行销毁回调。

**示例**：
```go
cache.Destroy()
```

---

### GetOrSetFunc

```go
func (c *baseCache[K, V]) GetOrSetFunc(key K, fn func() V) V
```

**功能**：获取键对应的值，如果不存在则执行函数生成值并存储。

**参数**：
- `key K` - 键
- `fn func() V` - 生成值的函数

**返回值**：
- `V` - 获取到的值或新生成的值

**特性**：
- 线程安全的懒加载模式
- 避免缓存穿透问题
- 使用双重检查锁定模式

**示例**：
```go
user := cache.GetOrSetFunc("user:123", func() User {
    // 从数据库加载用户信息
    return loadUserFromDB(123)
})
```

---

## 完整使用示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/llyb120/yoya/cachex"
)

type User struct {
    ID   int
    Name string
    Age  int
}

func main() {
    // 创建缓存
    cache := cachex.NewBaseCache[string, User](cachex.CacheOption{
        Expire:           10 * time.Minute,    // 10分钟后销毁整个缓存
        DefaultKeyExpire: 2 * time.Minute,     // 键默认2分钟过期
        CheckInterval:    1 * time.Minute,     // 1分钟检查一次过期键
        Destroy: func() {
            fmt.Println("Cache has been destroyed")
        },
    })
    defer cache.Destroy()

    // 基本操作
    user1 := User{ID: 1, Name: "Alice", Age: 25}
    cache.Set("user:1", user1)

    // 设置自定义过期时间
    user2 := User{ID: 2, Name: "Bob", Age: 30}
    cache.SetExpire("user:2", user2, 5*time.Minute)

    // 获取值
    if value, found := cache.Get("user:1"); found {
        fmt.Printf("Found user: %+v\n", value)
    }

    // 批量获取
    users := cache.Gets("user:1", "user:2", "user:3")
    fmt.Printf("Found %d users\n", len(users))

    // 懒加载模式
    user3 := cache.GetOrSetFunc("user:3", func() User {
        fmt.Println("Loading user 3 from database...")
        return User{ID: 3, Name: "Charlie", Age: 35}
    })
    fmt.Printf("User 3: %+v\n", user3)

    // 再次获取，应该从缓存返回
    user3Again := cache.GetOrSetFunc("user:3", func() User {
        fmt.Println("This should not be printed")
        return User{}
    })
    fmt.Printf("User 3 (cached): %+v\n", user3Again)

    // 删除操作
    cache.Del("user:1")
    
    // 验证删除
    if _, found := cache.Get("user:1"); !found {
        fmt.Println("User 1 has been deleted")
    }

    // 清空缓存
    cache.Clear()
    fmt.Println("Cache cleared")
}
```

---

## 高级特性

### 1. 自动过期清理

缓存会在后台启动协程定期清理过期的键值对，清理间隔由 `CheckInterval` 配置。

### 2. 整体生命周期管理

通过 `Expire` 配置可以设置整个缓存的生存时间，到期后自动销毁。

### 3. 线程安全

所有操作都是线程安全的，使用读写锁保证并发访问的正确性。

### 4. 内存友好

定期清理过期键值对，避免内存泄漏。

---

## 性能特性

- **读操作**：使用读锁，支持并发读取
- **写操作**：使用写锁，保证数据一致性
- **批量操作**：`Gets` 和 `Del` 支持批量操作，减少锁竞争
- **懒加载**：`GetOrSetFunc` 使用双重检查锁定，避免重复计算

---

## 注意事项

1. **键类型限制**：键类型必须是可比较的（comparable）
2. **内存管理**：长时间运行的程序应该合理设置过期时间和检查间隔
3. **销毁操作**：程序退出前应该调用 `Destroy()` 方法清理资源
4. **并发安全**：所有方法都是线程安全的，可以在多协程环境下使用 