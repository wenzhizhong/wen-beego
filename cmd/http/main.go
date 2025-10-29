package main

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/middleware"
	cmdCommon "WenBeego/cmd/common"
	_ "WenBeego/routers"
	"fmt"
	"strconv"

	_ "github.com/beego/beego/v2/core/config/yaml"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	cmdCommon.RunBefore()

	// beego自定义配置
	_ = beego.LoadAppConfig("yaml", global.ConfigDir+"/app.yaml")
	beego.InsertFilter("/*", beego.BeforeRouter, new(middleware.AccessMiddleware).LimitTimes())
	beego.AddViewPath(global.AppDir + "/index/views")
	beego.AddViewPath(global.AppDir + "/admin_plat/views")
	beego.BConfig.WebConfig.StaticDir["/static"] = global.StaticDir
	fmt.Println("beego.BConfig.RunMode:", beego.BConfig.RunMode)
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = global.AppDir + "/swagger"
	}

	// 注册自己资源服务
	cmdCommon.InitCommonSource("")
	cmdCommon.InitMqClient()

	// 启动服务
	httpport, _ := beego.AppConfig.DIY("httpport")
	beego.Run(":" + strconv.Itoa(httpport.(int)))
}
