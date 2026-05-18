package upload

import (
	uploadService "WenBeego/apps/admin_plat/services/upload"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
)

type UploadController struct {
	commonControllers.AdminBaseController
	uploadService uploadService.UploadService
}

// 上传
// @Summary 上传
// @Description 上传
// @Tags 上传
// @Accept multipart/form-data
// @Produce  json
// @Param file formData file true "file"
// @Success 200 {object} dto.Response "返回结果"
// @Router /admin_plat/upload/upload [post]
func (c *UploadController) Upload() {
	userId := c.Ctx.Input.GetData("userId")
	unitId := c.Ctx.Input.GetData("unitId")
	file, fileHeader, err := c.GetFile("file")
	postData := c.Ctx.Request.PostForm
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}

	data, err := c.uploadService.Upload(userId.(string), unitId.(string), &file, fileHeader, postData, c.ModuleName)
	file.Close()
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}

	c.Data["json"] = helper.Response(200, "上传成功", data)
	c.ServeJSON()
}

// 分片上传
// @Summary 分片上传
// @Description 分片上传
// @Tags 上传
// @Accept multipart/form-data
// @Produce  json
// @Param file formData file true "file"
// @Param chunkNumber formData int true "chunkNumber"
// @Param chunkSize formData int true "chunkSize"
// @Param currentChunkSize formData int true "currentChunkSize"
// @Param totalSize formData int true "totalSize"
// @Param totalChunks formData int true "totalChunks"
// @Param filename formData string true "filename"
// @Success 200 {object} dto.Response "返回结果"
// @Router /admin_plat/upload/vue-slice-upload [post]
// @Security ApiKeyAuth

func (c *UploadController) VueSliceUpload() {
	userId := c.Ctx.Input.GetData("userId")
	unitId := c.Ctx.Input.GetData("unitId")
	file, fileHeader, err := c.GetFile("file")
	postData := c.Ctx.Request.PostForm

	if err != nil {
		c.Ctx.Output.SetStatus(210)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	data, err := c.uploadService.VueSliceUpload(userId.(string), unitId.(string), &file, fileHeader, postData, c.ModuleName)
	file.Close()

	if data.HttpCode == 0 {
		data.HttpCode = 200
	}
	c.Ctx.Output.SetStatus(data.HttpCode)
	if err != nil {
		c.Ctx.Output.SetStatus(data.HttpCode)
		global.Log.Error("VueSliceUpload() Error:\n" + err.Error())
		c.Data["json"] = helper.Response(500, err.Error(), data)
		c.ServeJSON()
		return
	}

	c.Data["json"] = helper.Response(200, "上传成功", data)
	c.ServeJSON()
}

// 分片上传检查
// @Summary 分片上传检查
// @Description 分片上传检查
// @Tags 上传
// @Accept multipart/form-data
// @Produce  json
// @Param identifier formData string true "identifier"
// @Param chunkNumber formData int true "chunkNumber"
// @Param totalChunks formData int true "totalChunks"
// @Success 200 {object} dto.Response "返回结果"
// @Router /admin_plat/upload/vue-slice-upload-check [post]
func (c *UploadController) VueSliceUploadCheck() {
	userId := c.Ctx.Input.GetData("userId")
	unitId := c.Ctx.Input.GetData("unitId")
	fileMd5 := c.GetString("identifier")
	sliceIndex := c.GetString("chunkNumber")
	sliceTotal := c.GetString("totalChunks")
	data, err := c.uploadService.VueSliceUploadCheck(userId.(string), unitId.(string), fileMd5, sliceIndex, sliceTotal)
	c.Ctx.Output.SetStatus(data.HttpCode)
	if err != nil {
		c.Ctx.Output.SetStatus(210)
		c.Data["json"] = helper.Response(500, err.Error(), data)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// 私有附件-链接签名
// @Summary 私有附件-链接签名
// @Description 私有附件-链接签名
// @Tags 上传
// @Accept application/x-wwww-form-urlencoded
// @Produce  json
// @Param urls path string true "urls"
// @Success 200 {object} dto.Response "返回结果"
// @Router /admin_plat/upload/link-sign [get]
// @Security ApiKeyAuth
func (c *UploadController) LinkSign() {
	host := c.Ctx.Request.Host
	urls := c.GetString("urls")
	data, err := c.uploadService.LinkSign(host, urls)

	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// 获取链接by ids
// @Summary 通过id获取访问链接
// @Description 通过id获取访问链接
// @Tags 上传
// @Accept application/x-wwww-form-urlencoded
// @Produce  json
// @Param ids path string true "ids"
// @Success 200 {object} dto.Response "返回结果"
// @Router /admin_plat/upload/get-link-by-id [get]
// @Security ApiKeyAuth
func (c *UploadController) GetLinkById() {
	host := c.Ctx.Request.Host
	ids := c.GetString("ids")
	data, err := c.uploadService.GetLinkById(host, ids)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}
