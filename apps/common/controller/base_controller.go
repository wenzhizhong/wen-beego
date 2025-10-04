package controller

import (
	"WenBeego/apps/common/helper"

	"github.com/beego/beego/v2/server/web"
)

type BaseController struct {
	web.Controller
	ModuleName string
}

func (c *BaseController) Prepare() {
	// 从路由中解析模块名
	c.ModuleName = helper.ParseModuleFromRoute(c.Ctx) // 使用正确的context类型
	// 设置视图路径
	c.TplExt = "html"
	c.ViewPath = "../../apps/" + c.ModuleName + "/views"
}
