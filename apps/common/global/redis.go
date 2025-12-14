package global

import (
	// "github.com/beego/beego/v2/client/cache"

	"github.com/redis/go-redis/v9"
)

// var RedisCache cache.Cache
var RedisCache *redis.Client
