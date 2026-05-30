package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	_ "WenBeego/routers"

	cmdCommon "WenBeego/cmd/common"

	_ "github.com/beego/beego/v2/core/config/yaml"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	fmt.Println("init ....。 ")
	_, file, _, _ := runtime.Caller(0)
	appPath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(appPath)

	_ = beego.LoadAppConfig("yaml", "conf/app.yaml")
	// beego.TestBeegoInit(appPath)
	path := filepath.Join(appPath, "conf", "app.yaml")
	os.Chdir(appPath)
	// 注册自己资源服务
	cmdCommon.InitCommonSource("pathUnitTest")
	cmdCommon.InitMqClient()

	beego.InitBeegoBeforeTest(path)
}
