package cachex

import (
	"context"
	"sync"
	"time"
)

// 一次性缓存，超过多久即会销毁

type baseCache[K comparable, V any] struct {
	mu     sync.RWMutex
	cache  map[K]cacheItemWrapper[V]
	opts   CacheOption
	cancel context.CancelFunc
}

type CacheOption struct {
	Expire           time.Duration
	DefaultKeyExpire time.Duration
	CheckInterval    time.Duration
	Destroy          func()
}

func NewBaseCache[K comparable, V any](opts CacheOption) *baseCache[K, V] {
	cache := &baseCache[K, V]{
		opts:  opts,
		cache: make(map[K]cacheItemWrapper[V]),
	}
	go cache.start()
	return cache
}

func (c *baseCache[K, V]) start() {
	if c.opts.Expire > 0 && c.opts.Destroy != nil {
		defer c.opts.Destroy()
	}

	var ctx context.Context
	if c.opts.Expire > 0 {
		ctx, c.cancel = context.WithTimeout(context.Background(), c.opts.Expire)
	} else {
		ctx, c.cancel = context.WithCancel(context.Background())
	}
	defer c.cancel()

	if c.opts.CheckInterval > 0 {
		// 小于等于0的时候永不过期
		go func() {
			ticker := time.NewTicker(c.opts.CheckInterval)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					func() {
						c.mu.Lock()
						// 执行检查操作
						defer c.mu.Unlock()
						mp := make(map[K]cacheItemWrapper[V])
						now := time.Now()
						for key, item := range c.cache {
							if item.canExpire && item.expire.After(now) {
								mp[key] = item
							}
						}
						c.cache = mp
					}()
				}
			}
		}()
	}

	<-ctx.Done()
}

func (c *baseCache[K, V]) Set(key K, value V) {
	c.SetExpire(key, value, c.opts.DefaultKeyExpire)
}

func (c *baseCache[K, V]) SetExpire(key K, value V, expire time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheItemWrapper[V]{
		value:     value,
		expire:    time.Now().Add(expire),
		canExpire: expire > 0,
	}
}

func (c *baseCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.get(key)
}

func (c *baseCache[K, V]) Gets(keys ...K) []V {
	c.mu.RLock()
	defer c.mu.RUnlock()
	values := make([]V, len(keys))
	for i, key := range keys {
		value, ok := c.get(key)
		if ok {
			values[i] = value
		}
	}
	return values
}

func (c *baseCache[K, V]) get(key K) (V, bool) {
	item, ok := c.cache[key]
	if !ok {
		return item.value, false
	}
	return item.value, true
}

func (c *baseCache[K, V]) Del(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

func (c *baseCache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[K]cacheItemWrapper[V])
}

func (c *baseCache[K, V]) Destroy() {
	c.mu.Lock()
	c.cancel()
	c.mu.Unlock()

	// 解锁，防止出问题
	if c.opts.Destroy != nil {
		c.opts.Destroy()
	}
}

func (c *baseCache[K, V]) GetOrSetFunc(key K, fn func() V) V {
	value, ok := c.Get(key)
	if !ok {
		c.mu.Lock()
		defer c.mu.Unlock()
		if value, ok = c.get(key); ok {
			return value
		}
		value = fn()
		c.cache[key] = cacheItemWrapper[V]{
			value:  value,
			expire: time.Now().Add(c.opts.DefaultKeyExpire),
		}
		return value
	}
	return value
}
