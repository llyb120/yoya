package gotool

import "github.com/llyb120/gotool/collection"

type Map[K comparable, V any] collection.Map[K, V]

func NewMap[K comparable, V any]() Map[K, V] {
	return collection.NewMap[K, V]()
}
