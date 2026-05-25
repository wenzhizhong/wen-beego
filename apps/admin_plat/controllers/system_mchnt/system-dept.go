package system_mchnt

import (
	systemService "WenBeego/apps/admin_plat/services/system_mchnt"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto_vo/dept_dto"
	"WenBeego/apps/common/dto_vo/page_dto"
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
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Router /admin_plat/admin_mchnt/system-dept/get [get]

func (c *DeptController) Get() {
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
	deptDto := page_dto.SystemDeptListReqDto{}
	deptDto.BaseParamDto = baseParamDto
	deptDto.ReqDataListDto = reqDataListDto
	deptDto.Name = c.GetString("name")
	deptDto.SelectUnitIds = helper.Ternary(c.GetString("selectUnitIds") != "", strings.Split(c.GetString("selectUnitIds"), ","), []string{baseParamDto.UnitId})

	data, err := c.deptService.GetUnitDeptList(deptDto)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// 组织架构树形列表
// @Summary 组织架构树形列表
// @Description 组织架构树形列表
// @Tags 组织架构
// @Accept  json
// @Produce  json
// @Param   selectUnitIds query string true "selectUnitIds"
// @Success 200 {object}  interface{} "返回结果"
// @Router /admin_plat/admin_mchnt/system-dept/get-dept-tree [get]
func (c *DeptController) GetUnitDeptTree() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	selectUnitIds := helper.Ternary(c.GetString("selectUnitIds") != "", strings.Split(c.GetString("selectUnitIds"), ","), []string{baseParamDto.UnitId})
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	data, err := c.deptService.GetUnitDeptTree(baseParamDto, selectUnitIds)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// 搜索可用的组织架构负责人
// @Summary 搜索可用的组织架构负责人
// @Description 搜索可用的组织架构负责人
// @Tags 组织架构
// @Accept  json
// @Produce  json
// @Param   selectUnitIds query string true "selectUnitIds"
// @Success 200 {object} interface{} "返回结果"
// @Router /admin_plat/admin_mchnt/system-dept/get-dept-principal [get]
func (c *DeptController) GetUnitDeptPrincipal() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)

	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	reqDataListDto, err2 := helper.GetReqDataListDto(&c.Controller)
	if err2 != nil {
		c.Data["json"] = helper.Response(500, err2.Error(), nil)
		c.ServeJSON()
		return
	}

	deptPrincipalDto := page_dto.SystemDeptPrincipalReqDto{}
	deptPrincipalDto.BaseParamDto = baseParamDto
	deptPrincipalDto.ReqDataListDto = reqDataListDto
	deptPrincipalDto.SelectUnitIds = helper.Ternary(c.GetString("selectUnitIds") != "", strings.Split(c.GetString("selectUnitIds"), ","), []string{baseParamDto.UnitId})
	deptPrincipalDto.Keyword = c.GetString("keyword")

	data, err := c.deptService.GetUnitDeptPrincipal(baseParamDto, deptPrincipalDto)
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
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Route /admin_mchnt/system-dept/add

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
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Route /admin_mchnt/system-dept/edit
func (c *DeptController) Edit() {
	c.save("edit")
}
func (c *DeptController) save(optType string) {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
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

// 删除组织架构
// @Summary 删除组织架构
// @Description 删除组织架构
// @Tags 系统管理-组织架构
// @Accept application/json
// @Produce application/json
// @Param reqDto body dept_dto.UnitDeptDto true "删除组织架构参数"
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Route /admin_mchnt/system-dept/del
func (c *DeptController) Del() {
	deptDto, err1 := helper.GetReqBody[dept_dto.UnitDeptDto](c.Ctx)
	baseParamDto, err2 := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	if err1 != nil {
		c.Data["json"] = helper.Response(500, err1.Error(), nil)
		c.ServeJSON()
		return
	}
	if err2 != nil {
		c.Data["json"] = helper.Response(500, err2.Error(), nil)
		c.ServeJSON()
		return
	}
	err := c.deptService.Del(baseParamDto, deptDto)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}
