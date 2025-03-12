### MultiError
在 go 中经常需要返回多个错误，MultiError 可以轻松帮你合并多个错误

#### 示例
```
var merr MultiError
merr.Add(errors.New("error 1"))
merr.Add(errors.New("error 2"))
merr.Add(errors.New("error 3"))
if !merr.HasError() {
    t.Errorf("Expected error, got nil")
}
if merr.Error() != "error 1\nerror 2\nerror 3" {
    t.Errorf("Expected error string, got %s", merr.Error())
}
```

#### 方法列表
- Add(err error) 添加一个错误
- HasError() bool 判断是否存在错误
- Error() string 返回所有错误信息


#### Try/TryDo
有的时候，我们并不想显式处理错误，我只需要进行简单的计算。
Try 和 TryDo 是两个非常相似的函数，它们都接受一个函数作为参数，并返回一个错误。
该函数将会在安全环境下运行，即使遇到panic，也会被捕获，并且返回错误

**参数**
- fn 任意func类型

**返回**
Try 函数只会返回一个错误
- err error 如果fn返回错误，则返回错误，否则返回nil
TryDo 函数可以返回一个值和一个错误
- val 入参函数的返回类型
- err error 如果fn返回错误，则返回错误，否则返回nil

**示例**
```
func Test(a, b int) int {
    time.Sleep(1*time.Second)
    return a + b
}

err := errx.Try(func() error {
    panic("panic")
    return Test(1, 2)
})

val, err := errx.TryDo(func() (int, error) {
    return Test(1, 2)
})
```

