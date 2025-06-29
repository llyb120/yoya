## refx

`refx` 包提供了安全、易用的反射操作封装，简化了 Go 语言中的动态编程需求。相比直接使用 `reflect` 包，`refx` 提供了更友好的 API 和更好的错误处理。

---

## API 详细说明

### Get

```go
func Get(obj any, fieldName string) (any, error)
```

**功能**：安全地获取对象字段值，支持私有字段访问。

**参数**：
- `obj any` - 目标对象
- `fieldName string` - 字段名

**返回值**：
- `any` - 字段值
- `error` - 获取过程中的错误

**示例**：
```go
type User struct {
    Name string
    age  int // 私有字段
}

user := User{Name: "Alice", age: 25}

// 获取公有字段
name, err := refx.Get(user, "Name")
if err == nil {
    fmt.Printf("Name: %v\n", name) // Alice
}

// 获取私有字段（使用 unsafe）
age, err := refx.Get(user, "age")
if err == nil {
    fmt.Printf("Age: %v\n", age) // 25
}
```

---

### Set

```go
func Set(obj any, fieldName string, value any) error
```

**功能**：安全地设置对象字段值，支持私有字段修改。

**参数**：
- `obj any` - 目标对象（必须是指针）
- `fieldName string` - 字段名
- `value any` - 新值

**返回值**：
- `error` - 设置过程中的错误

**示例**：
```go
user := &User{Name: "Alice", age: 25}

// 设置公有字段
err := refx.Set(user, "Name", "Bob")
if err == nil {
    fmt.Printf("New name: %s\n", user.Name) // Bob
}

// 设置私有字段
err = refx.Set(user, "age", 30)
if err == nil {
    // 验证设置成功
    age, _ := refx.Get(user, "age")
    fmt.Printf("New age: %v\n", age) // 30
}
```

---

### Call

```go
func Call(obj any, methodName string, args ...any) ([]any, error)
```

**功能**：动态调用对象方法。

**参数**：
- `obj any` - 目标对象
- `methodName string` - 方法名
- `args ...any` - 方法参数

**返回值**：
- `[]any` - 方法返回值列表
- `error` - 调用过程中的错误

**示例**：
```go
type Calculator struct{}

func (c Calculator) Add(a, b int) int {
    return a + b
}

func (c Calculator) Divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}

calc := Calculator{}

// 调用简单方法
results, err := refx.Call(calc, "Add", 10, 20)
if err == nil {
    sum := results[0].(int)
    fmt.Printf("10 + 20 = %d\n", sum) // 30
}

// 调用返回多值的方法
results, err = refx.Call(calc, "Divide", 10.0, 3.0)
if err == nil {
    quotient := results[0].(float64)
    callErr := results[1]
    fmt.Printf("10 / 3 = %f, error: %v\n", quotient, callErr)
}
```

---

### GetFields

```go
func GetFields(obj any) (map[string]FieldInfo, error)
```

**功能**：获取对象的所有字段信息，包括类型和访问器。

**参数**：
- `obj any` - 目标对象

**返回值**：
- `map[string]FieldInfo` - 字段信息映射
- `error` - 获取过程中的错误

**FieldInfo 结构**：
```go
type FieldInfo struct {
    Type   reflect.Type
    Getter func() (any, error)
    Setter func(any) error
}
```

**示例**：
```go
type Person struct {
    Name    string
    Age     int
    Email   string
    private bool
}

person := &Person{Name: "Alice", Age: 30, Email: "alice@example.com"}

fields, err := refx.GetFields(person)
if err != nil {
    log.Fatal(err)
}

for fieldName, fieldInfo := range fields {
    fmt.Printf("字段: %s, 类型: %s\n", fieldName, fieldInfo.Type)
    
    // 使用 Getter 获取值
    value, err := fieldInfo.Getter()
    if err == nil {
        fmt.Printf("  值: %v\n", value)
    }
}

// 使用 Setter 修改值
if nameField, ok := fields["Name"]; ok {
    err := nameField.Setter("Bob")
    if err == nil {
        fmt.Printf("名字已更改为: %s\n", person.Name)
    }
}
```

---

### GetMethods

```go
func GetMethods(obj any) (map[string]MethodInfo, error)
```

**功能**：获取对象的所有方法信息，包括调用器。

**参数**：
- `obj any` - 目标对象

**返回值**：
- `map[string]MethodInfo` - 方法信息映射
- `error` - 获取过程中的错误

**MethodInfo 结构**：
```go
type MethodInfo struct {
    Type   reflect.Type
    Caller func(...any) ([]any, error)
}
```

**示例**：
```go
type Service struct {
    name string
}

func (s Service) GetName() string {
    return s.name
}

func (s *Service) SetName(name string) {
    s.name = name
}

func (s Service) Process(data string, count int) (string, int) {
    return strings.Repeat(data, count), len(data) * count
}

service := &Service{name: "MyService"}

methods, err := refx.GetMethods(service)
if err != nil {
    log.Fatal(err)
}

for methodName, methodInfo := range methods {
    fmt.Printf("方法: %s, 类型: %s\n", methodName, methodInfo.Type)
}

// 使用 Caller 调用方法
if processMethod, ok := methods["Process"]; ok {
    results, err := processMethod.Caller("hello", 3)
    if err == nil {
        result := results[0].(string)
        length := results[1].(int)
        fmt.Printf("处理结果: %s, 长度: %d\n", result, length)
    }
}
```

---

## 完整使用示例

```go
package main

import (
    "fmt"
    "log"
    "reflect"
    "github.com/llyb120/yoya/refx"
)

type User struct {
    ID       int    `json:"id"`
    Name     string `json:"name"`
    Email    string `json:"email"`
    password string // 私有字段
}

func (u User) GetDisplayName() string {
    return fmt.Sprintf("%s <%s>", u.Name, u.Email)
}

func (u *User) SetPassword(password string) error {
    if len(password) < 6 {
        return fmt.Errorf("密码长度不能少于6位")
    }
    u.password = password
    return nil
}

func (u User) Validate() (bool, []string) {
    var errors []string
    
    if u.Name == "" {
        errors = append(errors, "姓名不能为空")
    }
    if u.Email == "" {
        errors = append(errors, "邮箱不能为空")
    }
    
    return len(errors) == 0, errors
}

func main() {
    user := &User{
        ID:    1,
        Name:  "Alice",
        Email: "alice@example.com",
    }
    
    fmt.Println("=== 字段操作示例 ===")
    fieldOperations(user)
    
    fmt.Println("\n=== 方法调用示例 ===")
    methodCalls(user)
    
    fmt.Println("\n=== 批量字段操作示例 ===")
    batchFieldOperations(user)
    
    fmt.Println("\n=== 实际应用场景 ===")
    realWorldExample()
}

func fieldOperations(user *User) {
    // 读取公有字段
    name, err := refx.Get(user, "Name")
    if err == nil {
        fmt.Printf("用户名: %v\n", name)
    }
    
    // 读取私有字段
    password, err := refx.Get(user, "password")
    if err == nil {
        fmt.Printf("密码: %v\n", password)
    }
    
    // 修改公有字段
    err = refx.Set(user, "Name", "Bob")
    if err == nil {
        fmt.Printf("修改后的用户名: %s\n", user.Name)
    }
    
    // 修改私有字段
    err = refx.Set(user, "password", "secret123")
    if err == nil {
        // 验证修改成功
        newPassword, _ := refx.Get(user, "password")
        fmt.Printf("设置的密码: %v\n", newPassword)
    }
}

func methodCalls(user *User) {
    // 调用无参方法
    results, err := refx.Call(user, "GetDisplayName")
    if err == nil {
        displayName := results[0].(string)
        fmt.Printf("显示名称: %s\n", displayName)
    }
    
    // 调用有参方法
    results, err = refx.Call(user, "SetPassword", "newpassword123")
    if err == nil {
        if len(results) > 0 && results[0] != nil {
            fmt.Printf("设置密码错误: %v\n", results[0])
        } else {
            fmt.Println("密码设置成功")
        }
    }
    
    // 调用返回多值的方法
    results, err = refx.Call(user, "Validate")
    if err == nil {
        isValid := results[0].(bool)
        errors := results[1].([]string)
        fmt.Printf("验证结果: %v, 错误: %v\n", isValid, errors)
    }
}

func batchFieldOperations(user *User) {
    // 获取所有字段信息
    fields, err := refx.GetFields(user)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("所有字段:")
    for fieldName, fieldInfo := range fields {
        value, err := fieldInfo.Getter()
        if err == nil {
            fmt.Printf("  %s: %v (类型: %s)\n", 
                fieldName, value, fieldInfo.Type)
        }
    }
    
    // 批量设置字段
    updates := map[string]any{
        "ID":    999,
        "Email": "newemail@example.com",
    }
    
    for fieldName, newValue := range updates {
        if field, ok := fields[fieldName]; ok {
            err := field.Setter(newValue)
            if err == nil {
                fmt.Printf("已更新 %s = %v\n", fieldName, newValue)
            }
        }
    }
}

func realWorldExample() {
    // 模拟 ORM 场景：动态设置结构体字段
    data := map[string]any{
        "ID":    100,
        "Name":  "Charlie",
        "Email": "charlie@example.com",
    }
    
    user := &User{}
    
    // 动态填充结构体
    fields, _ := refx.GetFields(user)
    for key, value := range data {
        if field, ok := fields[key]; ok {
            err := field.Setter(value)
            if err != nil {
                fmt.Printf("设置字段 %s 失败: %v\n", key, err)
            }
        }
    }
    
    fmt.Printf("动态填充的用户: %+v\n", user)
    
    // 模拟 API 调用场景：动态方法调用
    methods, _ := refx.GetMethods(user)
    
    // 调用验证方法
    if validateMethod, ok := methods["Validate"]; ok {
        results, err := validateMethod.Caller()
        if err == nil {
            isValid := results[0].(bool)
            errors := results[1].([]string)
            
            if isValid {
                fmt.Println("用户数据验证通过")
            } else {
                fmt.Printf("用户数据验证失败: %v\n", errors)
            }
        }
    }
}
```

---

## 高级特性

### 1. 私有字段访问
使用 `unsafe` 包安全地访问和修改私有字段：

```go
type Config struct {
    host string // 私有字段
    port int    // 私有字段
}

config := &Config{}
refx.Set(config, "host", "localhost")
refx.Set(config, "port", 8080)

host, _ := refx.Get(config, "host") // "localhost"
```

### 2. 类型安全的字段访问
```go
fields, _ := refx.GetFields(obj)
if nameField, ok := fields["Name"]; ok {
    // 类型检查
    if nameField.Type.Kind() == reflect.String {
        nameField.Setter("new name")
    }
}
```

### 3. 方法签名检查
```go
methods, _ := refx.GetMethods(obj)
if method, ok := methods["Process"]; ok {
    methodType := method.Type
    fmt.Printf("方法有 %d 个参数\n", methodType.NumIn())
    fmt.Printf("方法有 %d 个返回值\n", methodType.NumOut())
}
```

---

## 性能特性

- **缓存优化**：字段和方法信息会被缓存，避免重复反射操作
- **错误安全**：所有操作都有完善的错误处理，不会导致 panic
- **类型检查**：在运行时进行类型兼容性检查

---

## 使用场景

1. **ORM 框架**：动态映射数据库字段到结构体
2. **配置管理**：从配置文件动态设置对象属性
3. **API 框架**：动态调用处理器方法
4. **序列化/反序列化**：自定义的 JSON/XML 处理
5. **测试工具**：访问私有字段进行单元测试

---

## 注意事项

1. **性能考虑**：反射操作比直接访问慢，避免在热点代码中频繁使用
2. **类型安全**：返回值需要进行类型断言
3. **并发安全**：所有操作都是线程安全的
4. **内存安全**：使用 `unsafe` 访问私有字段时要确保对象生命周期

---

## 安装
```bash
go get github.com/llyb120/yoya/refx
```

---

## 主要能力
| 功能 | 核心 API | 说明 |
| ---- | -------- | ---- |
| 读写字段 | `Get` `Set` | 根据字段名直接获取 / 设置，支持私有字段（`unsafe`） |
| 调用方法 | `Call` | 根据方法名动态调用并返回 `[]any` |
| 列出字段 | `GetFields` | 返回字段映射（含类型和 Getter/Setter） |
| 列出方法 | `GetMethods` | 返回方法映射（含调用封装） |

错误都会被封装为 `error` 返回，避免 panic 崩溃。

---

## 快速示例
```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/refx"
)

type User struct {
    Name string
}

func (u *User) Say(msg string) string {
    return fmt.Sprintf("%s: %s", u.Name, msg)
}

func main() {
    u := &User{Name: "Tom"}

    // 修改字段
    _ = refx.Set(u, "Name", "Jerry")

    // 调用方法
    res, _ := refx.Call(u, "Say", "hello")
    fmt.Println(res[0]) // Jerry: hello
}
```

---

## 高级：批量获取字段 Getter / Setter
```go
fields := refx.GetFields(&User{}, refx.IgnoreFunc)

setName := fields["Name"].Set
getName := fields["Name"].Get

_ = setName("Alice")
val, _ := getName()
fmt.Println(val) // Alice
```

---

## 许可协议
本项目遵循 MIT License。 