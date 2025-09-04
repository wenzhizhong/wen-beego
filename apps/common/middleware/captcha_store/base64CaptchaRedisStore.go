package captcha_store

import (
	"WenBeego/apps/common/global"
	"context"
	"time"

	"github.com/mojocn/base64Captcha"
)

type Base64CaptchaRedisStore struct {
	Expiration time.Duration
}

var _ base64Captcha.Store = (*Base64CaptchaRedisStore)(nil)

func (r *Base64CaptchaRedisStore) Set(id string, value string) error {
	return global.Redis.Put(context.Background(), getCacheKey(id), value, r.Expiration)
}

func (r *Base64CaptchaRedisStore) Get(id string, clear bool) string {
	val, err := global.Redis.Get(context.Background(), getCacheKey(id))
	if err != nil {
		global.Log.Error("redis get error:", err)
		return ""
	}
	if val == nil {
		return ""
	}

	// 先断言为 []byte，再转换为 string
	byteVal, ok := val.([]byte)
	if !ok {
		global.Log.Error("failed to assert value to []byte")
		return ""
	}
	if clear {
		global.Redis.Delete(context.Background(), getCacheKey(id))
	}

	return string(byteVal)
}

func (r *Base64CaptchaRedisStore) Verify(id, answer string, clear bool) bool {
	val := r.Get(id, false)
	if val == "" {
		return false
	}

	if clear {
		global.Redis.Delete(context.Background(), getCacheKey(id))
	}

	return val == answer
}
func getCacheKey(id string) string {
	return "captcha:" + id
}
