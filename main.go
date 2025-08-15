package main

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/middleware"
	_ "WenBeego/routers"
	"fmt"
	"strconv"

	_ "github.com/beego/beego/v2/core/config/yaml"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// beego自定义配置
	_ = beego.LoadAppConfig("yaml", "conf/app.yaml")
	beego.InsertFilter("/*", beego.BeforeRouter, new(middleware.AccessMiddleware).LimitTimes())
	beego.AddViewPath("apps/index/views")
	beego.AddViewPath("apps/admin/views")
	fmt.Println("beego.BConfig.RunMode:", beego.BConfig.RunMode)
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "apps/swagger"
	}

	// 注册自己资源服务
	initSourceService()

	// 启动服务
	httpport, _ := beego.AppConfig.DIY("httpport")
	beego.Run(":" + strconv.Itoa(httpport.(int)))
}

func initSourceService() {
	// 输出终端ASCII艺术字符logo
	helper.Output_ASCII_ArtisticCharacters()

	// 注册系统日志中间件
	err := middleware.InitLog()
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
