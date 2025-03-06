package gotool

import (
	"github.com/llyb120/gotool/stlx"
)

// 这里需要集中高频使用的Api

type Map[K comparable, V any] stlx.Map[K, V]

type Set[K comparable] stlx.Set[K]

func NewMap[K comparable, V any]() Map[K, V] {
	return stlx.NewMap[K, V]()
}

func NewSet[K comparable]() Set[K] {
	return stlx.NewSet[K]()
}
