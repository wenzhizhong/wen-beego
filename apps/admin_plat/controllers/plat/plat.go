package plat

import (
	servicePlat "WenBeego/apps/admin_plat/services/plat"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/helper"
	"fmt"
)

type PlatController struct {
	commonControllers.AdminBaseController
	PlatService servicePlat.Plat
}

// 获取用户组织单位
// @Summary 获取用户组织单位
// @Description 获取用户组织单位
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} helper.Response "返回结果"
// @Router /admin_plat/plat/get-user-unit [get]
// @Security ApiKeyAuth
func (c *PlatController) GetUserUnit() {
	fmt.Println("GetUserUnit", c.Ctx)
	userId := c.Ctx.Input.GetData("userId")
	data, err := c.PlatService.GetUserUnit(userId.(string))
	if err != nil {
		c.Data["json"] = helper.Response{Code: 500, Message: err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = helper.Response{Code: 200, Message: "登录成功", Data: data}
	c.ServeJSON()
}

// 切换组织
// @Summary 切换组织
// @Description 切换组织
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} helper.Response "返回结果"
// @Router /admin_plat/plat/change-unit [get]
// @Security ApiKeyAuth
func (c *PlatController) ChangeUnit() {
	fmt.Println("ChangeUnit", c.Ctx)

	c.Data["json"] = helper.Response{Code: 200, Message: "登录成功"}
	c.ServeJSON()
}
