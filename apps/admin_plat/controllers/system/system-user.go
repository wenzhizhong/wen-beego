package system

// 系统管理-内部用户管理
import (
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/helper"

	systemService "WenBeego/apps/admin_plat/services/system"
	commonControllers "WenBeego/apps/common/controller"
)

type UserController struct {
	commonControllers.AdminBaseController
	UserService systemService.UserService
}

// 系统管理-获取用户列表
// @Summary 系统管理-获取用户列表
// @Description 系统管理-获取用户列表
// @Tags 系统管理-用户管理
// @Accept application/json
// @Produce application/json
// @Param data body dto.GetUserListDto true "请求参数"
// @Success 200 {object} dto.GetUserListDto "返回结果"
// @Router /admin_plat/system-user/get [get]

func (c *UserController) GetUserList() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	reqDataListDto, err2 := helper.GetReqDataListDto(&c.Controller)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	if err2 != nil {
		c.Data["json"] = helper.Response(500, err2.Error(), nil)
		c.ServeJSON()
		return
	}

	userDto := page_dto.SystemUserListReqDto{}
	userDto.BaseParamDto = baseParamDto
	userDto.ReqDataListDto = reqDataListDto
	userDto.SelectUnitIds = c.GetStrings("selectUnitIds")

	data, err := c.UserService.GetUserList(userDto)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

func (c *UserController) Add() {
}
func (c *UserController) Edit() {
}
func (c *UserController) Del() {
}
