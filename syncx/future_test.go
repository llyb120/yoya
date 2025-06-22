package syncx

import (
	"fmt"
	"testing"
	"unsafe"
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
	future, err := asyncAdd(&a, nil)
	if err() != nil {
		fmt.Println(err())
		return
	}

	fmt.Println(future())

	return

	base := uintptr(unsafe.Pointer(&future))
	fmt.Printf("simpleFuture 地址: 0x%x\n", base)

	// 搜索包含 ctx 指针的位置
	// ctxPtr := uintptr(unsafe.Pointer(ctx))
	// fmt.Printf("搜索包含 ctx 指针 0x%x 的位置:\n", ctxPtr)

	found := false
	for offset := -16; offset <= 8; offset += 8 {
		// func() {
		// defer func() { recover() }()

		val := *(*uintptr)(unsafe.Pointer(base + uintptr(offset)))
		// if val == ctxPtr {
		fmt.Printf("  偏移 %2d: 找到 ctx 指针!\n", offset)
		found = true

		// 验证可以通过这个指针访问
		foundCtx := (**futureContext[int])(unsafe.Pointer(val))
		fmt.Printf("    通过指针访问 data: %d\n", (*foundCtx).data)

		// 测试修改
		(*foundCtx).data = 999
		result := future()
		fmt.Printf("    修改后 simpleFuture(): %d\n", result)

		if result == 999 {
			fmt.Printf("    ✓ 偏移 %d 是正确的 ctx 指针位置!\n", offset)
		}

		// 恢复
		(*foundCtx).data = 42
		// }

		break
		// }()
	}

	if !found {
		fmt.Println("  未找到 ctx 指针")
	}

	// fmt.Println("\n=== 测试 Async2 ===")
	// asyncAdd := Async2[int](add)
	// future := asyncAdd(1, 2)
	// result := future()
	// fmt.Printf("future() = %d\n", result)

	// future.Set("type", "async2")

}
