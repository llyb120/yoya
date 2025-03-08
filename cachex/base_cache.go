package cachex

import (
	"context"
	"time"

	"github.com/llyb120/gotool/internal/lockx"
)

// 一次性缓存，超过多久即会销毁

type BaseCache[K comparable, V any] struct {
	mu    lockx.Lock
	cache map[K]cacheItemWrapper[V]
	opts  OnceCacheOption
}

type OnceCacheOption struct {
	Expire           time.Duration
	DefaultKeyExpire time.Duration
	CheckInterval    time.Duration
	Destroy          func()
}

func NewBaseCache[K comparable, V any](opts OnceCacheOption) *BaseCache[K, V] {
	cache := &BaseCache[K, V]{
		opts:  opts,
		cache: make(map[K]cacheItemWrapper[V]),
	}
	go cache.start()
	return cache
}

func (c *BaseCache[K, V]) start() {
	if c.opts.Expire > 0 && c.opts.Destroy != nil {
		defer c.opts.Destroy()
	}

	var ctx context.Context
	var cancel context.CancelFunc
	if c.opts.Expire > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), c.opts.Expire)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}
	defer cancel()

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

func (c *BaseCache[K, V]) Set(key K, value V) {
	c.SetExpire(key, value, c.opts.DefaultKeyExpire)
}

func (c *BaseCache[K, V]) SetExpire(key K, value V, expire time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = cacheItemWrapper[V]{
		value:     value,
		expire:    time.Now().Add(expire),
		canExpire: expire > 0,
	}
}

func (c *BaseCache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.cache[key]
	if !ok {
		return item.value, false
	}
	return item.value, true
}

func (c *BaseCache[K, V]) Del(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

func (c *BaseCache[K, V]) GetOrSetFunc(key K, fn func() V) V {
	value, ok := c.Get(key)
	if !ok {
		if value, ok = c.Get(key); ok {
			return value
		}
		c.mu.Lock()
		defer c.mu.Unlock()
		value = fn()
		c.cache[key] = cacheItemWrapper[V]{
			value:  value,
			expire: time.Now().Add(c.opts.DefaultKeyExpire),
		}
		return value
	}
	return value
}
