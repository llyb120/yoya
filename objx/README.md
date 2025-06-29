# objx - 对象操作和深度拷贝

objx包提供了强大的对象操作工具，包括深度拷贝、字段选择、对象遍历、类型转换等功能，支持循环引用检测和私有字段访问。

## 核心特性

- **深度拷贝**: 支持所有Go类型的深度克隆，包括私有字段
- **循环引用**: 自动检测和处理循环引用问题
- **字段选择**: 强大的CSS选择器风格的字段提取工具
- **对象遍历**: 递归遍历复杂对象结构
- **类型转换**: 安全的类型转换和赋值操作
- **零值处理**: 提供零值替换和默认值处理

## 主要函数

### 1. 深度拷贝函数

#### DeepClone - 泛型深度拷贝
```go
func DeepClone[T any](src T) (T, error)
```
深拷贝任意对象，包括私有字段，支持所有Go类型。

#### DeepCloneAny - 通用深度拷贝
```go
func DeepCloneAny(src any) (any, error)
```
深拷贝任意类型的对象，返回interface{}类型。

#### MustDeepClone - 必须成功的深度拷贝
```go
func MustDeepClone[T any](src T) T
```
深拷贝，如果出错则panic。

#### MustDeepCloneAny - 必须成功的通用深度拷贝
```go
func MustDeepCloneAny(src any) any
```
深拷贝任意类型，如果出错则panic。

**深度拷贝特性:**
- 支持所有Go类型：基本类型、指针、切片、映射、结构体、接口、通道、函数等
- 自动处理循环引用问题
- 支持私有字段拷贝
- 内存安全，使用unsafe包进行底层操作

### 2. 字段选择函数

#### Pick - 字段选择器
```go
func Pick[T any](src any, rules ...any) []T
```
使用CSS选择器风格的规则从对象中提取字段值。

**Pick选项:**
```go
const (
    Distinct pickOption = iota  // 去重选项
)
```

**Pick特性:**
- 支持CSS选择器语法
- 支持嵌套对象选择
- 支持多规则并发执行
- 支持结果去重
- 类型安全的结果返回

### 3. 对象遍历函数

#### Walk - 对象遍历
```go
func Walk(dest any, fn func(s any, k any, v any) any, opts ...walkOption)
```
递归遍历对象结构，对每个键值对执行指定函数。

### 4. 类型转换和赋值函数

#### Assign - 映射赋值
```go
func Assign[T comparable, K0 any, K1 any](dest map[T]K0, source map[T]K1)
```
将一个映射的内容赋值给另一个映射，支持类型转换。

#### Ensure - 类型确保
```go
func Ensure(objs ...any) bool
```
确保变量是指定类型的值，支持批量类型转换。

#### Cast - 类型转换
```go
func Cast(dest any, src any) error
```
安全的类型转换，将源对象转换为目标类型。

### 5. 零值处理函数

#### Or - 零值替换
```go
func Or[T any](v T, def T) T
```
如果v是零值，则返回def，否则返回v。

## 使用示例

### 深度拷贝示例

```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/objx"
)

type User struct {
    ID       int
    Name     string
    Email    string
    Settings map[string]interface{}
    Friends  []*User
}

func main() {
    // 创建原始对象
    user1 := &User{
        ID:    1,
        Name:  "Alice",
        Email: "alice@example.com",
        Settings: map[string]interface{}{
            "theme":    "dark",
            "language": "en",
            "notifications": map[string]bool{
                "email": true,
                "sms":   false,
            },
        },
    }
    
    user2 := &User{
        ID:    2,
        Name:  "Bob",
        Email: "bob@example.com",
    }
    
    // 设置循环引用
    user1.Friends = []*User{user2}
    user2.Friends = []*User{user1}
    
    // 深度拷贝
    clonedUser, err := objx.DeepClone(user1)
    if err != nil {
        fmt.Printf("Clone error: %v\n", err)
        return
    }
    
    fmt.Printf("Original: %+v\n", user1)
    fmt.Printf("Cloned: %+v\n", clonedUser)
    
    // 修改克隆对象不会影响原对象
    clonedUser.Name = "Alice Clone"
    clonedUser.Settings["theme"] = "light"
    
    fmt.Printf("After modification:\n")
    fmt.Printf("Original name: %s, theme: %s\n", 
        user1.Name, user1.Settings["theme"])
    fmt.Printf("Cloned name: %s, theme: %s\n", 
        clonedUser.Name, clonedUser.Settings["theme"])
    
    // 验证循环引用处理
    fmt.Printf("Original user1 == user1.Friends[0].Friends[0]: %t\n", 
        user1 == user1.Friends[0].Friends[0])
    fmt.Printf("Cloned user != original user: %t\n", 
        clonedUser != user1)
}
```

### 字段选择示例

```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/objx"
)

type Company struct {
    Name    string
    Address string
    Employees []Employee
}

type Employee struct {
    ID       int
    Name     string
    Position string
    Salary   float64
    Manager  *Employee
    Department Department
}

type Department struct {
    Name   string
    Budget float64
}

func main() {
    // 创建测试数据
    manager := Employee{
        ID:       1,
        Name:     "John Manager",
        Position: "Manager",
        Salary:   80000,
        Department: Department{
            Name:   "Engineering",
            Budget: 1000000,
        },
    }
    
    company := Company{
        Name:    "Tech Corp",
        Address: "123 Tech Street",
        Employees: []Employee{
            manager,
            {
                ID:       2,
                Name:     "Alice Developer",
                Position: "Developer",
                Salary:   70000,
                Manager:  &manager,
                Department: Department{
                    Name:   "Engineering",
                    Budget: 1000000,
                },
            },
            {
                ID:       3,
                Name:     "Bob Designer",
                Position: "Designer",
                Salary:   65000,
                Manager:  &manager,
                Department: Department{
                    Name:   "Design",
                    Budget: 500000,
                },
            },
        },
    }
    
    // 1. 提取所有员工姓名
    names := objx.Pick[string](company, "Employees.Name")
    fmt.Printf("Employee names: %v\n", names)
    
    // 2. 提取所有薪资
    salaries := objx.Pick[float64](company, "Employees.Salary")
    fmt.Printf("Salaries: %v\n", salaries)
    
    // 3. 提取部门名称（去重）
    departments := objx.Pick[string](company, "Employees.Department.Name", objx.Distinct)
    fmt.Printf("Departments: %v\n", departments)
    
    // 4. 提取多个字段
    ids := objx.Pick[int](company, "Employees.ID")
    positions := objx.Pick[string](company, "Employees.Position")
    fmt.Printf("IDs: %v\n", ids)
    fmt.Printf("Positions: %v\n", positions)
    
    // 5. 复杂嵌套选择
    budgets := objx.Pick[float64](company, "Employees.Department.Budget", objx.Distinct)
    fmt.Printf("Department budgets: %v\n", budgets)
}
```

### 对象遍历示例

```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/objx"
)

type Config struct {
    Database map[string]interface{}
    Server   map[string]interface{}
    Features map[string]bool
}

func main() {
    config := Config{
        Database: map[string]interface{}{
            "host":     "localhost",
            "port":     5432,
            "username": "admin",
            "password": "secret",
        },
        Server: map[string]interface{}{
            "host": "0.0.0.0",
            "port": 8080,
            "ssl":  true,
        },
        Features: map[string]bool{
            "logging":    true,
            "monitoring": false,
            "caching":    true,
        },
    }
    
    fmt.Println("Configuration traversal:")
    
    // 遍历所有字段
    objx.Walk(config, func(parent any, key any, value any) any {
        fmt.Printf("Key: %v, Value: %v, Type: %T\n", key, value, value)
        return value
    })
    
    // 可以在遍历过程中修改值
    fmt.Println("\nModifying passwords during traversal:")
    objx.Walk(config, func(parent any, key any, value any) any {
        if key == "password" {
            return "***HIDDEN***"
        }
        return value
    })
}
```

### 类型转换示例

```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/objx"
)

func main() {
    // 1. 映射赋值示例
    source := map[string]int{
        "apple":  5,
        "banana": 3,
        "orange": 8,
    }
    
    dest := make(map[string]float64)
    objx.Assign(dest, source)
    
    fmt.Printf("Source: %v\n", source)
    fmt.Printf("Dest: %v\n", dest)
    
    // 2. 类型确保示例
    var a, b, c int
    values := []interface{}{42, "100", 3.14}
    
    success := objx.Ensure(
        &a, values[0],  // 42 -> int
        &b, values[1],  // "100" -> int (会尝试转换)
        &c, values[2],  // 3.14 -> int (会尝试转换)
    )
    
    fmt.Printf("Conversion success: %t\n", success)
    if success {
        fmt.Printf("a: %d, b: %d, c: %d\n", a, b, c)
    }
    
    // 3. 类型转换示例
    type Person struct {
        Name string
        Age  int
    }
    
    data := map[string]interface{}{
        "Name": "Alice",
        "Age":  30,
    }
    
    var person Person
    err := objx.Cast(&person, data)
    if err != nil {
        fmt.Printf("Cast error: %v\n", err)
    } else {
        fmt.Printf("Person: %+v\n", person)
    }
}
```

### 零值处理示例

```go
package main

import (
    "fmt"
    "github.com/llyb120/yoya/objx"
)

func main() {
    // 字符串零值处理
    var name string
    displayName := objx.Or(name, "Anonymous")
    fmt.Printf("Display name: %s\n", displayName)
    
    name = "Alice"
    displayName = objx.Or(name, "Anonymous")
    fmt.Printf("Display name: %s\n", displayName)
    
    // 数值零值处理
    var count int
    displayCount := objx.Or(count, 1)
    fmt.Printf("Display count: %d\n", displayCount)
    
    count = 5
    displayCount = objx.Or(count, 1)
    fmt.Printf("Display count: %d\n", displayCount)
    
    // 指针零值处理
    var ptr *string
    defaultStr := "default"
    result := objx.Or(ptr, &defaultStr)
    fmt.Printf("Result: %s\n", *result)
    
    actualStr := "actual"
    ptr = &actualStr
    result = objx.Or(ptr, &defaultStr)
    fmt.Printf("Result: %s\n", *result)
}
```

### 复杂业务场景示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/llyb120/yoya/objx"
)

type Order struct {
    ID          string
    CustomerID  string
    Items       []OrderItem
    TotalAmount float64
    Status      string
    CreatedAt   time.Time
    Metadata    map[string]interface{}
}

type OrderItem struct {
    ProductID string
    Name      string
    Price     float64
    Quantity  int
    Category  string
}

func main() {
    // 创建订单数据
    orders := []Order{
        {
            ID:         "order-1",
            CustomerID: "customer-1",
            Items: []OrderItem{
                {"prod-1", "Laptop", 999.99, 1, "Electronics"},
                {"prod-2", "Mouse", 29.99, 2, "Electronics"},
            },
            TotalAmount: 1059.97,
            Status:      "completed",
            CreatedAt:   time.Now().AddDate(0, 0, -1),
            Metadata: map[string]interface{}{
                "source":   "web",
                "campaign": "summer_sale",
            },
        },
        {
            ID:         "order-2",
            CustomerID: "customer-2",
            Items: []OrderItem{
                {"prod-3", "Book", 19.99, 3, "Books"},
                {"prod-4", "Pen", 2.99, 5, "Stationery"},
            },
            TotalAmount: 74.92,
            Status:      "pending",
            CreatedAt:   time.Now(),
            Metadata: map[string]interface{}{
                "source": "mobile",
                "notes":  "Gift wrap requested",
            },
        },
    }
    
    // 1. 克隆订单数据进行分析
    ordersCopy, err := objx.DeepClone(orders)
    if err != nil {
        fmt.Printf("Clone error: %v\n", err)
        return
    }
    
    // 2. 提取所有产品名称
    productNames := objx.Pick[string](orders, "Items.Name")
    fmt.Printf("All products: %v\n", productNames)
    
    // 3. 提取所有类别（去重）
    categories := objx.Pick[string](orders, "Items.Category", objx.Distinct)
    fmt.Printf("Categories: %v\n", categories)
    
    // 4. 提取所有价格
    prices := objx.Pick[float64](orders, "Items.Price")
    fmt.Printf("All prices: %v\n", prices)
    
    // 5. 提取元数据来源
    sources := objx.Pick[string](orders, "Metadata.source")
    fmt.Printf("Order sources: %v\n", sources)
    
    // 6. 使用Or处理可能的零值
    for i, order := range ordersCopy {
        // 确保状态有默认值
        order.Status = objx.Or(order.Status, "unknown")
        
        // 确保每个商品都有默认类别
        for j, item := range order.Items {
            order.Items[j].Category = objx.Or(item.Category, "Uncategorized")
        }
        
        ordersCopy[i] = order
    }
    
    fmt.Printf("Processed orders: %+v\n", ordersCopy)
    
    // 7. 类型转换示例：将订单转换为简化格式
    type SimpleOrder struct {
        ID     string
        Total  float64
        Status string
    }
    
    var simpleOrders []SimpleOrder
    for _, order := range orders {
        var simple SimpleOrder
        err := objx.Cast(&simple, map[string]interface{}{
            "ID":     order.ID,
            "Total":  order.TotalAmount,
            "Status": order.Status,
        })
        if err == nil {
            simpleOrders = append(simpleOrders, simple)
        }
    }
    
    fmt.Printf("Simple orders: %+v\n", simpleOrders)
}
```

## 性能特性

- **内存安全**: 使用unsafe包进行底层操作，但保证内存安全
- **循环引用**: 高效的循环引用检测算法
- **并发支持**: Pick函数支持多规则并发执行
- **零拷贝**: 某些操作避免不必要的内存拷贝
- **类型优化**: 针对不同类型进行专门优化

## 注意事项

1. **私有字段**: DeepClone可以拷贝私有字段，但需要谨慎使用
2. **循环引用**: 自动处理循环引用，但可能影响性能
3. **函数拷贝**: 函数类型无法深拷贝，会直接返回原值
4. **并发安全**: 除了Pick函数，其他函数不保证并发安全
5. **内存使用**: 深拷贝会创建完整的对象副本，注意内存使用

## 适用场景

- **对象克隆**: 需要完整复制复杂对象结构的场景
- **数据提取**: 从复杂嵌套结构中提取特定字段
- **配置处理**: 处理复杂的配置对象和数据转换
- **API数据**: 处理来自API的复杂JSON数据结构
- **测试场景**: 创建测试数据的独立副本 