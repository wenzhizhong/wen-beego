package middleware

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"context"
	"fmt"
	"strconv"

	_ "github.com/beego/beego/v2/client/cache"
	// "github.com/beego/beego/v2/client/cache/redis"
	"github.com/redis/go-redis/v9"
)

type redisConfig struct {
	// DbNum           string `json:"dbNum"`
	DbNum           int    `json:"dbNum"`
	SkipEmptyPrefix string `json:"skipEmptyPrefix"`
	Key             string `json:"key"`
	Conn            string `json:"conn"`
	MaxIdle         int    `json:"maxIdle"`
	TimeoutStr      string `json:"timeout"`
	Username        string `json:"username"`
	Password        string `json:"password"`
}

func InitRedis() error {
	tmpConfig, err := getRedisConfig()
	if err != nil {
		return err
	}
	if tmpConfig["host"] != "" && tmpConfig["port"] != "" {
		fmt.Println("init redis...")

		var redisConfigObj redisConfig
		redisConfigObj.Conn = tmpConfig["host"] + ":" + tmpConfig["port"]
		redisConfigObj.Username = tmpConfig["username"]
		redisConfigObj.Password = tmpConfig["password"]
		redisConfigObj.MaxIdle, err = helper.String2Int(tmpConfig["maxIdle"])
		redisConfigObj.TimeoutStr = tmpConfig["timeout"]
		redisConfigObj.DbNum, err = strconv.Atoi(tmpConfig["dbNum"])
		redisConfigObj.Key = tmpConfig["key"]
		redisConfigObj.SkipEmptyPrefix = tmpConfig["skipEmptyPrefix"]
		if err != nil {
			return err
		}

		// redisConfig, err := json.Marshal(redisConfigObj)
		// if err != nil {
		// 	return err
		// }
		// redisCache := redis.NewRedisCache()
		// err2 := redisCache.StartAndGC(string(redisConfig))
		// if err2 != nil {
		// 	return err2
		// }
		// global.RedisCache = redisCache

		ctx := context.Background()

		rdb := redis.NewClient(&redis.Options{
			Addr:     redisConfigObj.Conn,
			Username: redisConfigObj.Username,
			Password: redisConfigObj.Password, // no password set
			DB:       redisConfigObj.DbNum,    // use default DB
			PoolSize: redisConfigObj.MaxIdle,
		})

		err = rdb.Ping(ctx).Err()
		if err != nil {
			return err
		}
		global.RedisCache = rdb

		fmt.Println("init redis done！")
	}
	return nil
}

func getRedisConfig() (map[string]string, error) {
	mapConfig, err := global.GetConfig("redis")
	if err != nil {
		return nil, err
	}
	redisConfig, err := helper.MapInterface2MapString(mapConfig)
	if err != nil {
		return nil, err
	}
	// if redisConfig["conn"] == "" {
	// 	return nil, errors.New("redis config error")
	// }

	return redisConfig, nil
}
