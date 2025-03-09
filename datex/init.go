package datex

import (
	"time"

	"github.com/llyb120/yoya/cachex"
	"github.com/llyb120/yoya/syncx"
)

var compareHolder *syncx.Holder[*cachex.BaseCache[string, time.Time]]

func init() {
	compareHolder = syncx.NewHolder[*cachex.BaseCache[string, time.Time]](func() *cachex.BaseCache[string, time.Time] {
		return cachex.NewBaseCache[string, time.Time](cachex.OnceCacheOption{
			Expire:           30 * time.Second,
			CheckInterval:    0,
			DefaultKeyExpire: 0,
			Destroy: func() {
				// 自动清理
				compareHolder.Del()
			},
		})
	})
}
