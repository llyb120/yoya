package supx

import (
	"container/list"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// SmartWeakMap 结合了弱引用和 LRU 策略的智能 Map
// - 当键被 GC 时自动清理
// - 超过容量时立即清理最老的条目
// - finalizer 触发时清理过期的条目
type SmartWeakMap[K any, V any] struct {
	mu             sync.RWMutex
	items          map[uintptr]*smartEntry[K, V]
	lruList        *list.List
	maxSize        int
	expireDuration time.Duration
	generation     int64
}

// smartEntry 存储条目信息
type smartEntry[K any, V any] struct {
	value      V
	keyPtr     uintptr
	generation int64
	lastAccess time.Time
	lruElement *list.Element
	weakMap    *SmartWeakMap[K, V]
}

// lruItem LRU 列表中的条目
type lruItem struct {
	keyPtr     uintptr
	generation int64
}

// NewSmartWeakMap 创建智能弱引用 Map
func NewSmartWeakMap[K any, V any](maxSize int, expireDuration time.Duration) *SmartWeakMap[K, V] {
	return &SmartWeakMap[K, V]{
		items:          make(map[uintptr]*smartEntry[K, V]),
		lruList:        list.New(),
		maxSize:        maxSize,
		expireDuration: expireDuration,
	}
}

// Set 设置键值对
func (wm *SmartWeakMap[K, V]) Set(key *K, value V) {
	if key == nil {
		panic("SmartWeakMap key cannot be nil")
	}

	keyPtr := uintptr(unsafe.Pointer(key))
	generation := atomic.AddInt64(&wm.generation, 1)
	now := time.Now()

	wm.mu.Lock()
	defer wm.mu.Unlock()

	// 如果已存在，更新并移到 LRU 前端
	if existingEntry, exists := wm.items[keyPtr]; exists {
		// 清理旧的注册
		unregisterSmartCleanup(keyPtr, existingEntry.generation)

		// 更新条目
		existingEntry.value = value
		existingEntry.lastAccess = now
		existingEntry.generation = generation

		// 移到 LRU 前端
		wm.lruList.MoveToFront(existingEntry.lruElement)
		existingEntry.lruElement.Value.(*lruItem).generation = generation

		// 重新注册（不需要重新设置 finalizer，因为对象没变）
		wm.registerEntryCleanup(keyPtr, generation)
		return
	}

	// 检查容量限制，立即清理最老的条目
	wm.enforceCapacity()

	// 创建新条目
	lruItem := &lruItem{
		keyPtr:     keyPtr,
		generation: generation,
	}

	entry := &smartEntry[K, V]{
		value:      value,
		keyPtr:     keyPtr,
		generation: generation,
		lastAccess: now,
		lruElement: wm.lruList.PushFront(lruItem),
		weakMap:    wm,
	}

	wm.items[keyPtr] = entry

	// 注册清理
	wm.registerEntryCleanup(keyPtr, generation)
	runtime.SetFinalizer(key, createSmartFinalizer[K](keyPtr))
}

// Get 获取值
func (wm *SmartWeakMap[K, V]) Get(key *K) (V, bool) {
	if key == nil {
		var zero V
		return zero, false
	}

	keyPtr := uintptr(unsafe.Pointer(key))

	wm.mu.Lock()
	defer wm.mu.Unlock()

	entry, exists := wm.items[keyPtr]
	if !exists {
		var zero V
		return zero, false
	}

	// 更新访问时间并移到 LRU 前端
	entry.lastAccess = time.Now()
	wm.lruList.MoveToFront(entry.lruElement)

	return entry.value, true
}

// Delete 手动删除
func (wm *SmartWeakMap[K, V]) Delete(key *K) bool {
	if key == nil {
		return false
	}

	keyPtr := uintptr(unsafe.Pointer(key))

	wm.mu.Lock()
	defer wm.mu.Unlock()

	return wm.deleteEntry(keyPtr)
}

// Len 返回当前条目数量
func (wm *SmartWeakMap[K, V]) Len() int {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return len(wm.items)
}

// 内部方法：强制执行容量限制
func (wm *SmartWeakMap[K, V]) enforceCapacity() {
	for len(wm.items) >= wm.maxSize {
		if elem := wm.lruList.Back(); elem != nil {
			lruItem := elem.Value.(*lruItem)
			wm.deleteEntry(lruItem.keyPtr)
		} else {
			break
		}
	}
}

// 内部方法：删除条目
func (wm *SmartWeakMap[K, V]) deleteEntry(keyPtr uintptr) bool {
	entry, exists := wm.items[keyPtr]
	if !exists {
		return false
	}

	// 从注册表中清理
	unregisterSmartCleanup(keyPtr, entry.generation)

	// 从 LRU 列表中移除
	wm.lruList.Remove(entry.lruElement)

	// 从 map 中移除
	delete(wm.items, keyPtr)

	return true
}

// 内部方法：注册条目清理
func (wm *SmartWeakMap[K, V]) registerEntryCleanup(keyPtr uintptr, generation int64) {
	registerSmartCleanup(keyPtr, generation, func() {
		// 直接清理对应的条目
		wm.mu.Lock()
		defer wm.mu.Unlock()

		if entry, exists := wm.items[keyPtr]; exists && entry.generation == generation {
			// 从 LRU 列表中移除
			wm.lruList.Remove(entry.lruElement)
			// 从 map 中移除
			delete(wm.items, keyPtr)
		}

		// 顺便清理过期的条目
		wm.cleanupExpiredEntries()
	})
}

// 新增的辅助方法：清理过期条目
func (wm *SmartWeakMap[K, V]) cleanupExpiredEntries() {
	if wm.expireDuration <= 0 {
		return
	}

	now := time.Now()
	cutoff := now.Add(-wm.expireDuration)

	// 从 LRU 尾部开始清理过期的条目
	for elem := wm.lruList.Back(); elem != nil; {
		lruItem := elem.Value.(*lruItem)
		entry := wm.items[lruItem.keyPtr]

		if entry == nil || entry.generation != lruItem.generation {
			// 条目已被清理或版本不匹配，移除 LRU 项
			nextElem := elem.Prev()
			wm.lruList.Remove(elem)
			elem = nextElem
			continue
		}

		if entry.lastAccess.Before(cutoff) {
			// 过期了，删除
			nextElem := elem.Prev()
			wm.deleteEntry(lruItem.keyPtr)
			elem = nextElem
		} else {
			// 由于 LRU 是按时间排序的，后面的都是更新的
			break
		}
	}
}

// 全局智能清理注册表
var (
	globalSmartCleanupMu       sync.RWMutex
	globalSmartCleanupRegistry = make(map[uintptr]map[int64]func())
)

// registerSmartCleanup 注册智能清理函数
func registerSmartCleanup(keyPtr uintptr, generation int64, cleanup func()) {
	globalSmartCleanupMu.Lock()
	defer globalSmartCleanupMu.Unlock()

	if globalSmartCleanupRegistry[keyPtr] == nil {
		globalSmartCleanupRegistry[keyPtr] = make(map[int64]func())
	}
	globalSmartCleanupRegistry[keyPtr][generation] = cleanup
}

// unregisterSmartCleanup 注销智能清理函数
func unregisterSmartCleanup(keyPtr uintptr, generation int64) {
	globalSmartCleanupMu.Lock()
	defer globalSmartCleanupMu.Unlock()

	if generations := globalSmartCleanupRegistry[keyPtr]; generations != nil {
		delete(generations, generation)
		if len(generations) == 0 {
			delete(globalSmartCleanupRegistry, keyPtr)
		}
	}
}

// cleanupSmartWeakMapEntry 执行智能弱引用清理
func cleanupSmartWeakMapEntry(keyPtr uintptr) {
	globalSmartCleanupMu.Lock()
	cleanups := globalSmartCleanupRegistry[keyPtr]
	delete(globalSmartCleanupRegistry, keyPtr)
	globalSmartCleanupMu.Unlock()

	for _, cleanup := range cleanups {
		cleanup()
	}
}

// createSmartFinalizer 创建智能 finalizer 函数
func createSmartFinalizer[K any](keyPtr uintptr) func(*K) {
	return func(*K) {
		cleanupSmartWeakMapEntry(keyPtr)
	}
}
