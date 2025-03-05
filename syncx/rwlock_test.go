package syncx

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// 测试写锁的可重入性
func TestLock_WriteLockReentrant(t *testing.T) {
	lock := &Lock{}

	// 第一次获取写锁
	lock.Lock()

	// 第二次获取写锁（重入）
	lock.Lock()

	// 第三次获取写锁（重入）
	lock.Lock()

	// 释放第三次获取的写锁
	lock.Unlock()

	// 释放第二次获取的写锁
	lock.Unlock()

	// 释放第一次获取的写锁
	lock.Unlock()

	// 验证锁已完全释放
	done := make(chan bool, 1)
	go func() {
		lock.Lock()
		lock.Unlock()
		done <- true
	}()

	select {
	case <-done:
		// 成功
	case <-time.After(time.Second):
		t.Fatal("锁未被完全释放")
	}
}

// 测试读锁的可重入性
func TestLock_ReadLockReentrant(t *testing.T) {
	lock := &Lock{}

	// 第一次获取读锁
	lock.RLock()

	// 第二次获取读锁（重入）
	lock.RLock()

	// 第三次获取读锁（重入）
	lock.RLock()

	// 释放第三次获取的读锁
	lock.RUnlock()

	// 释放第二次获取的读锁
	lock.RUnlock()

	// 释放第一次获取的读锁
	lock.RUnlock()

	// 验证锁已完全释放
	done := make(chan bool, 1)
	go func() {
		lock.Lock()
		lock.Unlock()
		done <- true
	}()

	select {
	case <-done:
		// 成功
	case <-time.After(time.Second):
		t.Fatal("锁未被完全释放")
	}
}

// 测试读锁升级为写锁
func TestLock_ReadToWriteUpgrade(t *testing.T) {
	lock := &Lock{}

	// 获取读锁
	lock.RLock()

	// 尝试升级为写锁
	lock.Lock()

	// 验证写锁已获取
	var writeAcquired int32
	done := make(chan bool, 1)
	go func() {
		// 另一个 goroutine 尝试获取读锁，应该被阻塞
		lock.RLock()
		atomic.StoreInt32(&writeAcquired, 1)
		lock.RUnlock()
		done <- true
	}()

	// 等待一段时间，确保上面的 goroutine 已经尝试获取读锁
	time.Sleep(100 * time.Millisecond)

	// 验证读锁未被获取（因为我们持有写锁）
	if atomic.LoadInt32(&writeAcquired) != 0 {
		t.Fatal("写锁未正确阻塞其他读锁")
	}

	// 释放写锁
	lock.Unlock()

	// 验证读锁已恢复
	done2 := make(chan bool, 1)
	go func() {
		// 另一个 goroutine 尝试获取写锁，应该被阻塞
		lock.Lock()
		lock.Unlock()
		done2 <- true
	}()

	// 等待一段时间，确保上面的 goroutine 已经尝试获取写锁
	time.Sleep(100 * time.Millisecond)

	// 释放读锁
	lock.RUnlock()

	// 验证写锁可以被获取
	select {
	case <-done2:
		// 成功
	case <-time.After(time.Second):
		t.Fatal("读锁未被正确恢复或释放")
	}

	// 验证第一个 goroutine 也完成了
	select {
	case <-done:
		// 成功
	case <-time.After(time.Second):
		t.Fatal("第一个 goroutine 未完成")
	}
}

// 测试持有写锁时获取读锁
func TestLock_WriteToReadDowngrade(t *testing.T) {
	lock := &Lock{}

	// 获取写锁
	lock.Lock()

	// 获取读锁（应该立即成功，因为同一个 goroutine）
	lock.RLock()

	// 释放读锁
	lock.RUnlock()

	// 释放写锁
	lock.Unlock()

	// 验证锁已完全释放
	done := make(chan bool, 1)
	go func() {
		lock.Lock()
		lock.Unlock()
		done <- true
	}()

	select {
	case <-done:
		// 成功
	case <-time.After(time.Second):
		t.Fatal("锁未被完全释放")
	}
}

// 测试 TryLock 功能
func TestLock_TryLock(t *testing.T) {
	lock := &Lock{}

	// 第一次尝试获取写锁应该成功
	if !lock.TryLock() {
		t.Fatal("第一次 TryLock 应该成功")
	}

	// 同一个 goroutine 再次尝试获取写锁应该成功（重入）
	if !lock.TryLock() {
		t.Fatal("重入的 TryLock 应该成功")
	}

	// 另一个 goroutine 尝试获取写锁应该失败
	var success bool
	done := make(chan bool, 1)
	go func() {
		success = lock.TryLock()
		done <- true
	}()

	<-done
	if success {
		t.Fatal("另一个 goroutine 的 TryLock 应该失败")
	}

	// 释放写锁
	lock.Unlock()
	lock.Unlock()

	// 现在另一个 goroutine 应该可以获取写锁
	go func() {
		success = lock.TryLock()
		lock.Unlock()
		done <- true
	}()

	<-done
	if !success {
		t.Fatal("释放锁后 TryLock 应该成功")
	}
}

// 测试 TryRLock 功能
func TestLock_TryRLock(t *testing.T) {
	lock := &Lock{}

	// 第一次尝试获取读锁应该成功
	if !lock.TryRLock() {
		t.Fatal("第一次 TryRLock 应该成功")
	}

	// 同一个 goroutine 再次尝试获取读锁应该成功（重入）
	if !lock.TryRLock() {
		t.Fatal("重入的 TryRLock 应该成功")
	}

	// 另一个 goroutine 尝试获取读锁应该成功
	var success bool
	done := make(chan bool, 1)
	go func() {
		success = lock.TryRLock()
		if success {
			lock.RUnlock()
		}
		done <- true
	}()

	<-done
	if !success {
		t.Fatal("另一个 goroutine 的 TryRLock 应该成功")
	}

	// 另一个 goroutine 尝试获取写锁应该失败
	go func() {
		success = lock.TryLock()
		done <- true
	}()

	<-done
	if success {
		t.Fatal("持有读锁时，另一个 goroutine 的 TryLock 应该失败")
	}

	// 释放读锁
	lock.RUnlock()
	lock.RUnlock()

	// 现在另一个 goroutine 应该可以获取写锁
	go func() {
		success = lock.TryLock()
		if success {
			lock.Unlock()
		}
		done <- true
	}()

	<-done
	if !success {
		t.Fatal("释放读锁后 TryLock 应该成功")
	}
}

// 测试读锁升级为写锁的边界情况
func TestLock_ReadToWriteUpgradeEdgeCases(t *testing.T) {
	lock := &Lock{}

	// 获取多个读锁
	lock.RLock()
	lock.RLock()
	lock.RLock()

	// 尝试升级为写锁
	lock.Lock()

	// 再次获取写锁（重入）
	lock.Lock()

	// 释放一次写锁
	lock.Unlock()

	// 验证写锁仍然有效
	var writeAcquired int32
	done := make(chan bool, 1)
	go func() {
		lock.RLock()
		atomic.StoreInt32(&writeAcquired, 1)
		lock.RUnlock()
		done <- true
	}()

	// 等待一段时间，确保上面的 goroutine 已经尝试获取读锁
	time.Sleep(100 * time.Millisecond)

	// 验证读锁未被获取（因为我们持有写锁）
	if atomic.LoadInt32(&writeAcquired) != 0 {
		t.Fatal("写锁未正确阻塞其他读锁")
	}

	// 释放第二次写锁
	lock.Unlock()

	// 验证读锁已恢复
	// 释放所有读锁
	lock.RUnlock()
	lock.RUnlock()
	lock.RUnlock()

	// 验证锁已完全释放
	select {
	case <-done:
		// 成功
	case <-time.After(time.Second):
		t.Fatal("锁未被完全释放")
	}
}

// 测试并发情况下的锁性能
func TestLock_Concurrent(t *testing.T) {
	lock := &Lock{}
	var counter int64

	// 并发读取
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				lock.RLock()
				_ = atomic.LoadInt64(&counter)
				lock.RUnlock()
			}
		}()
	}

	// 并发写入
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				lock.Lock()
				atomic.AddInt64(&counter, 1)
				lock.Unlock()
			}
		}()
	}

	wg.Wait()

	// 验证计数器的值
	if atomic.LoadInt64(&counter) != 100 {
		t.Fatalf("计数器的值应该是 100，实际是 %d", atomic.LoadInt64(&counter))
	}
}

// 测试错误使用情况
func TestLock_MisuseCase(t *testing.T) {
	lock := &Lock{}

	// 尝试释放未获取的写锁
	lock.Unlock()

	// 尝试释放未获取的读锁
	lock.RUnlock()

	// 获取写锁后，尝试释放读锁
	lock.Lock()
	lock.RUnlock()
	lock.Unlock()

	// 验证锁仍然可以正常工作
	lock.Lock()
	lock.Unlock()

	lock.RLock()
	lock.RUnlock()
}

// 测试读写锁混合使用
func TestLock_MixedUsage(t *testing.T) {
	lock := &Lock{}

	// 获取读锁
	lock.RLock()

	// 获取写锁（会升级）
	lock.Lock()

	// 再次获取读锁
	lock.RLock()

	// 再次获取写锁
	lock.Lock()

	// 按照相反的顺序释放锁
	lock.Unlock()
	lock.RUnlock()
	lock.Unlock()
	lock.RUnlock()

	// 验证锁已完全释放
	done := make(chan bool, 1)
	go func() {
		lock.Lock()
		lock.Unlock()
		done <- true
	}()

	select {
	case <-done:
		// 成功
	case <-time.After(time.Second):
		t.Fatal("锁未被完全释放")
	}
}
