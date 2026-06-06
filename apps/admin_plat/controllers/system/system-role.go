package system

import (
	systemService "WenBeego/apps/admin_plat/services/system"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/dto_vo/role_dto"
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/common/helper"
	"errors"
	"strings"
)

type RoleController struct {
	commonControllers.AdminBaseController
	roleService systemService.RoleService
}

// 获取角色列表
// @Summary 获取角色列表
// @Description 获取角色列表
// @Tags 系统管理-角色
// @Accept application/json
// @Produce application/json
// @Param selectUnitIds query string true "父级ID"
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Router /admin_plat/system-role/get [get]

func (c *RoleController) Get() {
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
	roleDto := page_dto.SystemRoleListReqDto{}
	roleDto.BaseParamDto = baseParamDto
	roleDto.ReqDataListDto = reqDataListDto
	roleDto.RoleName = c.GetString("role_name")
	roleDto.Status, _ = c.GetInt("status")
	roleDto.RoleClassifyName = c.GetString("role_classify_name")
	roleDto.SelectUnitIds = helper.Ternary(c.GetString("selectUnitIds") != "", strings.Split(c.GetString("selectUnitIds"), ","), []string{baseParamDto.UnitId})

	data, err := c.roleService.GetUnitRoleList(roleDto)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// 新增角色
// @Summary 新增角色
// @Description 新增角色
// @Tags 系统管理-角色
// @Accept application/json
// @Produce application/json
// @Param reqDto body role_dto.UnitRoleDto true "新增角色参数"
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Route /admin_plat/system-role/add

func (c *RoleController) Add() {
	c.save("add")
}

// 修改角色
// @Summary 修改角色
// @Description 修改角色
// @Tags 系统管理-角色
// @Accept application/json
// @Produce application/json
// @Param reqDto body role_dto.UnitRoleDto true "修改角色参数"
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Route /admin_plat/system-role/edit
func (c *RoleController) Edit() {
	c.save("edit")
}
func (c *RoleController) save(optType string) {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	roleDto, err := helper.GetReqBody[role_dto.UnitRoleDto](c.Ctx)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	var err2 error
	if optType == "add" && roleDto.Id != "" {
		err2 = errors.New("请调用编辑接口")
	} else if optType == "edit" && roleDto.Id == "" {
		err2 = errors.New("请调用添加接口")
	}
	if err2 != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err2.Error(), nil)
		c.ServeJSON()
		return
	}

	data, err := c.roleService.SaveUnitRole(baseParamDto, roleDto)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// 删除角色
// @Summary 删除角色
// @Description 删除角色
// @Tags 系统管理-角色
// @Accept application/json
// @Produce application/json
// @Param reqDto body role_dto.UnitRoleDto true "删除角色参数"
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Route /admin_plat/system-role/del
func (c *RoleController) Del() {
	roleDto, err1 := helper.GetReqBody[role_dto.UnitRoleDto](c.Ctx)
	baseParamDto, err2 := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	if err1 != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err1)
		c.Data["json"] = helper.Response(500, err1.Error(), nil)
		c.ServeJSON()
		return
	}
	if err2 != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err2)
		c.Data["json"] = helper.Response(500, err2.Error(), nil)
		c.ServeJSON()
		return
	}
	err := c.roleService.Del(baseParamDto, roleDto)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}

// 获取角色可选菜单
// @Summary 获取角色可选菜单
// @Description 获取角色可选菜单
// @Tags 系统管理-角色
// @Accept application/json
// @Produce application/json
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Router /admin_plat/system-role/role-menu [get]
func (c *RoleController) RoleMenu() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	selectUnitIds := helper.Ternary(c.GetString("selectUnitIds") != "", strings.Split(c.GetString("selectUnitIds"), ","), []string{baseParamDto.UnitId})
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	data, err := c.roleService.GetRoleMenu(baseParamDto, selectUnitIds)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// 获取角色已选菜单
// @Summary 获取角色已选菜单
// @Description 获取角色已选菜单
// @Tags 系统管理-角色
// @Accept application/json
// @Produce application/json
// @Param id query string true "角色id"
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Router /admin_plat/system-role/role-menu-ids [get]
func (c *RoleController) RoleMenuIds() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	roleId := c.GetString("id")
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	data, err := c.roleService.GetRoleMenuIds(baseParamDto, roleId)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}

	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// 保存角色菜单
// @Summary 保存角色菜单
// @Description 保存角色菜单
// @Tags 系统管理-角色
// @Accept application/json
// @Produce application/json
// @Param reqDto body role_dto.RoleMenuSaveDto true "保存角色菜单参数"
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Router /admin_plat/system-role/role-menu-save [post]

func (c *RoleController) RoleMenuSave() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	roleMenuSaveDto, err1 := helper.GetReqBody[role_dto.RoleMenuSaveDto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}

	err = c.roleService.RoleMenuSave(baseParamDto, roleMenuSaveDto)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}
