package captcha

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"context"
	"time"

	"github.com/mojocn/base64Captcha"
)

type Base64CaptchaRedisStore struct {
	Expiration time.Duration
}

var _ base64Captcha.Store = (*Base64CaptchaRedisStore)(nil)

func (r *Base64CaptchaRedisStore) Set(id string, value string) error {
	// return global.RedisCache.Put(context.Background(), getCacheKey(id), value, r.Expiration)
	return global.RedisCache.Set(context.Background(), getCacheKey(id), value, r.Expiration).Err()
}

func (r *Base64CaptchaRedisStore) Get(id string, clear bool) string {
	// val, err := global.RedisCache.Get(context.Background(), getCacheKey(id))
	val, err := global.RedisCache.Get(context.Background(), getCacheKey(id)).Result()
	if err != nil {
		global.Log.Error("redis get error:", err)
		return ""
	}

	if clear {
		global.RedisCache.Del(context.Background(), getCacheKey(id))
	}

	return val
}

func (r *Base64CaptchaRedisStore) Verify(id, answer string, clear bool) bool {
	val := r.Get(id, false)
	if val == "" {
		return false
	}

	if clear {
		global.RedisCache.Del(context.Background(), getCacheKey(id))
	}

	return val == answer
}
func getCacheKey(id string) string {
	key, _ := helper.GetCustomRedisKey("captcha:" + id)
	return key
}
