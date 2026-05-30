package common

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/middleware/blocker"
	"fmt"
	"os"
	"path/filepath"
	"runtime/debug"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func RunBefore() {
	// 全局目录
	initGlobalPath()
	// 自定义 beego
	customBeego()
}

func customBeego() {
	// beego自定义配置
	_ = beego.LoadAppConfig("yaml", global.ConfigDir+"/app.yaml")

	beego.InsertFilter("/*", beego.BeforeStatic, blocker.LimitTimes())
	beego.InsertFilter("/*", beego.BeforeRouter, blocker.LimitTimes())
	beego.AddViewPath(global.AppDir + "/index/views")
	beego.AddViewPath(global.AppDir + "/admin_plat/views")
	beego.BConfig.WebConfig.StaticDir["/static"] = global.StaticDir
	beego.BConfig.WebConfig.StaticDir["/uploads"] = global.UploadsDir
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = global.AppDir + "/swagger"
	}
	beego.BConfig.RecoverFunc = beegoRecoverFunc // 自定义替换掉默认defaultRecoverPanic方法

	fmt.Println("beego.BConfig.RunMode:", beego.BConfig.RunMode)
}

func initGlobalPath() {
	rootPath, _ := filepath.Abs("../../")
	appDir := filepath.Join(rootPath, "apps")
	configDir := filepath.Join(rootPath, "conf")
	staticDir := filepath.Join(rootPath, "static")
	routersDir := filepath.Join(rootPath, "routers")
	tempDir := filepath.Join(rootPath, "temp")
	uploadsDir := filepath.Join(rootPath, "uploads")

	global.RootPath = rootPath
	global.AppDir = appDir
	global.ConfigDir = configDir
	global.StaticDir = staticDir
	global.RoutersDir = routersDir
	global.TempDir = tempDir
	global.UploadsDir = uploadsDir

	createDir(tempDir)
	createDir(uploadsDir)

}
func createDir(path string) {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) || !fileInfo.IsDir() {
		os.MkdirAll(path, os.ModePerm)
	}
}

func beegoRecoverFunc(ctx *context.Context, config *beego.Config) {
	if err := recover(); err != nil {
		// 记录 panic 到日志
		runMode, _ := helper.AppRunmode()
		stack := string(debug.Stack())
		global.Log.Error("PANIC: %v\nStack: %s\n\n", err, stack)

		// 返回错误响应。确保不会重复写入响应
		if !ctx.ResponseWriter.Started {
			errMsg := helper.Ternary(runMode == "dev", fmt.Sprintf("%v", stack), "Internal Server Error")
			ctx.Output.SetStatus(500)
			ctx.Output.JSON(map[string]interface{}{
				"error": errMsg,
				"code":  500,
			}, false, false)
		}
	}
}
