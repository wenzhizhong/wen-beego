package {{.MenuModule}}

import (
	commonControllers "WenBeego/apps/common/controller"
	_ "WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/{{.MenuModule}}"
	"WenBeego/apps/common/helper"
	{{.MenuModule}}Service "WenBeego/apps/{{.AppModule}}/services/{{.MenuModule}}"
	"encoding/json"
	"net/url"
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

	dtoStr := c.GetString("dto")
	dtoStr, _ = url.QueryUnescape(dtoStr)
	var searchDto {{.MenuModule}}_dto.{{.ModelName}}Dto
	if dtoStr != "" {
		json.Unmarshal([]byte(dtoStr), &searchDto)
	}
	{{- if .HasUnitId}}
	searchDto.SelectUnitIds = helper.Ternary(searchDto.SelectUnitIds != "", searchDto.SelectUnitIds, baseParamDto.UnitId)
	{{- end}}

	data, err := c.Service.GetList(baseParamDto, reqDataListDto.PageSize, reqDataListDto.Offset, searchDto)
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
	dtoData := {{.MenuModule}}_dto.{{.ModelName}}Dto{}
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
