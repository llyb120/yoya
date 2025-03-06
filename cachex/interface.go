package cachex

import "time"

type Cache[T any] interface {
	Get(key string) (value T, ok bool)
	Set(key string, value T)
	SetExpire(key string, value T, expire time.Duration)
	Del(key string)
}

type cacheItemWrapper[T any] struct {
	value     T
	expire    time.Time
	canExpire bool
}
