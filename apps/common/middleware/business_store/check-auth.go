package business_store

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"context"
	"strconv"
	"strings"

	"github.com/beego/beego/v2/client/cache/redis"
)

const aumid_expired = 2 * 60

const AUMID_UPS = "AUMID_UPS"
const AUMID_US = "AUMID_US"
const AUMID_URS = "AUMID_URS"
const AUMID_URP = "AUMID_URP"

var ALL_AUMID_PRIFIXS = []string{
	AUMID_UPS,
	AUMID_US,
	AUMID_URS,
	AUMID_URP,
}

// get key
func getAumidUpsKey(userId string, unitId string) string {
	return AUMID_UPS + ":" + helper.Md5(userId+unitId)
}
func getAumidUsKey(userId string, unitId string) string {
	return AUMID_US + ":" + helper.Md5(userId+unitId)
}

func getAumidUrsKey(unitUserId string) string {
	return AUMID_URS + ":" + unitUserId
}

func getAumidUrpKey(moduleName, unitUserId, unitId string) string {
	return AUMID_URP + ":" + helper.Md5(moduleName+unitUserId+unitId)
}

// do get value
func doGet(redisKey string, defaultVal int) (status int, err error) {
	exits := ""
	exits, err = helper.RedisGet(redisKey)
	if err == nil && exits != "" {
		status, _ = strconv.Atoi(exits)
	} else {
		status = defaultVal
	}
	return
}

// 验证用户资料状态缓存
func GetAumidUps(userId string, unitId string, defaultVal int) (status int, err error) {
	redisKey := getAumidUpsKey(userId, unitId)
	status, err = doGet(redisKey, defaultVal)
	return
}
func SetAumidUps(userId string, unitId string, status int) error {
	redisKey := getAumidUpsKey(userId, unitId)
	err := helper.RedisPut(redisKey, status, aumid_expired)
	return err
}

// 验证组织单位状态缓存
func GetAumidUs(userId string, unitId string, defaultVal int) (status int, err error) {
	redisKey := getAumidUsKey(userId, unitId)
	status, err = doGet(redisKey, defaultVal)
	return
}

func SetAumidUs(userId string, unitId string, status int) error {
	redisKey := getAumidUsKey(userId, unitId)
	err := helper.RedisPut(redisKey, status, aumid_expired)
	return err
}

// 验证用户角色状态缓存
func GetAumidUrs(unitUserId string, defaultVal int) (status int, err error) {
	redisKey := getAumidUrsKey(unitUserId)
	status, err = doGet(redisKey, defaultVal)
	return
}
func SetAumidUrs(unitUserId string, status int) error {
	redisKey := getAumidUrsKey(unitUserId)
	err := helper.RedisPut(redisKey, status, aumid_expired)
	return err
}

// api权限缓存
func GetAumidUrp(moduleName, unitUserId, unitId string) (exits string, err error) {
	redisKey := getAumidUrpKey(moduleName, unitUserId, unitId)
	exits, err = helper.RedisGet(redisKey)
	return
}
func SetAumidUrp(moduleName, unitUserId, unitId, value string) error {
	redisKey := getAumidUrpKey(moduleName, unitUserId, unitId)
	err := helper.RedisPut(redisKey, value, aumid_expired)
	return err
}

// 清空认证缓存
func ClearAumid() error {
	if redisCache, ok := global.Redis.(*redis.Cache); ok {
		for _, prefix := range ALL_AUMID_PRIFIXS {
			key, _ := helper.GetCustomRedisKey(prefix)
			keys, err := redisCache.Scan(redis.DefaultKey + ":" + key + ":*")
			if err != nil {
				global.Log.Error("redisCache.Scan error:", err)
				continue
			}
			for _, k := range keys {
				k = strings.Replace(k, redis.DefaultKey+":", "", -1)
				err = redisCache.Delete(context.Background(), k)
				if err != nil {
					global.Log.Error("redisCache.Delete error:", err)
					continue
				}
			}
		}
	}
	return nil
}
