package helper

import (
	"WenBeego/apps/common/global"
	"context"
	"encoding/json"
	"errors"
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
func RedisPut(key string, value interface{}, timeoutAfter int) error {
	ctx := context.Background()
	key, err := getKey(key)
	if err != nil {
		return err
	}

	switch v := value.(type) {
	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		_ = v
	default:
		valueStr, err := json.Marshal(value)
		if err != nil {
			return err
		}
		value = string(valueStr)
	}
	return global.Redis.Put(ctx, key, value, time.Duration(timeoutAfter)*time.Second)
}

// redis get value from redis
func RedisGet(key string) (string, error) {
	ctx := context.Background()
	key, err := getKey(key)
	if err != nil {
		return "", err
	}
	data, err := global.Redis.Get(ctx, key)
	if err != nil {
		return "", err
	} else if data == nil {
		return "", nil
	} else {
		if bytes, ok := data.([]byte); ok {
			return string(bytes), nil
		}
	}
	return "", errors.New("redis get error")
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
