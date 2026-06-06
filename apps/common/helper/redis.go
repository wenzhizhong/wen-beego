package helper

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/global/app_error"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

// redis key
func GetCustomRedisKey(key string) (string, error) {
	appName, err1 := AppName()
	runmode, err2 := AppRunmode()
	if err1 != nil || err2 != nil {
		err := Ternary(err1 != nil, err1, err2)
		return "", err
	}
	return appName + ":" + runmode + ":" + key, nil
}

// redis put value to redis
func RedisPut(key string, value interface{}, timeoutAfter int) error {
	ctx := context.Background()
	key, err := GetCustomRedisKey(key)
	if err != nil {
		return err
	}

	switch v := value.(type) {
	case string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool:
		_ = v
	default:
		valueStr, err := json.Marshal(value)
		if err != nil {
			return app_error.NewHelperError(err)
		}
		value = string(valueStr)
	}
	// return global.RedisCache.Put(ctx, key, value, time.Duration(timeoutAfter)*time.Second)
	return global.RedisCache.Set(ctx, key, value, time.Duration(timeoutAfter)*time.Second).Err()
}

// redis get value from redis
func RedisGet(key string) (string, error) {
	ctx := context.Background()
	key, err := GetCustomRedisKey(key)
	if err != nil {
		return "", err
	}
	// data, err := global.RedisCache.Get(ctx, key).Result()
	// if err != nil {
	// 	return "", err
	// } else if data == nil {
	// 	return "", nil
	// } else {
	// 	if bytes, ok := data.([]byte); ok {
	// 		return string(bytes), nil
	// 	}
	// }
	// return "", errors.New("redis get error")
	var data string = ""
	data, err = global.RedisCache.Get(ctx, key).Result()
	err = Ternary(err == nil || err == redis.Nil, nil, err)
	if err != nil && err != redis.Nil {
		return "", app_error.NewHelperError(err)
	}
	return data, err
}

// redis delete value from redis
func RedisDel(key string) error {
	ctx := context.Background()
	key, err := GetCustomRedisKey(key)
	if err != nil {
		return err
	}
	// return global.RedisCache.Delete(ctx, key)
	err = global.RedisCache.Del(ctx, key).Err()
	err = Ternary(err == nil || err == redis.Nil, nil, err)
	if err != nil {
		return app_error.NewHelperError(err)
	}
	return nil
}
