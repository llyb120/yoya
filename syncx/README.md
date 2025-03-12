### Async
Async函数用于将任意函数转换为异步方法，转换后的方法入参完全相同，出参则会返回一个泛型的指针
设想一个场景，我需要获取一个时间传递给其它函数使用（假设该操作非常耗时），但其它函数又不会立刻使用，中间会有一些额外的操作
通常我们会声明一个变量，然后使用协程进行赋值，而传参又会非常麻烦
使用Async和Await，以上操作会变得异常简单

**参数**
- fn 任意func类型

**返回**
- fn 入参func被包装后的类型，该函数会返回一个泛型指针

### Await
Await会等待Async返回的指针，如果函数运行中又错误，则会返回对应的错误

**参数**
- ptrs 可变参数，Async返回的泛型指针，可以同时等待多个
- timeout 可选，超时时间，如果超过时间还没有完成，则会返回超时的error

**返回**
- err error 所有目标等待完成后，如果没有错误，则返回nil，如果有错误则立刻返回（可能是超时的error)

**示例**

- 基础示例

```
func Test(a, b int) int {
    time.Sleep(10*time.Second)
    return a + b
}

var AsyncTest = Async[int](Test)

func main(){
    ret := AsyncTest(1,2) // ret是一个空的int指针，等Test函数计算完成后会自动更新指针值
    // do sth...
    DoSomeThing(ret)
    // do sth...
}

func DoSomeThing(a *int){
    if err := Await(a); err != nil {
        log.Println(err)
        return
    }
    // do sth...
    log.Println(*a)
}

```

- 超时

```
func Test(a, b int) int {
    time.Sleep(10*time.Second)
    return a + b
}

var AsyncTest = Async[int](Test)

func main(){
    ret := AsyncTest(1,2) // ret是一个空的int指针，等Test函数计算完成后会自动更新指针值
    // do sth...
    DoSomeThing(ret)
    // do sth...
}

func DoSomeThing(a *int){
    if err := Await(a, 5*time.Second); err != nil {
        log.Println(err)
        return
    }
    // do sth...
    // 这里不可能触发，因为必定会超时
    log.Println(*a)
}
```