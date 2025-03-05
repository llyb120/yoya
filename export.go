package gotool

import "github.com/llyb120/gotool/collection"

// 这里需要集中高频使用的Api

type Map[K comparable, V any] collection.Map[K, V]

func NewMap[K comparable, V any]() Map[K, V] {
	return collection.NewMap[K, V]()
}
