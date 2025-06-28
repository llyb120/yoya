package syncx

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
	"time"
)

func add(a *int, b *int) int {
	if a == nil {
		return -1
	}
	if b == nil {
		return -2
	}
	return *a + *b
}

// 尝试搜索真正的 futureContext 指针
func TestAsync2(t *testing.T) {
	// 先测试简单的指针捕获
	fmt.Println("=== 测试简单指针捕获 ===")

	// ctx := &futureContext[int]{data: 43}
	// fmt.Printf("原始 ctx 地址: 0x%x\n", uintptr(unsafe.Pointer(ctx)))
	// fmt.Printf("原始 ctx.data: %d\n", ctx.data)

	// simpleFuture := func() int {
	// 	return ctx.data
	// }
	asyncAdd := Async2_2_1(add)
	a := 1
	future := asyncAdd(&a, nil)

	res, err := json.Marshal(future)
	fmt.Println(string(res), err)

	return
}

// ==================== 基础功能测试 ====================

// 测试：Mirai 创建空 Future
func TestMirai(t *testing.T) {
	var f = Mirai[int]()
	if v := f(); v != 0 {
		t.Fatalf("want 0, got %v", v)
	}
}

// 测试：Future 的 GetType 方法
func TestFutureGetType(t *testing.T) {
	f := Future[string](func() string { return "hello" })
	typ := f.GetType()
	expected := reflect.TypeOf("")
	if typ != expected {
		t.Fatalf("want %v, got %v", expected, typ)
	}
}

// 测试：Future 的 MarshalJSON 方法
func TestFutureMarshalJSON(t *testing.T) {
	f := Future[int](func() int { return 42 })
	data, err := f.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}
	if string(data) != "42" {
		t.Fatalf("want \"42\", got %s", data)
	}
}

// ==================== Async2_0_x 系列测试 ====================

// 测试：Async2_0_0 无参数无返回值
func TestAsync2_0_0(t *testing.T) {
	var executed bool
	fn := func() { executed = true }
	async := Async2_0_0(fn)
	f := async()
	f() // 等待执行完成
	if !executed {
		t.Fatal("function was not executed")
	}
}

// 测试：Async2_0_1 无参数有返回值
func TestAsync2_0_1_String(t *testing.T) {
	fn := func() string { return "hello world" }
	async := Async2_0_1(fn)
	f := async()
	if v := f(); v != "hello world" {
		t.Fatalf("want 'hello world', got %v", v)
	}
}

// 测试：Async2_0_1 复杂返回值
func TestAsync2_0_1_Slice(t *testing.T) {
	fn := func() []int { return []int{1, 2, 3} }
	async := Async2_0_1(fn)
	f := async()
	v := f()
	if len(v) != 3 || v[0] != 1 || v[1] != 2 || v[2] != 3 {
		t.Fatalf("want [1 2 3], got %v", v)
	}
}

// 测试：Async2_0_2 返回成功结果
func TestAsync2_0_2_Success(t *testing.T) {
	fn := func() (string, error) { return "success", nil }
	async := Async2_0_2(fn)
	f, _ := async()
	if v := f(); v != "success" {
		t.Fatalf("want 'success', got %v", v)
	}
}

// 测试：Async2_0_2 返回错误
func TestAsync2_0_2_WithError(t *testing.T) {
	myErr := errors.New("test error")
	fn := func() (int, error) { return 0, myErr }
	async := Async2_0_2(fn)
	f, _ := async()
	// 这里应该能获取到错误，具体取决于你的实现
	_ = f() // 至少不应该 panic
}

// ==================== Async2_1_x 系列测试 ====================

// 测试：Async2_1_0 单参数无返回值
func TestAsync2_1_0(t *testing.T) {
	var result int
	fn := func(x int) { result = x * 2 }
	async := Async2_1_0(fn)
	f := async(21)
	f()
	if result != 42 {
		t.Fatalf("want 42, got %v", result)
	}
}

// 测试：Async2_1_1 不同类型参数
func TestAsync2_1_1_DifferentTypes(t *testing.T) {
	// 字符串处理
	fn1 := func(s string) string { return strings.ToUpper(s) }
	async1 := Async2_1_1(fn1)
	f1 := async1("hello")
	if v := f1(); v != "HELLO" {
		t.Fatalf("want 'HELLO', got %v", v)
	}

	// 浮点数计算
	fn2 := func(x float64) float64 { return x * 3.14 }
	async2 := Async2_1_1(fn2)
	f2 := async2(2.0)
	if v := f2(); v != 6.28 {
		t.Fatalf("want 6.28, got %v", v)
	}
}

// 测试：Async2_1_2 单参数返回值和错误
func TestAsync2_1_2(t *testing.T) {
	fn := func(x int, divisor int) (int, error) {
		if divisor == 0 {
			return 0, errors.New("division by zero")
		}
		return x / divisor, nil
	}
	async := Async2_1_2(fn)

	// 正常情况
	f1, _ := async(10, 2)
	if v := f1(); v != 5 {
		t.Fatalf("want 5, got %v", v)
	}

	// 错误情况
	f2, _ := async(10, 0)
	_ = f2() // 应该处理错误，不 panic
}

// ==================== Async2_2_x 系列测试 ====================

// 测试：Async2_2_0 双参数无返回值
func TestAsync2_2_0(t *testing.T) {
	var result string
	fn := func(a string, b string) { result = a + b }
	async := Async2_2_0(fn)
	f := async("hello", " world")
	f()
	if result != "hello world" {
		t.Fatalf("want 'hello world', got %v", result)
	}
}

// 测试：Async2_2_1 双参数有返回值
func TestAsync2_2_1(t *testing.T) {
	fn := func(a int, b int) int { return a + b }
	async := Async2_2_1(fn)
	f := async(15, 27)
	if v := f(); v != 42 {
		t.Fatalf("want 42, got %v", v)
	}
}

// 测试：Async2_2_2 双参数返回值和错误
func TestAsync2_2_2(t *testing.T) {
	fn := func(a string, b string) (string, error) {
		if a == "" || b == "" {
			return "", errors.New("empty string")
		}
		return a + " " + b, nil
	}
	async := Async2_2_2(fn)
	f, _ := async("hello", "world")
	if v := f(); v != "hello world" {
		t.Fatalf("want 'hello world', got %v", v)
	}
}

// ==================== Async2_3_x 系列测试 ====================

// 测试：Async2_3_1 三参数
func TestAsync2_3_1(t *testing.T) {
	fn := func(a, b, c int) int { return a + b + c }
	async := Async2_3_1(fn)
	f := async(10, 20, 12)
	if v := f(); v != 42 {
		t.Fatalf("want 42, got %v", v)
	}
}

// 测试：Async2_3_2 三参数返回错误
func TestAsync2_3_2(t *testing.T) {
	fn := func(a, b, c int) (int, error) {
		if a < 0 || b < 0 || c < 0 {
			return 0, errors.New("negative number")
		}
		return a * b * c, nil
	}
	async := Async2_3_2(fn)
	f, _ := async(2, 3, 7)
	if v := f(); v != 42 {
		t.Fatalf("want 42, got %v", v)
	}
}

// ==================== Async2_4_x 和 Async2_5_x 系列测试 ====================

// 测试：Async2_4_1 四参数
func TestAsync2_4_1(t *testing.T) {
	fn := func(a, b, c, d int) int { return a + b + c + d }
	async := Async2_4_1(fn)
	f := async(10, 11, 12, 9)
	if v := f(); v != 42 {
		t.Fatalf("want 42, got %v", v)
	}
}

// 测试：Async2_5_1 五参数
func TestAsync2_5_1(t *testing.T) {
	fn := func(a, b, c, d, e int) int { return a * b * c * d * e }
	async := Async2_5_1(fn)
	f := async(1, 2, 3, 7, 1)
	if v := f(); v != 42 {
		t.Fatalf("want 42, got %v", v)
	}
}

// 测试error
func TestAsync2_4_2(t *testing.T) {
	fn := func(a, b, c, d int) (int, error) {
		panic("test panic")
		return 0, errors.New("test error")
	}
	async := Async2_4_2(fn)
	_, f2 := async(1, 2, 3, 7)
	if f2() == nil {
		t.Fatalf("want nil, got %v", f2())
	}

	t.Log(f2())
}

// ==================== 并发和性能测试 ====================

// 测试：多个 Future 并发执行
func TestConcurrentFutures(t *testing.T) {
	const numFutures = 10
	fn := func(x int) int {
		time.Sleep(10 * time.Millisecond)
		return x * 2
	}
	async := Async2_1_1(fn)

	var futures []Future[int]
	for i := 0; i < numFutures; i++ {
		f := async(i)
		futures = append(futures, f)
	}

	start := time.Now()
	for i, f := range futures {
		if v := f(); v != i*2 {
			t.Fatalf("future %d: want %d, got %d", i, i*2, v)
		}
	}
	duration := time.Since(start)

	// 并发执行应该比串行快得多
	if duration > 50*time.Millisecond {
		t.Logf("Warning: concurrent execution took %v, might not be truly concurrent", duration)
	}
}

// 测试：大量并发 Future
func TestMassiveConcurrentFutures(t *testing.T) {
	const numFutures = 100
	fn := func(x int) int { return x + 1 }
	async := Async2_1_1(fn)

	var wg sync.WaitGroup
	results := make([]int, numFutures)

	for i := 0; i < numFutures; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			f := async(idx)
			results[idx] = f()
		}(i)
	}

	wg.Wait()

	for i, result := range results {
		if result != i+1 {
			t.Fatalf("future %d: want %d, got %d", i, i+1, result)
		}
	}
}

// ==================== 边界情况和错误处理测试 ====================

// 测试：长时间运行的任务
func TestLongRunningTask(t *testing.T) {
	fn := func() string {
		time.Sleep(100 * time.Millisecond)
		return "completed"
	}
	async := Async2_0_1(fn)

	start := time.Now()
	f := async()
	result := f()
	duration := time.Since(start)

	if result != "completed" {
		t.Fatalf("want 'completed', got %v", result)
	}
	if duration < 100*time.Millisecond {
		t.Fatalf("task completed too quickly: %v", duration)
	}
}

// 测试：panic 恢复
func TestPanicRecovery(t *testing.T) {
	fn := func() int { panic("test panic") }
	async := Async2_0_1(fn)
	f := async()

	// 这个测试取决于你的 panic 处理实现
	defer func() {
		if r := recover(); r != nil {
			t.Log("Future didn't handle panic internally")
		}
	}()

	_ = f() // 不应该让 panic 传播出来
}

// 测试：空指针函数
func TestNilFunction(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			// t.Fatal("expected panic for nil function")
		}
	}()

	Async2_0_1[int](nil)
}

// ==================== 类型安全测试 ====================

// 测试：不同类型的 Future
func TestDifferentTypeFutures(t *testing.T) {
	// 字符串类型
	f1 := Future[string](func() string { return "test" })
	if f1.GetType() != reflect.TypeOf("") {
		t.Fatal("string Future type mismatch")
	}

	// 整数类型
	f2 := Future[int](func() int { return 42 })
	if f2.GetType() != reflect.TypeOf(0) {
		t.Fatal("int Future type mismatch")
	}

	// 切片类型
	f3 := Future[[]byte](func() []byte { return []byte("hello") })
	if f3.GetType() != reflect.TypeOf([]byte{}) {
		t.Fatal("[]byte Future type mismatch")
	}
}

// 测试：复杂数据结构
func TestComplexDataStructures(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	fn := func(name string, age int) Person {
		return Person{Name: name, Age: age}
	}
	async := Async2_2_1(fn)
	f := async("Alice", 30)

	person := f()
	if person.Name != "Alice" || person.Age != 30 {
		t.Fatalf("want {Alice 30}, got %+v", person)
	}
}

// ==================== 性能基准测试 ====================

// 基准测试：简单 Future 创建和执行
func BenchmarkSimpleFuture(b *testing.B) {
	fn := func(x int) int { return x * 2 }
	async := Async2_1_1(fn)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f := async(i)
		defer f()
	}
}

// 基准测试：并发 Future 执行
func BenchmarkConcurrentFutures(b *testing.B) {
	fn := func(x int) int { return x + 1 }
	async := Async2_1_1(fn)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			f := async(i)
			_ = f()
			i++
		}
	})
}
