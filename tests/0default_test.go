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
	apppath, _ := filepath.Abs(filepath.Dir(filepath.Join(file, ".."+string(filepath.Separator))))
	beego.TestBeegoInit(apppath)

	_ = beego.LoadAppConfig("yaml", "conf/app.yaml")
	// beego.TestBeegoInit(apppath)
	path := filepath.Join(apppath, "conf", "app.yaml")
	os.Chdir(apppath)
	// 注册自己资源服务
	cmdCommon.InitCommonSource("pathUnitTest")
	cmdCommon.InitMqClient()

	beego.InitBeegoBeforeTest(path)
}
