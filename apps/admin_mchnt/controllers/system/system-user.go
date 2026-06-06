package system

// 系统管理-内部用户管理
import (
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/dto_vo/user_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/common/helper"
	"strings"

	systemService "WenBeego/apps/admin_mchnt/services/system"
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
// @Param data selectUnitIds string false "单位ID列表"
// @Param data deptId string false "部门ID列表"
// @Param data roleId string false "角色ID列表"
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Router /admin_mchnt/system-user/get [get]

func (c *UserController) GetUserList() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	reqDataListDto, err2 := helper.GetReqDataListDto(&c.Controller)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	if err2 != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err2.Error(), nil)
		c.ServeJSON()
		return
	}

	userDto := page_dto.SystemUserListReqDto{}
	userDto.BaseParamDto = baseParamDto
	userDto.ReqDataListDto = reqDataListDto
	userDto.SelectUnitIds = helper.Ternary(c.GetString("selectUnitIds") != "", strings.Split(c.GetString("selectUnitIds"), ","), []string{baseParamDto.UnitId})
	userDto.DeptIds = helper.Ternary(c.GetString("deptId") != "", strings.Split(c.GetString("deptId"), ","), make([]string, 0))
	userDto.RoleIds = helper.Ternary(c.GetString("roleId") != "", strings.Split(c.GetString("roleId"), ","), make([]string, 0))

	data, err := c.UserService.GetUserList(userDto)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// 新增组织单位用户
// @Summary 新增组织单位用户
// @Description 新增组织单位用户
// @Tags 组织单位用户
// @Accept  json
// @Param   userDto body user_dto.UnitUserSaveDto true "请求参数"
// @Success 200 {object} dto.Response  "返回结果"
// @Router /admin_mchnt/system-user/add [post]
func (c *UserController) Add() {
	c.save("add")
}

// 编辑组织单位用户
// @Summary 编辑组织单位用户
// @Description 编辑组织单位用户
// @Tags 组织单位用户
// @Accept  json
// @Param   userDto body user_dto.UnitUserSaveDto true "请求参数"
// @Success 200 {object} dto.Response "返回结果"
// @Router /admin_mchnt/system-user/edit [post]
func (c *UserController) Edit() {
	c.save("edit")
}
func (c *UserController) save(optType string) {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	unitUserSaveDto, err0 := helper.GetReqBody[user_dto.UnitUserSaveDto](c.Ctx)
	userDto, err1 := helper.GetReqBody[user_dto.UserDto](c.Ctx)
	userProfileDto, err2 := helper.GetReqBody[user_dto.UserProfileDto](c.Ctx)
	unitUserDto, err3 := helper.GetReqBody[user_dto.UnitUserDto](c.Ctx)
	unitUserProfileDto, err4 := helper.GetReqBody[user_dto.UnitUserProfileDto](c.Ctx)

	unitUserSaveDto.UserDto = userDto
	unitUserSaveDto.UserProfileDto = userProfileDto
	unitUserSaveDto.UnitUserDto = unitUserDto
	unitUserSaveDto.UnitUserProfileDto = unitUserProfileDto

	if err != nil || err0 != nil || err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		err = helper.Ternary(err != nil, err, err0)
		err = helper.Ternary(err != nil, err, err1)
		err = helper.Ternary(err != nil, err, err2)
		err = helper.Ternary(err != nil, err, err3)
		err = helper.Ternary(err != nil, err, err4)
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	if optType == "edit" && userDto.Id == "" {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, "请调用添加接口", nil)
		c.ServeJSON()
		return
	}
	if optType == "add" && userDto.Id != "" {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, "请调用编辑接口", nil)
		c.ServeJSON()
		return
	}

	data, err := c.UserService.SaveUser(baseParamDto, &unitUserSaveDto)
	if err != nil {
		global.Log.Error("错误 %v", err)
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// 删除组织单位用户
// @Summary 删除组织单位用户
// @Description 删除组织单位用户
// @Tags 用户管理
// @Accept  json
// @Produce  json
// @Param   unitUserId query string true "组织单位用户ID"
// @Success 200 {object} dto.Response "{"code": 200, "data": [...]}"
// @Router /admin_mchnt/system-user/del [post]
func (c *UserController) Del() {
	type reqStruct struct {
		Id []string `json:"id"`
	}

	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	req, err0 := helper.GetReqBody[*reqStruct](c.Ctx)
	if err != nil || err0 != nil {
		err = helper.Ternary(err != nil, err, err0)
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.UserService.DelUnitUser(baseParamDto, req.Id)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}

// 角色树形列表
// @Summary 角色树形列表
// @Description 角色树形列表
// @Tags 角色
// @Accept  json
// @Produce  json
// @Param   selectUnitIds query string true "selectUnitIds"
// @Success 200 {object} dto.Response "{"code": 200, "data": [...]}"
// @Router /admin_mchnt/system-user/get-role-tree [get]
func (c *UserController) GetRoleTree() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	selectUnitIds := helper.Ternary(c.GetString("selectUnitIds") != "", strings.Split(c.GetString("selectUnitIds"), ","), []string{baseParamDto.UnitId})
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	data, err := c.UserService.GetUnitRoleTree(baseParamDto, selectUnitIds)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}
