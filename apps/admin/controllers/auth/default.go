package auth

import (
	// beego "github.com/beego/beego/v2/server/web"
	commonControllers "WenBeego/apps/common/controller"
)

type IndexController struct {
	// beego.Controller
	commonControllers.AdminBaseController
}

func (c *IndexController) Get() {
	c.Data["Website"] = "beego.vip"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "auth/index.tpl"
}
