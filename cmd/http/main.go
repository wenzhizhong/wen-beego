package main

import (
	"WenBeego/apps/common/global"
	cmdCommon "WenBeego/cmd/common"
	_ "WenBeego/routers"
	"os"
	"strconv"

	_ "github.com/beego/beego/v2/core/config/yaml"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/grace"
)

func main() {
	// 注册自己资源服务
	cmdCommon.RunBefore()
	cmdCommon.InitCommonSource("")
	cmdCommon.InitMqClient()
	cronManager := cmdCommon.InitCrontabTask()
	cronManager.Start()
	defer cronManager.Stop()

	// // 启动服务
	// httpport, _ := beego.AppConfig.DIY("httpport")
	// beego.Run(":" + strconv.Itoa(httpport.(int)))
	// fmt.Println("---->>>---")

	httpport, _ := beego.AppConfig.DIY("httpport")
	port := strconv.Itoa(httpport.(int))
	err := grace.ListenAndServe(":"+port, beego.BeeApp.Handlers)
	if err != nil {
		global.Log.Error("Server on %v stopped, err: %v", port, err)
	} else {
		global.Log.Info("Server on %v stopped", port)
	}
	os.Exit(0)
}
