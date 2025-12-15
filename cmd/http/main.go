package main

import (
	cmdCommon "WenBeego/cmd/common"
	_ "WenBeego/routers"
	"fmt"
	"strconv"

	_ "github.com/beego/beego/v2/core/config/yaml"
	beego "github.com/beego/beego/v2/server/web"
)

func main() {
	// 注册自己资源服务
	cmdCommon.RunBefore()
	cmdCommon.InitCommonSource("")
	cmdCommon.InitMqClient()
	cronManager := cmdCommon.InitCrontabTask()
	cronManager.Start()
	defer cronManager.Stop()

	// 启动服务
	httpport, _ := beego.AppConfig.DIY("httpport")
	beego.Run(":" + strconv.Itoa(httpport.(int)))
	fmt.Println("---->>>---")
}
