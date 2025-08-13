package home

import (
	commonControllers "WenBeego/apps/common/controller"
)

type IndexController struct {
	commonControllers.IndexBaseController
}

func (c *IndexController) Get() {
	c.Data["Website"] = "beego.vip1"
	c.Data["Email"] = "astaxie@gmail.com"
	c.Data["ModuleName"] = c.BaseController.ModuleName

	c.TplName = "home/index.html"
	c.Render()
}

func (c *IndexController) Post() {
	c.Data["json"] = map[string]string{
		"Website":    "beego.vip1",
		"Email":      "astaxie@gmail.com",
		"ModuleName": c.BaseController.ModuleName,
	}
	c.ServeJSON()
}
