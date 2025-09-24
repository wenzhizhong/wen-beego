package helper

import (
	"WenBeego/apps/common/global"
	"context"
	"time"
)

// redis key
func getKey(key string) (string, error) {
	runmode, err := AppRunmode()
	if err != nil {
		return "", err
	}
	return runmode + ":" + key, nil
}

// redis put value to redis
func RedisPut(key string, value interface{}, timeout time.Duration) error {
	ctx := context.Background()
	key, err := getKey(key)
	if err != nil {
		return err
	}
	return global.Redis.Put(ctx, key, value, timeout*time.Second)
}

// redis get value from redis
func RedisGet(key string) (interface{}, error) {
	ctx := context.Background()
	key, err := getKey(key)
	if err != nil {
		return nil, err
	}
	return global.Redis.Get(ctx, key)
}

// redis delete value from redis
func RedisDel(key string) error {
	ctx := context.Background()
	key, err := getKey(key)
	if err != nil {
		return err
	}
	return global.Redis.Delete(ctx, key)
}
