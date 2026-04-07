package common

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/middleware"
	"WenBeego/apps/common/middleware/business_store"
	"WenBeego/apps/common/middleware/crontab"
	"WenBeego/apps/common/middleware/database"
	"WenBeego/apps/common/middleware/mq"
)

// 初始化公共资源数据
func InitCommonSource(logType string) {
	// 输出终端ASCII艺术字符logo
	helper.Output_ASCII_ArtisticCharacters()

	// 注册系统日志中间件
	err := middleware.InitLog(logType)
	if err != nil {
		global.Log.Error(err.Error())
		panic(err)
	}
	// 注册数据库中间件
	err = database.InitDb()
	if err != nil {
		global.Log.Error(err.Error())
		panic(err)
	}
	// 注册缓存中间件
	err = middleware.InitRedis()
	if err != nil {
		global.Log.Error(err.Error())
		panic(err)
	}

	// 清空用户权限认证缓存
	business_store.ClearAumid()
}

// 初始化Mq client
func InitMqClient() {
	client := &mq.MqClient{}
	client.Init()
}

// 初始化定时任务
func InitCrontabTask() *crontab.CronManager {
	return crontab.GetCronManager()
}
