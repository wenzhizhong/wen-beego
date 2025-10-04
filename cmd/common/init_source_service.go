package common

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/middleware"
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
	err = middleware.InitDb()
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
}

func InitMqClient() {
	client := &middleware.MqClient{}
	client.Init()
}
