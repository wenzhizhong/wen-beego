package auth

//身份认证-用户组织单位
import (
	platService "WenBeego/apps/admin_plat/services/auth"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto_vo/auth_dto"
	"WenBeego/apps/common/helper"
	"fmt"
)

type PlatController struct {
	commonControllers.AdminBaseController
	PlatService platService.Plat
}

// 获取用户组织单位
// @Summary 获取用户组织单位
// @Description 获取用户组织单位
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response "返回结果"
// @Router /admin_plat/auth-plat/get-user-unit [get]
// @Security ApiKeyAuth
func (c *PlatController) GetUserUnitList() {
	userId := c.Ctx.Input.GetData("userId")
	host := c.Ctx.Request.Host
	data, err := c.PlatService.GetUserUnitList(userId.(string), host)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}

	c.Data["json"] = helper.Response(200, "获取组织单位成功", data)
	c.ServeJSON()
}

// 切换组织
// @Summary 切换组织
// @Description 切换组织
// @Tags admin
// @Accept json
// @Produce json
// @Success 200 {object} dto.Response "返回结果"
// @Router /admin_plat/auth-plat/change-unit [get]
// @Security ApiKeyAuth
func (c *PlatController) ChangeUnit() {
	fmt.Println("ChangeUnit", c.Ctx)
	userId := c.Ctx.Input.GetData("userId")
	ChangeUnitDto, err := helper.GetReqBody[auth_dto.ChangeUnitDto](c.Ctx)
	if err != nil {
		c.Data["json"] = helper.Response(0, err.Error(), nil)
		c.ServeJSON()
		return
	}
	data, err := c.PlatService.ChangeUnit(c.ModuleName, userId.(string), ChangeUnitDto)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}

	c.Data["json"] = helper.Response(200, "切换组织单位成功", data)
	c.ServeJSON()
}
