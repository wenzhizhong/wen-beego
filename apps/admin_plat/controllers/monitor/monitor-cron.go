package monitor

import (
	paltMonitorService "WenBeego/apps/admin_plat/services/monitor"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto/cron_dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/helper"
)

type CronController struct {
	commonControllers.AdminBaseController
	cronService paltMonitorService.CronService
}

// 定时任务列表
func (c *CronController) Get() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	reqDataListDto, err1 := helper.GetReqDataListDto(&c.Controller)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}

	reqDto := &page_dto.MonitorCronListReqDto{}
	reqDto.BaseParamDto = baseParamDto
	reqDto.ReqDataListDto = reqDataListDto
	data, err := c.cronService.GetCronList(reqDto)

	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}

	c.Data["json"] = helper.Response(200, "success", *data)
	c.ServeJSON()
}

// 获取可用定时任务列表
func (c *CronController) GetAvaibleCronList() {
	data, err := c.cronService.GetAvaibleCronList()
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// 添加定时任务
func (c *CronController) Add() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	dtoData, err1 := helper.GetReqBody[cron_dto.UnitCronDto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.cronService.AddTask(baseParamDto, dtoData)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}

// 修改定时任务
func (c *CronController) Edit() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	dtoData, err1 := helper.GetReqBody[cron_dto.UnitCronDto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.cronService.UpdateTask(baseParamDto, dtoData)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}

// 改变定时任务状态
func (c *CronController) ChangeStatus() {
	dtoData, _ := helper.GetReqBody[cron_dto.UnitCronDto](c.Ctx)
	if dtoData.Status == 1 {
		c.startTask()
	} else {
		c.stopTask()
	}
}

func (c *CronController) startTask() {

	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	dtoData, err1 := helper.GetReqBody[cron_dto.UnitCronDto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.cronService.StartTask(baseParamDto, dtoData)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}
func (c *CronController) stopTask() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	dtoData, err1 := helper.GetReqBody[cron_dto.UnitCronDto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.cronService.StopTask(baseParamDto, dtoData)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}

// 删除定时任务
func (c *CronController) Del() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	dtoData, err1 := helper.GetReqBody[cron_dto.UnitCronDto](c.Ctx)
	if err != nil || err1 != nil {
		err = helper.Ternary(err != nil, err, err1)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	err = c.cronService.DelTask(baseParamDto, dtoData)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", nil)
	c.ServeJSON()
}
