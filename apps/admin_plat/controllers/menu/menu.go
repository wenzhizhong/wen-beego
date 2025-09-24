package menu

import (
	serviceMenu "WenBeego/apps/admin_plat/services/menu"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/helper"
)

type MenuController struct {
	commonControllers.AdminBaseController
	MenuService serviceMenu.MenuService
}

// 获取异步路由
// @Summary 获取异步路由
// @Description 获取异步路由
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} helper.Response "返回结果"
// @Router /admin_plat/menu/get-async-routes [get]
// @Security ApiKeyAuth
func (c *MenuController) GetAsyncRoutes() {
	userId := c.Ctx.Input.GetData("userId")
	unitId := c.Ctx.Input.GetData("unitId")
	data, err := c.MenuService.GetAsyncRoutes(c.ModuleName, unitId.(string), userId.(string))
	if err != nil {
		c.Data["json"] = helper.Response{Code: 500, Message: err.Error()}
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response{Code: 200, Message: "success", Data: data}
	c.ServeJSON()
}
