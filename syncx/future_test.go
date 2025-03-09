package syncx

import (
	"errors"
	"testing"
	"time"
)

// 测试正常函数调用
func TestAsyncBasic(t *testing.T) {
	// 测试没有参数的函数
	noArgFn := func() int {
		return 42
	}
	asyncFn := Async[int](noArgFn, 5*time.Second)

	result := asyncFn()
	if err := Await(result); err != nil {
		t.Errorf("期望无错误，但得到: %v", err)
	}
	if *result != 42 {
		t.Errorf("期望结果为42，但得到: %v", result)
	}

	// 测试有参数的函数
	addFn := func(a, b int) int {
		return a + b
	}
	result = Async[int](addFn, 5*time.Second)(10, 20)
	if err := Await(result); err != nil {
		t.Errorf("期望无错误，但得到: %v", err)
	}
	if *result != 30 {
		t.Errorf("期望结果为30，但得到: %v", result)
	}

	// 测试多个返回
	fn0 := func() int {
		time.Sleep(1 * time.Second)
		return 1
	}
	fn1 := func() int {
		time.Sleep(1 * time.Second)
		return 2
	}
	fn2 := func() int {
		panic("failed")
		return 0
	}
	r0 := Async[int](fn0)()
	r1 := Async[int](fn1)()
	r2 := Async[int](fn2)()
	if err := Await(r0, r1, r2); err == nil {
		t.Errorf("期望有错误，但得到: %v", err)
	}

}

// 测试长时间运行的函数
func TestAsyncLongRunning(t *testing.T) {
	slowFn := func() int {
		time.Sleep(100 * time.Millisecond)
		return 99
	}
	result := Async[int](slowFn, 5*time.Second)()
	if err := Await(result); err != nil {
		t.Errorf("期望无错误，但得到: %v", err)
	}
	if *result != 99 {
		t.Errorf("期望结果为99，但得到: %v", result)
	}
}

// 测试函数返回错误
func TestAsyncWithError(t *testing.T) {
	type Result struct {
		Value int
		Err   error
	}
	errorFn := func() Result {
		return Result{0, errors.New("测试错误")}
	}
	result := Async[Result](errorFn, 5*time.Second)()
	if err := Await(result); err != nil {
		t.Errorf("期望无错误，但得到: %v", err)
	}
	if result.Err == nil || result.Err.Error() != "测试错误" {
		t.Errorf("期望错误信息'测试错误'，但得到: %v", result.Err)
	}
}

// // 测试超时
// func TestAsyncTimeout(t *testing.T) {
// 	timeoutFn := func() int {
// 		time.Sleep(200 * time.Millisecond)
// 		return 1
// 	}
// 	// 设置较短的超时时间
// 	future := Async[int](timeoutFn, 50*time.Millisecond)()
// 	_, err := future.Get()
// 	if err == nil {
// 		t.Error("期望超时错误，但没有收到错误")
// 	}
// }

// // 测试参数类型不匹配
// func TestAsyncTypeMismatch(t *testing.T) {
// 	typedFn := func(a string) string {
// 		return a + "!"
// 	}
// 	// 传递错误类型的参数
// 	future := Async[string](typedFn, 5*time.Second)(123)
// 	_, err := future.Get()
// 	if err == nil {
// 		t.Error("期望类型不匹配错误，但没有收到错误")
// 	}
// }

// // 测试panic恢复
// func TestAsyncPanic(t *testing.T) {
// 	panicFn := func() int {
// 		panic("测试panic")
// 		return 0 // 不会执行
// 	}
// 	future := Async[int](panicFn, 5*time.Second)()
// 	_, err := future.Get()
// 	if err == nil {
// 		t.Error("期望panic错误，但没有收到错误")
// 	}
// 	if err != nil && err.Error() != "future panic: 测试panic" {
// 		t.Errorf("期望panic错误信息，但得到: %v", err)
// 	}
// }

// // 测试并发调用
// func TestAsyncConcurrent(t *testing.T) {
// 	addFn := func(a, b int) int {
// 		time.Sleep(50 * time.Millisecond)
// 		return a * b
// 	}
// 	// 并发启动多个future
// 	futures := make([]*Future[int], 5)
// 	for i := 0; i < 5; i++ {
// 		futures[i] = Async[int](addFn, 5*time.Second)(i, i+10)
// 	}
// 	// 收集所有结果
// 	results := make([]int, 5)
// 	for i, future := range futures {
// 		result, err := future.Get()
// 		if err != nil {
// 			t.Errorf("future %d 返回错误: %v", i, err)
// 		}
// 		results[i] = result
// 	}
// 	// 验证结果
// 	expected := []int{0, 11, 24, 39, 56}
// 	for i, exp := range expected {
// 		if results[i] != exp {
// 			t.Errorf("future %d 期望结果 %d，但得到: %d", i, exp, results[i])
// 		}
// 	}
// }

// // 测试多种返回类型
// func TestAsyncDifferentReturnTypes(t *testing.T) {
// 	// 测试返回字符串
// 	strFn := func() string {
// 		return "hello"
// 	}
// 	strFuture := Async[string](strFn, 5*time.Second)()
// 	strResult, err := strFuture.Get()
// 	if err != nil || strResult != "hello" {
// 		t.Errorf("字符串测试失败: %v, %v", strResult, err)
// 	}
// 	// 测试返回结构体
// 	type Person struct {
// 		Name string
// 		Age  int
// 	}
// 	structFn := func() Person {
// 		return Person{"张三", 30}
// 	}
// 	structFuture := Async[Person](structFn, 5*time.Second)()
// 	person, err := structFuture.Get()
// 	if err != nil || person.Name != "张三" || person.Age != 30 {
// 		t.Errorf("结构体测试失败: %v, %v", person, err)
// 	}
// }
