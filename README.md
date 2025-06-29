# YOYA - Go 语言扩展工具库

YOYA 是一个现代化的 Go 语言扩展工具库，提供了丰富的功能模块来简化日常开发工作。每个子包都专注于特定领域，采用泛型设计，类型安全且性能优化。

---

## 📦 包结构总览

| 包名 | 功能领域 | 核心特性 | 文档链接 |
|------|----------|----------|----------|
| [**black**](./black/README.md) | 对象序列化 | 高性能字节转换、unsafe优化 | [详细文档](./black/README.md) |
| [**cachex**](./cachex/README.md) | 本地缓存 | 过期管理、并发安全、懒加载 | [详细文档](./cachex/README.md) |
| [**corex**](./corex/README.md) | 核心扩展 | 标准库补充（开发中） | [详细文档](./corex/README.md) |
| [**errx**](./errx/README.md) | 错误处理 | 多错误合并、异常捕获 | [详细文档](./errx/README.md) |
| [**lsx**](./lsx/README.md) | 列表处理 | 函数式编程、并发处理 | [详细文档](./lsx/README.md) |
| [**objx**](./objx/README.md) | 对象操作 | 深拷贝、字段选择、遍历 | [详细文档](./objx/README.md) |
| [**refx**](./refx/README.md) | 反射增强 | 安全反射、动态调用 | [详细文档](./refx/README.md) |
| [**stlx**](./stlx/README.md) | 数据结构 | 高级集合、有序容器 | [详细文档](./stlx/README.md) |
| [**strx**](./strx/README.md) | 字符串处理 | 模式匹配、通配符 | [详细文档](./strx/README.md) |
| [**supx**](./supx/README.md) | 支持工具 | 延迟执行、数据比较 | [详细文档](./supx/README.md) |
| [**syncx**](./syncx/README.md) | 并发控制 | 异步编程、Future模式 | [详细文档](./syncx/README.md) |
| [**tickx**](./tickx/README.md) | 时间处理 | 日期操作、格式化 | [详细文档](./tickx/README.md) |

---

## 🚀 快速开始

### 安装
```bash
go get github.com/llyb120/yoya
```

### 基础示例
```go
package main

import (
    "fmt"
    "time"
    
    "github.com/llyb120/yoya/lsx"
    "github.com/llyb120/yoya/cachex"
    "github.com/llyb120/yoya/objx"
)

func main() {
    // 1. 使用 lsx 进行函数式编程
    numbers := []int{1, 2, 3, 4, 5}
    doubled := lsx.Map(numbers, func(n, i int) int {
        return n * 2
    })
    fmt.Printf("Doubled: %v\n", doubled)
    
    // 2. 使用 cachex 进行缓存
    cache := cachex.NewBaseCache[string, string](cachex.CacheOption{
        DefaultKeyExpire: 5 * time.Minute,
    })
    defer cache.Destroy()
    
    cache.Set("key1", "value1")
    if value, found := cache.Get("key1"); found {
        fmt.Printf("Cached: %s\n", value)
    }
    
    // 3. 使用 objx 进行对象操作
    original := map[string]any{
        "name": "Alice",
        "age":  30,
        "city": "Beijing",
    }
    
    clone, _ := objx.Clone(original)
    picked, _ := objx.Pick(original, "name", "age")
    
    fmt.Printf("Original: %v\n", original)
    fmt.Printf("Clone: %v\n", clone)
    fmt.Printf("Picked: %v\n", picked)
}
```

---

## 🎯 核心特性

### 类型安全
- 全面使用 Go 1.18+ 泛型
- 编译时类型检查
- 减少运行时错误

### 性能优化
- 零拷贝优化（black包）
- 并发处理支持（lsx包）
- 内存友好设计

### 易于使用
- 链式调用支持
- 丰富的示例代码
- 详细的API文档

### 生产就绪
- 完整的单元测试
- 并发安全设计
- 错误处理机制

---

## 📚 详细文档

### 核心功能包

#### [black - 对象序列化](./black/README.md)
高性能的对象与字节序列转换工具，支持结构体、切片、映射的序列化。
```go
// 零拷贝字符串转换
str := black.Byte2Str([]byte("hello"))

// 对象序列化
data, _ := black.ToBytes(&user)
restored, _ := black.FromBytes[User](data)
```

#### [lsx - 函数式编程](./lsx/README.md)
完整的切片处理工具集，支持 Map、Filter、Reduce 等操作，内置并发支持。
```go
// 并发映射
results := lsx.Map(data, processor, lsx.Async)

// 链式操作
lsx.Filter(&data, condition)
lsx.Sort(&data, comparator)
lsx.Distinct(&data)
```

#### [cachex - 本地缓存](./cachex/README.md)
轻量级、类型安全的本地缓存实现，支持过期管理和懒加载。
```go
cache := cachex.NewBaseCache[string, User](cachex.CacheOption{
    DefaultKeyExpire: 10 * time.Minute,
    CheckInterval:    1 * time.Minute,
})

user := cache.GetOrSetFunc("user:123", loadUserFromDB)
```

### 专业工具包

#### [objx - 对象操作](./objx/README.md)
深度对象操作工具，支持克隆、字段选择、遍历等功能。
```go
// 深拷贝（支持循环引用）
clone, _ := objx.Clone(complexObject)

// 字段选择
subset, _ := objx.Pick(data, "name", "email", "profile.age")
```

#### [refx - 反射增强](./refx/README.md)
安全的反射操作封装，简化动态调用和字段访问。
```go
// 安全的字段访问
value, _ := refx.Get(obj, "fieldName")
_ = refx.Set(obj, "fieldName", newValue)

// 动态方法调用
results, _ := refx.Call(obj, "methodName", arg1, arg2)
```

#### [stlx - 高级数据结构](./stlx/README.md)
丰富的数据结构实现，包括有序映射、跳表、双向映射等。
```go
// 有序映射
om := stlx.NewOrderedMap[string, int]()
om.Set("first", 1)
om.Set("second", 2)

// 跳表
sl := stlx.NewSkipList[int]()
sl.Insert(10)
sl.Insert(5)
```

### 辅助工具包

#### [syncx - 并发控制](./syncx/README.md)
异步编程工具，支持 Future 模式和协程管理。
```go
// 异步执行
asyncFunc := syncx.Async[int](expensiveComputation)
result := asyncFunc(params)

// 等待结果
_ = syncx.Await(result, 5*time.Second)
```

#### [errx - 错误处理](./errx/README.md)
多错误合并和异常捕获工具。
```go
// 错误合并
var merr errx.MultiError
merr.Add(err1)
merr.Add(err2)

// 异常捕获
err := errx.Try(func() error {
    // 可能 panic 的代码
    return riskyOperation()
})
```

---

## 🔧 开发指南

### 代码风格
- 遵循 Go 官方代码规范
- 使用有意义的函数和变量名
- 提供完整的文档注释

### 测试覆盖
- 单元测试覆盖率 > 80%
- 并发安全测试
- 基准测试优化

### 性能考虑
- 避免不必要的内存分配
- 使用对象池模式
- 支持并发处理

---

## 🤝 贡献指南

我们欢迎社区贡献！请遵循以下步骤：

1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

### 贡献要求
- 添加相应的单元测试
- 更新相关文档
- 遵循现有代码风格
- 通过所有 CI 检查

---

## 📄 许可证

本项目基于 MIT 许可证开源。详情请参阅 [LICENSE](LICENSE) 文件。

---

## 🙏 致谢

感谢所有为 YOYA 项目做出贡献的开发者们！

---

## 📞 联系方式

- 项目主页: https://github.com/llyb120/yoya
- 问题反馈: https://github.com/llyb120/yoya/issues
- 功能请求: https://github.com/llyb120/yoya/discussions

---

*最后更新：2024年12月* 