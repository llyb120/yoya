package lockx

import (
	"sync"
	"sync/atomic"

	"github.com/petermattis/goid"
)

// Lock 是一个可重入的读写锁
// 它允许同一个 goroutine 多次获取锁而不会死锁
// 支持读锁升级为写锁
type Lock struct {
	mu sync.Mutex

	// 条件变量，用于等待锁状态变化
	writerCond *sync.Cond
	readerCond *sync.Cond

	// 写锁持有者的 goroutine ID，如果没有写锁持有者则为 0
	writer int64
	// 写锁的重入计数
	writerCount int

	// 每个 goroutine 持有的读锁计数
	readers map[int64]int
	// 活跃读锁总数
	readerCount int32

	// 等待写锁的 goroutines 数量
	waitingWriters int32

	once sync.Once
}

func (l *Lock) init() {
	l.once.Do(func() {
		l.readers = make(map[int64]int)
		l.writerCond = sync.NewCond(&l.mu)
		l.readerCond = sync.NewCond(&l.mu)
	})
}

// Lock 获取写锁
// 如果当前 goroutine 已经持有写锁，则增加重入计数
// 如果当前 goroutine 持有读锁，则升级为写锁
func (l *Lock) Lock() {
	l.init()
	gid := goid.Get()

	l.mu.Lock()
	// 如果当前 goroutine 已经持有写锁，增加重入计数
	if atomic.LoadInt64(&l.writer) == gid {
		l.writerCount++
		l.mu.Unlock()
		return
	}

	// 检查当前 goroutine 是否持有读锁（锁升级场景）
	readCount, hasReadLock := l.readers[gid]

	if hasReadLock {
		// 标记写锁等待者
		atomic.AddInt32(&l.waitingWriters, 1)
		// 保存并取消读锁持有记录
		delete(l.readers, gid)
		atomic.AddInt32(&l.readerCount, -int32(readCount))

		// 等待所有其他读者释放锁
		for atomic.LoadInt32(&l.readerCount) > 0 && atomic.LoadInt64(&l.writer) == 0 {
			l.writerCond.Wait() // 使用条件变量等待
		}

		// 获取写锁并设置写锁持有者
		atomic.StoreInt64(&l.writer, gid)
		l.writerCount = 1
		// 记录已升级的读锁数，以便在释放写锁时可以恢复读锁
		l.readers[gid] = -readCount // 负值表示这是从读锁升级的
		atomic.AddInt32(&l.waitingWriters, -1)
		l.mu.Unlock()
		return
	}

	// 标记有等待的写锁
	atomic.AddInt32(&l.waitingWriters, 1)

	// 等待条件: 没有活跃的读锁且没有其他写锁持有者
	for atomic.LoadInt32(&l.readerCount) > 0 || atomic.LoadInt64(&l.writer) != 0 {
		l.writerCond.Wait() // 使用条件变量等待
	}

	// 获取写锁
	atomic.StoreInt64(&l.writer, gid)
	l.writerCount = 1
	atomic.AddInt32(&l.waitingWriters, -1)
	l.mu.Unlock()
}

// Unlock 释放写锁
// 如果是重入的写锁，则减少重入计数
// 如果是从读锁升级的写锁，则在完全释放后恢复读锁
func (l *Lock) Unlock() {
	l.init()
	gid := goid.Get()

	l.mu.Lock()
	// 检查当前 goroutine 是否持有写锁
	if atomic.LoadInt64(&l.writer) != gid {
		// 不是写锁持有者，不能释放写锁
		l.mu.Unlock()
		return
	}

	// 减少写锁重入计数
	l.writerCount--
	if l.writerCount > 0 {
		// 还有重入的写锁，不完全释放
		l.mu.Unlock()
		return
	}

	// 检查是否是从读锁升级的
	readCount, exists := l.readers[gid]
	if exists && readCount < 0 {
		// 恢复读锁
		actualReadCount := -readCount
		l.readers[gid] = actualReadCount
		atomic.AddInt32(&l.readerCount, int32(actualReadCount))
		atomic.StoreInt64(&l.writer, 0)
		// 通知等待的读者和写者
		l.readerCond.Broadcast()
		l.writerCond.Signal()
		l.mu.Unlock()
		return
	}

	// 完全释放写锁
	atomic.StoreInt64(&l.writer, 0)
	// 通知等待的读者和写者
	l.readerCond.Broadcast()
	l.writerCond.Signal()
	l.mu.Unlock()
}

// RLock 获取读锁
// 如果当前 goroutine 已经持有写锁，则不实际获取读锁，只增加读锁计数
// 如果当前 goroutine 已经持有读锁，则增加重入计数
func (l *Lock) RLock() {
	l.init()
	gid := goid.Get()

	l.mu.Lock()
	// 如果当前 goroutine 持有写锁，则不需要实际获取读锁
	if atomic.LoadInt64(&l.writer) == gid {
		l.readers[gid]++
		l.mu.Unlock()
		return
	}

	// 如果当前 goroutine 已经持有读锁，增加重入计数
	if count, ok := l.readers[gid]; ok && count > 0 {
		l.readers[gid]++
		atomic.AddInt32(&l.readerCount, 1)
		l.mu.Unlock()
		return
	}

	// 如果有等待的写锁，读锁要等待（防止写锁饥饿）
	for atomic.LoadInt32(&l.waitingWriters) > 0 || atomic.LoadInt64(&l.writer) != 0 {
		l.readerCond.Wait() // 使用条件变量等待
	}

	// 获取读锁
	l.readers[gid] = 1
	atomic.AddInt32(&l.readerCount, 1)
	l.mu.Unlock()
}

// RUnlock 释放读锁
// 如果是重入的读锁，则减少重入计数
func (l *Lock) RUnlock() {
	l.init()
	gid := goid.Get()

	l.mu.Lock()
	// 检查当前 goroutine 是否持有读锁
	count, ok := l.readers[gid]
	if !ok || count <= 0 {
		// 不持有读锁或者是从读锁升级的写锁，不能释放
		l.mu.Unlock()
		return
	}

	// 减少读锁重入计数
	l.readers[gid]--
	atomic.AddInt32(&l.readerCount, -1)

	if l.readers[gid] > 0 {
		// 还有重入的读锁，不完全释放
		l.mu.Unlock()
		return
	}

	// 完全释放读锁
	delete(l.readers, gid)
	// 如果没有读锁了，通知等待的写者
	if atomic.LoadInt32(&l.readerCount) == 0 {
		l.writerCond.Broadcast() // 使用Broadcast而不是Signal确保所有等待的写者都被通知
	}
	l.mu.Unlock()
}

// TryLock 尝试获取写锁，如果获取失败则返回 false
func (l *Lock) TryLock() bool {
	l.init()
	gid := goid.Get()

	l.mu.Lock()
	defer l.mu.Unlock()

	// 如果当前 goroutine 已经持有写锁，增加重入计数
	if atomic.LoadInt64(&l.writer) == gid {
		l.writerCount++
		return true
	}

	// 检查当前 goroutine 是否持有读锁（锁升级场景）
	readCount, hasReadLock := l.readers[gid]

	// 如果有活跃的读锁（非当前 goroutine）或其他写锁持有者，则获取失败
	otherReaders := atomic.LoadInt32(&l.readerCount)
	if hasReadLock {
		otherReaders -= int32(readCount)
	}

	if otherReaders > 0 || (atomic.LoadInt64(&l.writer) != 0 && atomic.LoadInt64(&l.writer) != gid) {
		return false
	}

	// 可以获取写锁
	if hasReadLock {
		// 取消读锁记录
		delete(l.readers, gid)
		atomic.AddInt32(&l.readerCount, -int32(readCount))
	}

	// 获取写锁
	atomic.StoreInt64(&l.writer, gid)
	l.writerCount = 1

	if hasReadLock {
		// 记录已升级的读锁数
		l.readers[gid] = -readCount
	}

	return true
}

// TryRLock 尝试获取读锁，如果获取失败则返回 false
func (l *Lock) TryRLock() bool {
	l.init()
	gid := goid.Get()

	l.mu.Lock()
	defer l.mu.Unlock()

	// 如果当前 goroutine 持有写锁，则不需要实际获取读锁
	if atomic.LoadInt64(&l.writer) == gid {
		l.readers[gid]++
		return true
	}

	// 如果当前 goroutine 已经持有读锁，增加重入计数
	if count, ok := l.readers[gid]; ok && count > 0 {
		l.readers[gid]++
		atomic.AddInt32(&l.readerCount, 1)
		return true
	}

	// 如果有等待的写锁或已经有写锁持有者，则获取失败
	if atomic.LoadInt32(&l.waitingWriters) > 0 || atomic.LoadInt64(&l.writer) != 0 {
		return false
	}

	// 获取读锁
	l.readers[gid] = 1
	atomic.AddInt32(&l.readerCount, 1)
	return true
}
