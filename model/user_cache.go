package model

import (
	"context"
	"time"

	"github.com/code-art/gin-im/util"
)

// 设置在线用户到Redis缓存
func SetUserOnlineInfo(key string, val []byte, ttl time.Duration) {
	ctx := context.Background()
	util.Redigo.Set(ctx, key, val, ttl)
}
