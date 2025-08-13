package middleware

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"encoding/json"
	"fmt"

	_ "github.com/beego/beego/v2/client/cache"
	"github.com/beego/beego/v2/client/cache/redis"
)

type redisConfig struct {
	DbNum           string `json:"dbNum"`
	SkipEmptyPrefix string `json:"skipEmptyPrefix"`
	Key             string `json:"key"`
	Conn            string `json:"conn"`
	MaxIdle         string `json:"maxIdle"`
	TimeoutStr      string `json:"timeout"`
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
		redisConfigObj.Password = tmpConfig["password"]
		redisConfigObj.MaxIdle = tmpConfig["maxIdle"]
		redisConfigObj.TimeoutStr = tmpConfig["timeout"]
		redisConfigObj.DbNum = tmpConfig["dbNum"]
		redisConfigObj.Key = tmpConfig["key"]
		redisConfigObj.SkipEmptyPrefix = tmpConfig["skipEmptyPrefix"]

		redisConfig, err := json.Marshal(redisConfigObj)
		if err != nil {
			return err
		}

		redisCache := redis.NewRedisCache()
		err2 := redisCache.StartAndGC(string(redisConfig))
		if err2 != nil {
			return err2
		}
		global.Redis = redisCache
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
