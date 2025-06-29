## black

`black` 包提供了高性能的对象与字节序列转换功能，支持结构体、切片、映射的序列化与反序列化。

---

## API 详细说明

### Byte2Str

```go
func Byte2Str(bs []byte) string
```

**功能**：使用 `unsafe` 零拷贝地将 `[]byte` 转换为 `string`。

**参数**：
- `bs []byte` - 要转换的字节切片

**返回值**：
- `string` - 转换后的字符串

**特性**：
- 零内存拷贝，性能极高
- 使用 `unsafe` 包直接操作内存指针
- **注意**：返回的字符串与原字节切片共享内存，修改原字节切片会影响字符串

**示例**：
```go
bs := []byte("hello world")
str := black.Byte2Str(bs)
fmt.Println(str) // "hello world"
```

---

### ToBytes

```go
func ToBytes(obj any) ([]byte, error)
```

**功能**：将对象转换为字节切片，支持结构体、切片、映射类型。

**参数**：
- `obj any` - 要转换的对象

**返回值**：
- `[]byte` - 序列化后的字节数据
- `error` - 转换过程中的错误

**支持的类型**：
- **结构体**：必须传入指针类型，使用 `unsafe` 进行内存拷贝
- **切片**：使用 `gob` 编码
- **映射**：使用 `gob` 编码

**错误情况**：
- 传入结构体但不是指针类型：返回 `"struct must be pointer"` 错误
- 不支持的数据类型：返回 `"invalid data size"` 错误

**示例**：
```go
// 结构体转换
type User struct {
    ID   int64
    Name string
}
user := &User{ID: 1, Name: "Alice"}
data, err := black.ToBytes(user)

// 切片转换
nums := []int{1, 2, 3, 4, 5}
data, err := black.ToBytes(nums)

// 映射转换
m := map[string]int{"a": 1, "b": 2}
data, err := black.ToBytes(m)
```

---

### FromBytes

```go
func FromBytes[T any](bs []byte) (T, error)
```

**功能**：将字节切片反序列化为指定类型的对象。

**类型参数**：
- `T` - 目标类型

**参数**：
- `bs []byte` - 要反序列化的字节数据

**返回值**：
- `T` - 反序列化后的对象
- `error` - 反序列化过程中的错误

**支持的类型**：
- **结构体**：使用 `unsafe` 进行内存操作
- **切片**：使用 `gob` 解码
- **映射**：使用 `gob` 解码

**错误情况**：
- 不支持的数据类型：返回 `"invalid data size"` 错误
- `gob` 解码失败：返回相应的解码错误

**示例**：
```go
// 反序列化结构体
user, err := black.FromBytes[User](data)

// 反序列化切片
nums, err := black.FromBytes[[]int](data)

// 反序列化映射
m, err := black.FromBytes[map[string]int](data)
```

---

## 完整使用示例

```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/black"
)

type Person struct {
    ID   int64
    Name string
    Age  int
}

func main() {
    // 1. 结构体序列化与反序列化
    original := &Person{ID: 100, Name: "Bob", Age: 25}
    
    // 序列化
    data, err := black.ToBytes(original)
    if err != nil {
        panic(err)
    }
    
    // 反序列化
    restored, err := black.FromBytes[Person](data)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Original: %+v\n", original)
    fmt.Printf("Restored: %+v\n", restored)
    
    // 2. 切片操作
    numbers := []int{10, 20, 30, 40, 50}
    numData, _ := black.ToBytes(numbers)
    restoredNums, _ := black.FromBytes[[]int](numData)
    fmt.Printf("Numbers: %v -> %v\n", numbers, restoredNums)
    
    // 3. 映射操作
    scores := map[string]int{"Alice": 85, "Bob": 92, "Charlie": 78}
    scoreData, _ := black.ToBytes(scores)
    restoredScores, _ := black.FromBytes[map[string]int](scoreData)
    fmt.Printf("Scores: %v -> %v\n", scores, restoredScores)
    
    // 4. 零拷贝字符串转换
    byteData := []byte("Hello, World!")
    str := black.Byte2Str(byteData)
    fmt.Printf("String: %s\n", str)
}
```

---

## 性能特性

- **结构体转换**：使用 `unsafe` 包直接操作内存，性能极高，但需要注意内存对齐和字节序
- **切片/映射转换**：使用 Go 标准库的 `gob` 包，兼容性好，支持复杂嵌套结构
- **零拷贝转换**：`Byte2Str` 不进行内存拷贝，适用于临时字符串操作

---

## 注意事项

1. **结构体必须使用指针**：`ToBytes` 处理结构体时必须传入指针类型
2. **内存安全**：结构体转换使用 `unsafe` 包，需要确保内存布局一致
3. **字节序**：结构体转换依赖于系统字节序，跨平台使用需注意
4. **共享内存**：`Byte2Str` 返回的字符串与原字节切片共享内存

---

## 安装
```bash
go get github.com/llyb120/yoya/black
```

---

## API 一览
| 函数 | 说明 |
| ---- | ---- |
| `Byte2Str(bs []byte) string` | 使用 `unsafe` 零拷贝地将 `[]byte` 转为 `string` |
| `ToBytes(v any) ([]byte, error)` | 将切片 / 结构体 / map 转为 `[]byte` |
| `FromBytes[T any](bs []byte) (T, error)` | 与 `ToBytes` 相反，将字节切片解码为目标类型 |

> **注意**：
> 1. 当传入结构体时必须是 **指针** 类型，否则会返回错误。
> 2. 结构体转换使用 `unsafe` 做内存复制，请确保字节序与内存布局一致。

---

## 使用示例

### 1. 结构体 ↔ bytes
```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/black"
)

type User struct {
    Id   int64
    Name string
}

func main() {
    u := &User{Id: 1, Name: "yoya"}

    // 结构体 → []byte
    data, err := black.ToBytes(u)
    if err != nil {
        panic(err)
    }

    // []byte → 结构体
    v, err := black.FromBytes[User](data)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", v)
}
```

### 2. Slice ↔ bytes
```go
nums := []int{1,2,3}
bs, _ := black.ToBytes(nums)
clone, _ := black.FromBytes[[]int](bs)
fmt.Println(clone) // [1 2 3]
```

### 3. Map ↔ bytes
```go
mp := map[string]int{"a":1,"b":2}
bs, _ := black.ToBytes(mp)
clone, _ := black.FromBytes[map[string]int](bs)
fmt.Println(clone)
```

---

## 许可协议
本项目遵循 MIT License。 