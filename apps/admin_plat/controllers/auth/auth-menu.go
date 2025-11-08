package auth

//身份认证-用户系统菜单/系统权限
import (
	menuService "WenBeego/apps/admin_plat/services/auth"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/helper"
)

type MenuController struct {
	commonControllers.AdminBaseController
	MenuService menuService.MenuService
}

// 获取异步路由
// @Summary 获取异步路由
// @Description 获取异步路由
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response "返回结果"
// @Router /admin_plat/auth-menu/get-async-routes [get]
// @Security ApiKeyAuth
func (c *MenuController) GetAsyncRoutes() {
	// userId := c.Ctx.Input.GetData("userId")
	unitId := c.Ctx.Input.GetData("unitId")
	unitUserId := c.Ctx.Input.GetData("unitUserId")
	data, err := c.MenuService.GetAsyncRoutes(c.ModuleName, unitId.(string), unitUserId.(string))
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}
