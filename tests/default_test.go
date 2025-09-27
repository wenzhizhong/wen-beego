package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/middleware"
	_ "WenBeego/routers"

	_ "github.com/beego/beego/v2/core/config/yaml"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	fmt.Println("init ....。 ")
	_, file, _, _ := runtime.Caller(0)
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)

	_ = beego.LoadAppConfig("yaml", "conf/app.yaml")
	// beego.TestBeegoInit(apppath)
	path := filepath.Join(apppath, "conf", "app.yaml")
	os.Chdir(apppath)
	// 注册自己资源服务
	initSourceService()
	beego.InitBeegoBeforeTest(path)
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
