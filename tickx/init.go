package datex

import (
	"time"

	"github.com/llyb120/yoya/cachex"
	"github.com/llyb120/yoya/syncx"
)

var compareHolder syncx.Holder[cachex.Cache[string, time.Time]]

func init() {
	compareHolder.InitFunc = func() cachex.Cache[string, time.Time] {
		return cachex.NewBaseCache[string, time.Time](cachex.CacheOption{
			Expire:           30 * time.Second,
			CheckInterval:    0,
			DefaultKeyExpire: 0,
			Destroy: func() {
				// 自动清理
				compareHolder.Del()
			},
		})
	}
}
