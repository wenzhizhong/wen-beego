package {{.MenuModule}}

import (
	commonControllers "WenBeego/apps/common/controller"
	_ "WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/{{.MenuModule}}_dto"
	"WenBeego/apps/common/helper"
	{{.MenuModule}}Service "WenBeego/apps/{{.AppModule}}/services/{{.MenuModule}}"
)

type {{.ModelName}}Controller struct {
	commonControllers.AdminBaseController
	Service {{.MenuModule}}Service.{{.ModelName}}Service
}

func (c *{{.ModelName}}Controller) Get() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	reqDataListDto, err1 := helper.GetReqDataListDto(&c.Controller)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	keyword := c.GetString("keyword")
	data, err := c.Service.GetList(baseParamDto, reqDataListDto.PageSize, reqDataListDto.Offset, keyword)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", data)
	}
	c.ServeJSON()
}

func (c *{{.ModelName}}Controller) Add() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	dtoData, err1 := helper.GetReqBody[{{.MenuModule}}_dto.{{.ModelName}}Dto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.Service.Add(baseParamDto, dtoData)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", nil)
	}
	c.ServeJSON()
}

func (c *{{.ModelName}}Controller) Edit() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	dtoData, err1 := helper.GetReqBody[{{.MenuModule}}_dto.{{.ModelName}}Dto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.Service.Edit(baseParamDto, dtoData)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", nil)
	}
	c.ServeJSON()
}

func (c *{{.ModelName}}Controller) Del() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	dtoData, err1 := helper.GetReqBody[{{.MenuModule}}_dto.{{.ModelName}}Dto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.Service.Del(baseParamDto, dtoData.Id)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", nil)
	}
	c.ServeJSON()
}

func (c *{{.ModelName}}Controller) Detail() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	
	id := c.Controller.GetString("id")
	dtoData := system_dto.GenerateFormDto{}
	dtoData.Id = id
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	data, err := c.Service.GetDetail(baseParamDto, dtoData.Id)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", data)
	}
	c.ServeJSON()
}
