package systemtools

import (
	genService "WenBeego/apps/admin_plat/services/systemtools"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto_vo/generate_code_dto"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/common/helper"
	"path/filepath"
)

type GenerateCodeController struct {
	commonControllers.AdminBaseController
	GenService genService.GenerateCodeService
}

func (c *GenerateCodeController) GetDbTables() {
	data, err := c.GenService.GetDbTableList()
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", data)
	}
	c.ServeJSON()
}

func (c *GenerateCodeController) GetDbTableDetail() {
	dtoData, err := helper.GetReqBody[generate_code_dto.GetTableDetailDto](c.Ctx)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	data, err := c.GenService.GetDbTableDetail(dtoData.TableName)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", data)
	}
	c.ServeJSON()
}

func (c *GenerateCodeController) Save() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	dtoData, err1 := helper.GetReqBody[generate_code_dto.SaveFormDetailDto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.GenService.SaveFormDetail(baseParamDto, dtoData)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", nil)
	}
	c.ServeJSON()
}

func (c *GenerateCodeController) Del() {
	dtoData, err := helper.GetReqBody[generate_code_dto.DelFormDetailDto](c.Ctx)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.GenService.DelFormDetail(dtoData.Ids)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", nil)
	}
	c.ServeJSON()
}

func (c *GenerateCodeController) List() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	reqDataListDto, err1 := helper.GetReqDataListDto(&c.Controller)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	reqDto := &page_dto.GenerateCodeListReqDto{}
	reqDto.BaseParamDto = baseParamDto
	reqDto.ReqDataListDto = reqDataListDto
	reqDto.Keyword = c.GetString("keyword")

	data, err := c.GenService.GetGenerateCodeList(*reqDto)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", data)
	}
	c.ServeJSON()
}

func (c *GenerateCodeController) GetGenParams() {
	data, err := c.GenService.GetGenerateCodeParam()
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", data)
	}
	c.ServeJSON()
}

func (c *GenerateCodeController) Run() {
	dtoData, err := helper.GetReqBody[generate_code_dto.GenCodeRunDto](c.Ctx)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	data, err := c.GenService.GenerateCode(dtoData)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
	} else {
		c.Data["json"] = helper.Response(200, "success", data)
	}
	c.ServeJSON()
}

func (c *GenerateCodeController) Download() {
	dtoData, err := helper.GetReqBody[generate_code_dto.DownloadCodeDto](c.Ctx)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	zipPath, err := c.GenService.DownloadCode(dtoData.ZipPath)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	fullPath := filepath.Join(global.UploadsDir, "public", "code", zipPath)
	c.Ctx.Output.Download(fullPath)
}
