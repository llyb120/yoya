package cachex

import "time"

type Cache[K comparable, V any] interface {
	Get(key K) (value V, ok bool)
	Gets(keys ...K) []V
	Set(key K, value V)
	SetExpire(key K, value V, expire time.Duration)
	Del(key ...K)
	Clear()
	Destroy()
	GetOrSetFunc(key K, fn func() V) V
}

type cacheItemWrapper[T any] struct {
	value     T
	expire    time.Time
	canExpire bool
}
