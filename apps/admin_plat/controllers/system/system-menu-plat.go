package system

// 系统管理-平台菜单管理
import (
	systemService "WenBeego/apps/admin_plat/services/system"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto/menu_dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/helper"
	"errors"
)

type MenuPlatController struct {
	commonControllers.AdminBaseController
	menuService systemService.MenuService
}

// 系统管理-平台菜单管理
// @Summary 平台菜单管理
// @Description 平台菜单管理
// @Tags 系统管理-平台菜单管理
// @Accept application/json
// @Produce application/json
// @Param parentUnitId query string true "父级ID"
// @Success 200 {object} dto.RespDataListDto
// @Router /admin_plat/system-menu/get-plat [get]
func (c *MenuPlatController) Get() {
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
	pageDto := page_dto.SystemMenuListReqDto{}
	pageDto.BaseParamDto = baseParamDto
	pageDto.ReqDataListDto = reqDataListDto
	pageDto.Title = c.GetString("title")

	data, err := c.menuService.GetMenuList(pageDto, "admin_plat")
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}
func (c *MenuPlatController) save(optType string) {
	menuDto, err1 := helper.GetReqBody[menu_dto.MenuDto](c.Ctx)
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
	if optType == "add" && menuDto.Id != "" {
		err3 = errors.New("请调用编辑接口")
	} else if optType == "edit" && menuDto.Id == "" {
		err3 = errors.New("请调用添加接口")
	}
	if err3 != nil {
		c.Data["json"] = helper.Response(500, err3.Error(), nil)
		c.ServeJSON()
		return
	}

	data, err := c.menuService.Save(baseParamDto, menuDto, "admin_plat")
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// 系统管理-新增平台菜单管理
// @Summary 新增平台菜单管理
// @Description 新增平台菜单管理
// @Tags 系统管理-平台菜单管理
// @Accept application/json
// @Produce application/json
// @Param menuDto body menu_dto.MenuDto true "平台菜单管理"
// @Success 200 {object} dto.Response
// @Router /admin_plat/system-menu/add-plat [post]

func (c *MenuPlatController) Add() {
	c.save("add")
}

// 系统管理-编辑平台菜单管理
// @Summary 编辑平台菜单管理
// @Description 编辑平台菜单管理
// @Tags 系统管理-平台菜单管理
// @Accept application/json
// @Produce application/json
// @Param menuDto body menu_dto.MenuDto true "平台菜单管理"
// @Success 200 {object} dto.Response
// @Router /admin_plat/system-menu/edit-plat [post]

func (c *MenuPlatController) Edit() {
	c.save("edit")
}

// @Summary 删除平台菜单管理
// @Description 删除平台菜单管理
// @Tags 系统管理-平台菜单管理
// @Accept application/json
// @Produce application/json
// @Param ids body dto.menuDto true "平台菜单管理"
// @Success 200 {object} dto.Response
// @Router /admin_plat/system-menu/del-plat [post]
func (c *MenuPlatController) Del() {
	menuDto, err1 := helper.GetReqBody[menu_dto.MenuDto](c.Ctx)
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
	err := c.menuService.Del(baseParamDto, menuDto, "admin_plat")
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}
