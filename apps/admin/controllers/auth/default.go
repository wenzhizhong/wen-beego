package auth

import (
	// beego "github.com/beego/beego/v2/server/web"
	commonControllers "WenBeego/apps/common/controller"

	_ "WenBeego/apps/common/helper"
)

type IndexController struct {
	// beego.Controller
	commonControllers.AdminBaseController
}

// ShowAccount godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  helper.Response
// @Router       /accounts/{id} [get]
func (c *IndexController) Get() {
	c.Data["Website"] = "beego.vip"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "auth/index.tpl"
}
