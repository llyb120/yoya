package syncx

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
	rw sync.RWMutex

	// 写锁持有者的 goroutine ID，如果没有写锁持有者则为 0
	writer int64
	// 写锁的重入计数
	writerCount int

	// 每个 goroutine 持有的读锁计数
	readers map[int64]int

	once sync.Once
}

// NewLock 创建一个新的可重入读写锁
//
//	func NewLock() *Lock {
//		return &Lock{
//			readers: make(map[int64]int),
//		}
//	}
func (l *Lock) init() {
	l.once.Do(func() {
		l.readers = make(map[int64]int)
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

	// 检查当前 goroutine 是否持有读锁
	readCount, hasReadLock := l.readers[gid]

	// 如果持有读锁，需要先释放读锁，然后获取写锁
	if hasReadLock {
		// 暂存读锁计数
		delete(l.readers, gid)
		l.mu.Unlock()

		// 先释放读锁，避免死锁
		// 注意：这里会暂时释放读锁，可能导致其他 goroutine 获取到读锁或写锁
		// 但这是必要的，否则会导致死锁
		l.rw.RUnlock()

		// 获取写锁
		l.rw.Lock()

		l.mu.Lock()
		// 设置写锁持有者和计数
		atomic.StoreInt64(&l.writer, gid)
		l.writerCount = 1
		// 保存之前的读锁计数，以便在解锁时恢复
		l.readers[gid] = -readCount // 负值表示这是从读锁升级的
		l.mu.Unlock()
		return
	}

	l.mu.Unlock()

	// 正常获取写锁
	l.rw.Lock()

	l.mu.Lock()
	atomic.StoreInt64(&l.writer, gid)
	l.writerCount = 1
	l.mu.Unlock()
}

// Unlock 释放写锁
// 如果是重入的写锁，则减少重入计数
// 如果是从读锁升级的写锁，则在完全释放后恢复读锁
func (l *Lock) Unlock() {
	l.init()
	gid := goid.Get()

	l.mu.Lock()
	defer l.mu.Unlock()

	// 检查当前 goroutine 是否持有写锁
	if atomic.LoadInt64(&l.writer) != gid {
		// 不是写锁持有者，不能释放写锁
		return
	}

	// 减少写锁重入计数
	l.writerCount--
	if l.writerCount > 0 {
		// 还有重入的写锁，不完全释放
		return
	}

	// 检查是否是从读锁升级的
	readCount, exists := l.readers[gid]
	if exists && readCount < 0 {
		// 恢复读锁
		l.readers[gid] = -readCount
		atomic.StoreInt64(&l.writer, 0)
		l.rw.Unlock()

		// 重新获取读锁
		l.rw.RLock()
		return
	}

	// 完全释放写锁
	atomic.StoreInt64(&l.writer, 0)
	l.rw.Unlock()
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
		l.mu.Unlock()
		return
	}

	l.mu.Unlock()

	// 正常获取读锁
	l.rw.RLock()

	l.mu.Lock()
	l.readers[gid]++
	l.mu.Unlock()
}

// RUnlock 释放读锁
// 如果是重入的读锁，则减少重入计数
func (l *Lock) RUnlock() {
	l.init()
	gid := goid.Get()

	l.mu.Lock()
	defer l.mu.Unlock()

	// 检查当前 goroutine 是否持有读锁
	count, ok := l.readers[gid]
	if !ok || count <= 0 {
		// 不持有读锁或者是从读锁升级的写锁，不能释放
		return
	}

	// 减少读锁重入计数
	l.readers[gid]--
	if l.readers[gid] > 0 {
		// 还有重入的读锁，不完全释放
		return
	}

	// 完全释放读锁
	delete(l.readers, gid)

	// 如果当前 goroutine 持有写锁，则不需要实际释放读锁
	if atomic.LoadInt64(&l.writer) == gid {
		return
	}

	l.rw.RUnlock()
}

// TryLock 尝试获取写锁，如果获取失败则返回 false
func (l *Lock) TryLock() bool {
	l.init()
	gid := goid.Get()

	l.mu.Lock()
	// 如果当前 goroutine 已经持有写锁，增加重入计数
	if atomic.LoadInt64(&l.writer) == gid {
		l.writerCount++
		l.mu.Unlock()
		return true
	}

	// 检查当前 goroutine 是否持有读锁
	readCount, hasReadLock := l.readers[gid]

	// 如果持有读锁，需要尝试升级为写锁
	if hasReadLock {
		// 暂存读锁计数
		delete(l.readers, gid)
		l.mu.Unlock()

		// 先释放读锁
		l.rw.RUnlock()

		// 尝试获取写锁
		if !l.rw.TryLock() {
			// 获取失败，恢复读锁
			l.rw.RLock()
			l.mu.Lock()
			l.readers[gid] = readCount
			l.mu.Unlock()
			return false
		}

		l.mu.Lock()
		// 设置写锁持有者和计数
		atomic.StoreInt64(&l.writer, gid)
		l.writerCount = 1
		// 保存之前的读锁计数，以便在解锁时恢复
		l.readers[gid] = -readCount // 负值表示这是从读锁升级的
		l.mu.Unlock()
		return true
	}

	l.mu.Unlock()

	// 尝试获取写锁
	if !l.rw.TryLock() {
		return false
	}

	l.mu.Lock()
	atomic.StoreInt64(&l.writer, gid)
	l.writerCount = 1
	l.mu.Unlock()
	return true
}

// TryRLock 尝试获取读锁，如果获取失败则返回 false
func (l *Lock) TryRLock() bool {
	l.init()
	gid := goid.Get()

	l.mu.Lock()
	// 如果当前 goroutine 持有写锁，则不需要实际获取读锁
	if atomic.LoadInt64(&l.writer) == gid {
		l.readers[gid]++
		l.mu.Unlock()
		return true
	}

	// 如果当前 goroutine 已经持有读锁，增加重入计数
	if count, ok := l.readers[gid]; ok && count > 0 {
		l.readers[gid]++
		l.mu.Unlock()
		return true
	}

	l.mu.Unlock()

	// 尝试获取读锁
	if !l.rw.TryRLock() {
		return false
	}

	l.mu.Lock()
	l.readers[gid]++
	l.mu.Unlock()
	return true
}
