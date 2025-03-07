package lockx

import (
	"sync"
	"testing"
)

// 基准测试：标准 RWMutex 的 Lock/Unlock
func BenchmarkStdRWMutex_Lock(b *testing.B) {
	var mu sync.RWMutex
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mu.Lock()
		mu.Unlock()
	}
}

// 基准测试：标准 RWMutex 的 RLock/RUnlock
func BenchmarkStdRWMutex_RLock(b *testing.B) {
	var mu sync.RWMutex
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mu.RLock()
		mu.RUnlock()
	}
}

// 基准测试：可重入锁的 Lock/Unlock（非重入情况）
func BenchmarkLock_Lock(b *testing.B) {
	lock := &Lock{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Unlock()
	}
}

// 基准测试：可重入锁的 RLock/RUnlock（非重入情况）
func BenchmarkLock_RLock(b *testing.B) {
	lock := &Lock{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.RLock()
		lock.RUnlock()
	}
}

// 基准测试：可重入锁的 Lock/Unlock（重入情况）
func BenchmarkLock_LockReentrant(b *testing.B) {
	lock := &Lock{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.Lock() // 重入
		lock.Unlock()
		lock.Unlock()
	}
}

// 基准测试：可重入锁的 RLock/RUnlock（重入情况）
func BenchmarkLock_RLockReentrant(b *testing.B) {
	lock := &Lock{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.RLock()
		lock.RLock() // 重入
		lock.RUnlock()
		lock.RUnlock()
	}
}

// 基准测试：可重入锁的读锁升级为写锁
func BenchmarkLock_ReadToWriteUpgrade(b *testing.B) {
	lock := &Lock{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.RLock()
		lock.Lock() // 升级为写锁
		lock.Unlock()
		lock.RUnlock()
	}
}

// 基准测试：可重入锁的写锁获取读锁
func BenchmarkLock_WriteToReadDowngrade(b *testing.B) {
	lock := &Lock{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		lock.Lock()
		lock.RLock() // 在持有写锁的情况下获取读锁
		lock.RUnlock()
		lock.Unlock()
	}
}

// 基准测试：并发读取
func BenchmarkRWLock_ConcurrentRead(b *testing.B) {
	// 标准 RWMutex
	b.Run("StdRWMutex", func(b *testing.B) {
		var mu sync.RWMutex
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				mu.RLock()
				mu.RUnlock()
			}
		})
	})

	// 可重入锁
	b.Run("Lock", func(b *testing.B) {
		lock := &Lock{}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				lock.RLock()
				lock.RUnlock()
			}
		})
	})
}

// 基准测试：并发读写
func BenchmarkRWLock_ConcurrentReadWrite(b *testing.B) {
	// 标准 RWMutex
	b.Run("StdRWMutex", func(b *testing.B) {
		var mu sync.RWMutex
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			// 每个 goroutine 交替进行读写操作
			counter := 0
			for pb.Next() {
				if counter%10 == 0 {
					// 写操作（10% 的概率）
					mu.Lock()
					mu.Unlock()
				} else {
					// 读操作（90% 的概率）
					mu.RLock()
					mu.RUnlock()
				}
				counter++
			}
		})
	})

	// 可重入锁
	b.Run("Lock", func(b *testing.B) {
		lock := &Lock{}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			// 每个 goroutine 交替进行读写操作
			counter := 0
			for pb.Next() {
				if counter%10 == 0 {
					// 写操作（10% 的概率）
					lock.Lock()
					lock.Unlock()
				} else {
					// 读操作（90% 的概率）
					lock.RLock()
					lock.RUnlock()
				}
				counter++
			}
		})
	})
}
