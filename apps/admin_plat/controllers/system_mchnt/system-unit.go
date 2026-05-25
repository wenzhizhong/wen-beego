package system_mchnt

// 系统管理-内部组织管理
import (
	systemService "WenBeego/apps/admin_plat/services/system_mchnt"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/dto_vo/unit_dto"
	"WenBeego/apps/common/helper"
	"errors"
)

type UnitController struct {
	commonControllers.AdminBaseController
	UnitService systemService.UnitService
}

// 系统管理-获取内部组织管理
// @Summary 获取内部组织管理
// @Description 获取内部组织管理
// @Tags 系统管理-内部组织管理
// @Accept application/json
// @Produce application/json
// @Param parentUnitId query string true "父级ID"
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Router /admin_plat/admin_mchnt/system-unit/get [get]
func (c *UnitController) Get() {
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
	unitDto := page_dto.SystemUnitListReqDto{}
	unitDto.BaseParamDto = baseParamDto
	unitDto.ReqDataListDto = reqDataListDto
	unitDto.Name = c.GetString("name")
	unitDto.Code = c.GetString("code")
	unitDto.Status, _ = c.GetInt("status")

	data, err := c.UnitService.GetUnitList(unitDto)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}
func (c *UnitController) save(optType string) {
	unitDto, err1 := helper.GetReqBody[unit_dto.UnitDto](c.Ctx)
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

	var err3 error
	if optType == "add" && unitDto.Id != "" {
		err3 = errors.New("请调用编辑接口")
	} else if optType == "edit" && unitDto.Id == "" {
		err3 = errors.New("请调用添加接口")
	}
	if err3 != nil {
		c.Data["json"] = helper.Response(500, err3.Error(), nil)
		c.ServeJSON()
		return
	}

	data, err := c.UnitService.Save(baseParamDto, unitDto)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// 系统管理-新增内部组织管理
// @Summary 新增内部组织管理
// @Description 新增内部组织管理
// @Tags 系统管理-内部组织管理
// @Accept application/json
// @Produce application/json
// @Param unitDto body unit_dto.UnitDto true "内部组织管理"
// @Success 200 {object} dto.Response "返回结果"
// @Router /admin_plat/admin_mchnt/system-unit/add [post]

func (c *UnitController) Add() {
	c.save("add")
}

// 系统管理-编辑内部组织管理
// @Summary 编辑内部组织管理
// @Description 编辑内部组织管理
// @Tags 系统管理-内部组织管理
// @Accept application/json
// @Produce application/json
// @Param unitDto body unit_dto.UnitDto true "内部组织管理"
// @Success 200 {object} dto.Response "返回结果"
// @Router /admin_plat/admin_mchnt/system-unit/edit [post]

func (c *UnitController) Edit() {
	c.save("edit")
}

// @Summary 删除内部组织管理
// @Description 删除内部组织管理
// @Tags 系统管理-内部组织管理
// @Accept application/json
// @Produce application/json
// @Param ids body unit_dto.UnitDto true "内部组织管理"
// @Success 200 {object} dto.Response "返回结果"
// @Router /admin_plat/admin_mchnt/system-unit/del [post]
func (c *UnitController) Del() {
	unitDto, err1 := helper.GetReqBody[unit_dto.UnitDto](c.Ctx)
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
	err := c.UnitService.Del(baseParamDto, unitDto)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}
