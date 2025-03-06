package stlx

import (
	"math/rand"
	"time"

	"github.com/llyb120/gotool/syncx"
)

const (
	maxLevel    = 32  // 跳表的最大层数
	probability = 0.5 // 节点升级到更高层的概率
)

// skipListEntryNode 表示跳表中的一个节点
type skipListEntryNode[K comparable, V any] struct {
	key   K
	value V
	// forward 存储每一层的下一个节点
	forward []*skipListEntryNode[K, V]
}

// newskipListEntryNode 创建一个新的跳表节点
func newSkipListEntryNode[K comparable, V any](key K, value V, level int) *skipListEntryNode[K, V] {
	return &skipListEntryNode[K, V]{
		key:     key,
		value:   value,
		forward: make([]*skipListEntryNode[K, V], level),
	}
}

// SkipMap 是一个协程安全的跳表实现
// 跳表是一种可以用来快速查找的数据结构，类似于平衡树
// 它通过维护多层的链表，使得查找、插入和删除操作的平均时间复杂度为 O(log n)
type SkipMap[K comparable, V any] struct {
	mu       syncx.Lock
	header   *skipListEntryNode[K, V] // 头节点，不存储实际数据
	level    int                      // 当前跳表的最大层数
	length   int                      // 跳表中的元素数量
	less     func(a, b K) bool        // 比较函数，用于确定元素顺序
	randSeed *rand.Rand               // 随机数生成器，用于确定节点层数
}

// NewSkipList 创建一个新的跳表
// less 函数用于比较两个键的大小，如果 a < b 则返回 true
func NewSkipMap[K comparable, V any](less func(a, b K) bool) *SkipMap[K, V] {
	if less == nil {
		return nil
	}

	sl := &SkipMap[K, V]{
		header:   newSkipListEntryNode[K, V](*new(K), *new(V), maxLevel),
		level:    1,
		less:     less,
		randSeed: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	return sl
}

// randomLevel 随机生成一个层数
// 每一层的概率为 probability
func (sl *SkipMap[K, V]) randomLevel() int {
	level := 1
	for level < maxLevel && sl.randSeed.Float64() < probability {
		level++
	}
	return level
}

// Set 添加或更新键值对
func (sl *SkipMap[K, V]) Set(key K, value V) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	// 用于记录每一层需要更新的节点
	update := make([]*skipListEntryNode[K, V], maxLevel)
	current := sl.header

	// 从最高层开始查找
	for i := sl.level - 1; i >= 0; i-- {
		// 在当前层查找最后一个小于 key 的节点
		for current.forward[i] != nil && sl.less(current.forward[i].key, key) {
			current = current.forward[i]
		}
		update[i] = current
	}

	// 移动到第0层的下一个节点
	current = current.forward[0]

	// 如果找到了键，则更新值
	if current != nil && current.key == key {
		current.value = value
		return
	}

	// 生成一个随机层数
	newLevel := sl.randomLevel()

	// 如果新层数大于当前层数，则更新头节点的 forward 指针
	if newLevel > sl.level {
		for i := sl.level; i < newLevel; i++ {
			update[i] = sl.header
		}
		sl.level = newLevel
	}

	// 创建新节点
	newNode := newSkipListEntryNode(key, value, newLevel)

	// 更新所有受影响的节点的 forward 指针
	for i := 0; i < newLevel; i++ {
		newNode.forward[i] = update[i].forward[i]
		update[i].forward[i] = newNode
	}

	sl.length++
}

// Get 获取键对应的值
func (sl *SkipMap[K, V]) Get(key K) (V, bool) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	current := sl.header

	// 从最高层开始查找
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && sl.less(current.forward[i].key, key) {
			current = current.forward[i]
		}
	}

	current = current.forward[0]

	// 如果找到了键，则返回值
	if current != nil && current.key == key {
		return current.value, true
	}

	// 未找到键，返回零值
	return *new(V), false
}

// Del 删除键对应的值，并返回被删除的值
func (sl *SkipMap[K, V]) Del(key K) V {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	// 用于记录每一层需要更新的节点
	update := make([]*skipListEntryNode[K, V], maxLevel)
	current := sl.header

	// 从最高层开始查找
	for i := sl.level - 1; i >= 0; i-- {
		for current.forward[i] != nil && sl.less(current.forward[i].key, key) {
			current = current.forward[i]
		}
		update[i] = current
	}

	current = current.forward[0]

	// 如果找到了键，则删除节点
	if current != nil && current.key == key {
		value := current.value

		// 更新所有受影响的节点的 forward 指针
		for i := 0; i < sl.level; i++ {
			if update[i].forward[i] != current {
				break
			}
			update[i].forward[i] = current.forward[i]
		}

		// 更新跳表的最大层数
		for sl.level > 1 && sl.header.forward[sl.level-1] == nil {
			sl.level--
		}

		sl.length--
		return value
	}

	// 未找到键，返回零值
	return *new(V)
}

// Len 返回跳表中的元素数量
func (sl *SkipMap[K, V]) Len() int {
	sl.mu.RLock()
	defer sl.mu.RUnlock()
	return sl.length
}

// Clear 清空跳表
func (sl *SkipMap[K, V]) Clear() {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.header = newSkipListEntryNode[K, V](*new(K), *new(V), maxLevel)
	sl.level = 1
	sl.length = 0
}

// Keys 返回跳表中的所有键，按顺序排列
func (sl *SkipMap[K, V]) Keys() []K {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	keys := make([]K, 0, sl.length)
	current := sl.header.forward[0]

	for current != nil {
		keys = append(keys, current.key)
		current = current.forward[0]
	}

	return keys
}

// Vals 返回跳表中的所有值，按键的顺序排列
func (sl *SkipMap[K, V]) Vals() []V {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	values := make([]V, 0, sl.length)
	current := sl.header.forward[0]

	for current != nil {
		values = append(values, current.value)
		current = current.forward[0]
	}

	return values
}

// For 遍历跳表中的所有键值对
// 如果回调函数返回 false，则停止遍历
func (sl *SkipMap[K, V]) For(fn func(key K, value V) bool) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	current := sl.header.forward[0]

	for current != nil {
		if !fn(current.key, current.value) {
			break
		}
		current = current.forward[0]
	}
}
