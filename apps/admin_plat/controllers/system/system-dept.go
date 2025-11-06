package system

import (
	systemService "WenBeego/apps/admin_plat/services/system"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto/dept_dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/helper"
	"errors"
	"strings"
)

type DeptController struct {
	commonControllers.AdminBaseController
	deptService systemService.DeptService
}

// 获取组织架构列表
// @Summary 获取组织架构列表
// @Description 获取组织架构列表
// @Tags 系统管理-组织架构
// @Accept application/json
// @Produce application/json
// @Param selectUnitIds query string true "父级ID"
// @Success 200 {object} dto.RespDataListDto
// @Router /admin_plat/system-dept/get [get]

func (c *DeptController) Get() {
	userId := c.Ctx.Input.GetData("userId")
	unitId := c.Ctx.Input.GetData("unitId")

	baseParamDto, err := helper.GetBaseParamDto(c.Ctx.Request.Host, c.ModuleName, unitId.(string), userId.(string))
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
	deptDto := page_dto.SystemDeptListReqDto{}
	deptDto.BaseParamDto = baseParamDto
	deptDto.ReqDataListDto = reqDataListDto
	deptDto.Name = c.GetString("name")
	deptDto.SelectUnitIds = helper.Ternary(c.GetString("selectUnitIds") != "", strings.Split(c.GetString("selectUnitIds"), ","), []string{})

	data, err := c.deptService.GetUnitDeptList(deptDto)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// 新增组织架构
// @Summary 新增组织架构
// @Description 新增组织架构
// @Tags 系统管理-组织架构
// @Accept application/json
// @Produce application/json
// @Param reqDto body dept_dto.UnitDeptDto true "新增组织架构参数"
// @Success 200 {object} dto.RespDataDto

func (c *DeptController) Add() {
	c.save("add")
}

// 修改组织架构
// @Summary 修改组织架构
// @Description 修改组织架构
// @Tags 系统管理-组织架构
// @Accept application/json
// @Produce application/json
// @Param reqDto body dept_dto.UnitDeptDto true "修改组织架构参数"
// @Success 200 {object} dto.RespDataDto
func (c *DeptController) Edit() {
	c.save("edit")
}
func (c *DeptController) save(optType string) {
	userId := c.Ctx.Input.GetData("userId")
	unitId := c.Ctx.Input.GetData("unitId")
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx.Request.Host, c.ModuleName, unitId.(string), userId.(string))
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	deptDto, err := helper.GetReqBody[dept_dto.UnitDeptDto](c.Ctx)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	var err2 error
	if optType == "add" && deptDto.Id != "" {
		err2 = errors.New("请调用编辑接口")
	} else if optType == "edit" && deptDto.Id == "" {
		err2 = errors.New("请调用添加接口")
	}
	if err2 != nil {
		c.Data["json"] = helper.Response(500, err2.Error(), nil)
		c.ServeJSON()
		return
	}

	data, err := c.deptService.SaveUnitDept(baseParamDto, deptDto)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}
