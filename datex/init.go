package datex

import (
	"time"

	"github.com/llyb120/gotool/cachex"
	"github.com/llyb120/gotool/syncx"
	"github.com/petermattis/goid"
)

var compareHolder *syncx.Holder[int64, *cachex.BaseCache[string, time.Time]]

func init() {
	compareHolder = syncx.NewHolder[int64, *cachex.BaseCache[string, time.Time]](func() *cachex.BaseCache[string, time.Time] {
		return cachex.NewBaseCache[string, time.Time](cachex.OnceCacheOption{
			Expire:           30 * time.Second,
			CheckInterval:    0,
			DefaultKeyExpire: 0,
			Destroy: func() {
				// 自动清理
				compareHolder.Del(goid.Get())
			},
		})
	})
}
