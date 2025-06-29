# errx - 错误处理和多错误管理

errx包提供了增强的错误处理功能，包括多错误收集、异常捕获转换等实用工具，让错误处理变得更加优雅和高效。

## 核心特性

- **MultiError**: 多错误收集和管理，支持并发安全的错误聚合
- **Try/TryDo**: 异常捕获和转换，将panic转换为error返回
- **线程安全**: 所有操作都支持并发访问
- **错误聚合**: 支持批量错误收集和统一处理

## 主要类型和函数

### 1. MultiError - 多错误管理

```go
type MultiError struct {
    // 私有字段，线程安全
}

// MultiError方法
func (e *MultiError) Error() string           // 实现error接口
func (e *MultiError) Add(err error)           // 添加错误
func (e *MultiError) HasError() bool          // 检查是否有错误
```

**MultiError特性:**
- 线程安全的错误收集
- 自动过滤nil错误
- 支持错误信息聚合显示
- 实现标准error接口

**MultiError方法详解:**

#### Error() string
返回所有错误的聚合信息，每个错误占一行。

#### Add(err error)
添加一个错误到集合中。如果传入的err为nil，则会被忽略。该方法是线程安全的。

#### HasError() bool
检查是否包含任何错误。返回true表示有错误，false表示没有错误。

### 2. Try函数 - 异常捕获

```go
// 捕获panic并转换为error
func Try(fn func() error) (err error)

// 捕获panic并返回结果和error
func TryDo[T any](fn func() (T, error)) (v T, err error)
```

**Try函数特性:**
- 自动捕获panic并转换为error
- 支持泛型返回值
- 保持原有的error不变
- 提供安全的函数执行环境

**Try函数详解:**

#### Try(fn func() error) error
执行可能产生panic的函数，如果发生panic，将其捕获并转换为error返回。如果函数正常执行并返回error，则直接返回该error。

#### TryDo[T any](fn func() (T, error)) (T, error)
执行可能产生panic的函数并返回结果。如果发生panic，返回零值和错误。如果函数正常执行，返回实际结果和error。

## 使用示例

### MultiError使用示例

```go
package main

import (
    "fmt"
    "errors"
    "sync"
    "github.com/llyb120/yoya/errx"
)

func main() {
    // 基本使用
    var multiErr errx.MultiError
    
    // 添加多个错误
    multiErr.Add(errors.New("第一个错误"))
    multiErr.Add(errors.New("第二个错误"))
    multiErr.Add(nil) // nil会被忽略
    multiErr.Add(errors.New("第三个错误"))
    
    // 检查是否有错误
    if multiErr.HasError() {
        fmt.Printf("发现错误:\n%s\n", multiErr.Error())
    }
    
    // 并发安全示例
    concurrentExample()
    
    // 批量操作示例
    batchOperationExample()
}

func concurrentExample() {
    var multiErr errx.MultiError
    var wg sync.WaitGroup
    
    // 启动多个协程，每个都可能产生错误
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            
            // 模拟可能失败的操作
            if id%3 == 0 {
                multiErr.Add(fmt.Errorf("协程 %d 执行失败", id))
            }
        }(i)
    }
    
    wg.Wait()
    
    if multiErr.HasError() {
        fmt.Printf("并发操作中的错误:\n%s\n", multiErr.Error())
    } else {
        fmt.Println("所有并发操作都成功了")
    }
}

func batchOperationExample() {
    // 模拟批量文件处理
    files := []string{"file1.txt", "file2.txt", "file3.txt", "file4.txt"}
    var multiErr errx.MultiError
    
    for _, file := range files {
        err := processFile(file)
        multiErr.Add(err) // 自动忽略nil错误
    }
    
    if multiErr.HasError() {
        fmt.Printf("文件处理中的错误:\n%s\n", multiErr.Error())
    } else {
        fmt.Println("所有文件处理成功")
    }
}

func processFile(filename string) error {
    // 模拟文件处理，某些文件可能失败
    if filename == "file2.txt" || filename == "file4.txt" {
        return fmt.Errorf("无法处理文件: %s", filename)
    }
    return nil
}
```

### Try函数使用示例

```go
package main

import (
    "fmt"
    "strconv"
    "github.com/llyb120/yoya/errx"
)

func main() {
    // Try函数示例
    tryExample()
    
    // TryDo函数示例
    tryDoExample()
    
    // 实际应用场景
    realWorldExample()
}

func tryExample() {
    fmt.Println("=== Try函数示例 ===")
    
    // 正常情况
    err := errx.Try(func() error {
        fmt.Println("正常执行的函数")
        return nil
    })
    fmt.Printf("正常执行结果: %v\n", err)
    
    // 返回错误的情况
    err = errx.Try(func() error {
        return fmt.Errorf("函数返回的错误")
    })
    fmt.Printf("函数错误结果: %v\n", err)
    
    // panic的情况
    err = errx.Try(func() error {
        panic("发生了panic")
    })
    fmt.Printf("panic捕获结果: %v\n", err)
}

func tryDoExample() {
    fmt.Println("\n=== TryDo函数示例 ===")
    
    // 正常情况
    result, err := errx.TryDo(func() (string, error) {
        return "成功结果", nil
    })
    fmt.Printf("正常执行: result=%s, err=%v\n", result, err)
    
    // 返回错误的情况
    result, err = errx.TryDo(func() (string, error) {
        return "", fmt.Errorf("函数返回的错误")
    })
    fmt.Printf("函数错误: result=%s, err=%v\n", result, err)
    
    // panic的情况
    result, err = errx.TryDo(func() (string, error) {
        panic("发生了panic")
    })
    fmt.Printf("panic捕获: result=%s, err=%v\n", result, err)
    
    // 数值计算示例
    value, err := errx.TryDo(func() (int, error) {
        return strconv.Atoi("123")
    })
    fmt.Printf("数值转换: value=%d, err=%v\n", value, err)
    
    // panic的数值计算
    value, err = errx.TryDo(func() (int, error) {
        var arr []int
        return arr[10], nil // 会引发panic
    })
    fmt.Printf("数组越界: value=%d, err=%v\n", value, err)
}

func realWorldExample() {
    fmt.Println("\n=== 实际应用场景 ===")
    
    // 批量数据处理，收集所有错误
    data := []string{"123", "456", "abc", "789", "xyz"}
    var multiErr errx.MultiError
    var results []int
    
    for _, item := range data {
        value, err := errx.TryDo(func() (int, error) {
            return strconv.Atoi(item)
        })
        
        if err != nil {
            multiErr.Add(fmt.Errorf("转换 '%s' 失败: %w", item, err))
        } else {
            results = append(results, value)
        }
    }
    
    fmt.Printf("成功转换的数值: %v\n", results)
    if multiErr.HasError() {
        fmt.Printf("转换过程中的错误:\n%s\n", multiErr.Error())
    }
    
    // 复杂操作的错误处理
    complexOperationExample()
}

func complexOperationExample() {
    fmt.Println("\n=== 复杂操作错误处理 ===")
    
    var multiErr errx.MultiError
    
    // 操作1: 可能panic的数组操作
    err := errx.Try(func() error {
        arr := []int{1, 2, 3}
        _ = arr[10] // 数组越界
        return nil
    })
    multiErr.Add(err)
    
    // 操作2: 可能panic的指针操作
    err = errx.Try(func() error {
        var ptr *int
        *ptr = 42 // 空指针引用
        return nil
    })
    multiErr.Add(err)
    
    // 操作3: 正常的错误返回
    err = errx.Try(func() error {
        return fmt.Errorf("业务逻辑错误")
    })
    multiErr.Add(err)
    
    // 操作4: 正常执行
    err = errx.Try(func() error {
        fmt.Println("这个操作正常执行")
        return nil
    })
    multiErr.Add(err) // nil会被忽略
    
    if multiErr.HasError() {
        fmt.Printf("复杂操作中的所有错误:\n%s\n", multiErr.Error())
    }
}
```

### 实际业务场景示例

```go
package main

import (
    "fmt"
    "sync"
    "time"
    "github.com/llyb120/yoya/errx"
)

// 模拟数据验证服务
type ValidationService struct {
    validators []func(data interface{}) error
}

func (vs *ValidationService) ValidateAll(data interface{}) error {
    var multiErr errx.MultiError
    
    for i, validator := range vs.validators {
        err := errx.Try(func() error {
            return validator(data)
        })
        
        if err != nil {
            multiErr.Add(fmt.Errorf("验证器%d失败: %w", i+1, err))
        }
    }
    
    if multiErr.HasError() {
        return &multiErr
    }
    return nil
}

// 模拟批量API调用
func batchAPICall(urls []string) error {
    var multiErr errx.MultiError
    var wg sync.WaitGroup
    var mu sync.Mutex
    
    for _, url := range urls {
        wg.Add(1)
        go func(url string) {
            defer wg.Done()
            
            result, err := errx.TryDo(func() (string, error) {
                return callAPI(url)
            })
            
            mu.Lock()
            if err != nil {
                multiErr.Add(fmt.Errorf("调用 %s 失败: %w", url, err))
            } else {
                fmt.Printf("API调用成功: %s -> %s\n", url, result)
            }
            mu.Unlock()
        }(url)
    }
    
    wg.Wait()
    
    if multiErr.HasError() {
        return &multiErr
    }
    return nil
}

func callAPI(url string) (string, error) {
    // 模拟API调用
    time.Sleep(time.Millisecond * 100)
    
    if url == "http://bad-api.com" {
        panic("API服务器崩溃")
    }
    
    if url == "http://error-api.com" {
        return "", fmt.Errorf("API返回错误")
    }
    
    return fmt.Sprintf("响应来自 %s", url), nil
}

func main() {
    // 数据验证示例
    fmt.Println("=== 数据验证示例 ===")
    vs := &ValidationService{
        validators: []func(interface{}) error{
            func(data interface{}) error {
                if data == nil {
                    return fmt.Errorf("数据不能为空")
                }
                return nil
            },
            func(data interface{}) error {
                // 这个验证器会panic
                panic("验证器内部错误")
            },
            func(data interface{}) error {
                return fmt.Errorf("业务规则验证失败")
            },
        },
    }
    
    err := vs.ValidateAll("test data")
    if err != nil {
        fmt.Printf("验证失败:\n%s\n", err.Error())
    }
    
    // 批量API调用示例
    fmt.Println("\n=== 批量API调用示例 ===")
    urls := []string{
        "http://api1.com",
        "http://api2.com",
        "http://bad-api.com",
        "http://error-api.com",
        "http://api3.com",
    }
    
    err = batchAPICall(urls)
    if err != nil {
        fmt.Printf("批量API调用中的错误:\n%s\n", err.Error())
    }
}
```

## 性能特性

- **线程安全**: MultiError使用读写锁保证并发安全
- **内存优化**: 错误信息按需聚合，避免不必要的字符串拼接
- **零分配**: Try函数在正常路径上没有额外的内存分配
- **panic恢复**: 高效的panic捕获和转换机制

## 注意事项

1. **错误聚合**: MultiError.Error()会将所有错误信息连接，大量错误可能产生很长的字符串
2. **nil过滤**: MultiError.Add()会自动忽略nil错误，无需手动检查
3. **panic转换**: Try函数会将所有panic转换为error，包括runtime panic
4. **并发安全**: MultiError支持并发访问，但Try函数本身不提供并发控制
5. **错误包装**: 建议使用fmt.Errorf和%w动词来包装错误，保持错误链

## 适用场景

- **批量操作**: 需要收集多个操作中的所有错误
- **并发处理**: 多个协程执行任务，需要聚合所有错误
- **数据验证**: 多个验证规则，需要返回所有验证失败的信息
- **异常安全**: 调用可能panic的第三方代码时提供安全保护
- **错误收集**: 复杂业务流程中需要收集和展示所有错误信息

