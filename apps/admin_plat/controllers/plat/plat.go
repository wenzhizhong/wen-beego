package plat

import (
	servicePlat "WenBeego/apps/admin_plat/services/plat"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto"
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
func (c *PlatController) GetUserUnitList() {
	userId := c.Ctx.Input.GetData("userId")
	data, err := c.PlatService.GetUserUnitList(userId.(string))
	if err != nil {
		c.Data["json"] = helper.Response{Code: 500, Message: err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = helper.Response{Code: 200, Message: "获取组织单位成功", Data: data}
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
	userId := c.Ctx.Input.GetData("userId")
	ChangeUnitDto, err := helper.GetReqBody[dto.ChangeUnitDto](c.Ctx)
	if err != nil {
		c.Data["json"] = helper.Response{Code: 0, Message: err.Error()}
		c.ServeJSON()
		return
	}
	data, err := c.PlatService.ChangeUnit("admin_plat", userId.(string), ChangeUnitDto.Id)
	if err != nil {
		c.Data["json"] = helper.Response{Code: 500, Message: err.Error()}
		c.ServeJSON()
		return
	}

	c.Data["json"] = helper.Response{Code: 200, Message: "切换组织单位成功", Data: data}
	c.ServeJSON()
}
