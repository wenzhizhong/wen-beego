package helper

import (
	"WenBeego/apps/common/global"
	"context"
	"time"
)

// redis put value to redis
func RedisPut(key string, value interface{}, timeout time.Duration) error {
	ctx := context.Background()
	return global.Redis.Put(ctx, key, value, timeout*time.Second)
}

// redis get value from redis
func RedisGet(key string) (interface{}, error) {
	ctx := context.Background()
	return global.Redis.Get(ctx, key)
}

// redis delete value from redis
func RedisDel(key string) error {
	ctx := context.Background()
	return global.Redis.Delete(ctx, key)
}
