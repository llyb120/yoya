package stlx

import (
	"math/rand"
	"time"
)

// skipListNode 表示跳表中的一个节点
type skipListNode[T any] struct {
	value   T
	forward []*skipListNode[T]
}

// newSkipListNode 创建一个新的跳表节点
func newSkipListNode[T comparable](value T, level int) *skipListNode[T] {
	return &skipListNode[T]{
		value:   value,
		forward: make([]*skipListNode[T], level),
	}
}

// SkipList 是一个协程安全的跳表实现
// 跳表是一种可以用来快速查找的数据结构，类似于平衡树
// 它通过维护多层的链表，使得查找、插入和删除操作的平均时间复杂度为 O(log n)
type SkipList[T comparable] struct {
	mu       lock
	header   *skipListNode[T]  // 头节点，不存储实际数据
	level    int               // 当前跳表的最大层数
	length   int               // 跳表中的元素数量
	less     func(a, b T) bool // 比较函数，用于确定元素顺序
	randSeed *rand.Rand        // 随机数生成器，用于确定节点层数
}

// NewSkipList 创建一个新的跳表
// less 函数用于比较两个键的大小，如果 a < b 则返回 true
func NewSkipList[T comparable](less func(a, b T) bool) *SkipList[T] {
	if less == nil {
		return nil
	}

	sl := &SkipList[T]{
		header:   newSkipListNode(*new(T), maxLevel),
		level:    1,
		less:     less,
		randSeed: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return sl
}

func NewSyncSkipList[T comparable](less func(a, b T) bool) *SkipList[T] {
	sl := NewSkipList[T](less)
	sl.mu.sync = true
	return sl
}

// randomLevel 随机生成一个层数
// 每一层的概率为 probability
func (sl *SkipList[T]) randomLevel() int {
	level := 1
	for level < maxLevel && sl.randSeed.Float64() < probability {
		level++
	}
	return level
}

// Add 插入一个新值
func (sl *SkipList[T]) Add(value T) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.add(value)
}

// Del 删除指定的值
func (sl *SkipList[T]) Del(value T) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	update := make([]*skipListNode[T], maxLevel)
	current := sl.header

	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && sl.less(current.forward[i].value, value) {
			current = current.forward[i]
		}
		update[i] = current
	}

	current = current.forward[0]
	if current != nil && current.value == value {
		for i := 0; i < sl.level; i++ {
			if update[i].forward[i] != current {
				break
			}
			update[i].forward[i] = current.forward[i]
		}

		for sl.level > 1 && sl.header.forward[sl.level-1] == nil {
			sl.level--
		}

		sl.length--
	}

}

// Has 检查值是否存在
func (sl *SkipList[T]) Has(value T) bool {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	current := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && sl.less(current.forward[i].value, value) {
			current = current.forward[i]
		}
	}

	current = current.forward[0]
	return current != nil && current.value == value
}

// Vals 返回有序切片
func (sl *SkipList[T]) Vals() []T {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	return sl.vals()
}

// Len 返回元素数量
func (sl *SkipList[T]) Len() int {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return sl.length
}

// Clear 清空跳表
func (sl *SkipList[T]) Clear() {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.clear()
}

// Get 返回指定下标的元素，如果下标超出范围则返回零值和false
func (sl *SkipList[T]) Get(index int) (T, bool) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	if index < 0 || index >= sl.length {
		return *new(T), false
	}

	current := sl.header.forward[0]
	for i := 0; i < index; i++ {
		current = current.forward[0]
	}

	return current.value, true
}

func (sl *SkipList[T]) For(fn func(value T) bool) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	sl.foreach(fn)
}
