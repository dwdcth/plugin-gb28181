package msgstore

import (
	"github.com/ReneKroon/ttlcache/v2"
	"time"
)

var store = ttlcache.NewCache()

func StoreMsg(id string, msg interface{}) {
	store.SetWithTTL(id, msg, 10*time.Second)
}

// 获取消息 超时ms
func GetMsg(id string, timeOut int) (interface{}, error) {
	if timeOut == 0 {
		return store.Get(id)
	}
	ticker := time.NewTicker(time.Microsecond * 200)
	for {
		select {
		case <-ticker.C:
			if value, err := store.Get(id); err == nil {
				ticker.Stop()
				return value, err
			}
		case <-time.After(time.Duration(timeOut) * time.Millisecond):
			ticker.Stop()
			return store.Get(id)
		}
	}
}
