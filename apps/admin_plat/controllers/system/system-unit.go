package system

// 系统管理-内部组织管理
import (
	systemService "WenBeego/apps/admin_plat/services/system"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/helper"
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
// @Success 200 {object} dto.RespDataListDto
// @Router /admin_plat/system-unit/get [get]
func (c *UnitController) GetUnitList() {
	userId := c.Ctx.Input.GetData("userId")
	unitId := c.Ctx.Input.GetData("unitId")

	baseParamDto, err := helper.GetBaseParamDto(c.ModuleName, unitId.(string), userId.(string))
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
