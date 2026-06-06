package monitor

import (
	monitorService "WenBeego/apps/admin_plat/services/monitor"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/dto_vo/queue_dlq_dto"
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/common/helper"
	"strconv"
)

type QueueDlqController struct {
	commonControllers.AdminBaseController
	QueueDlqService monitorService.QueueDlqService
}

// @Summary 获取死信队列失败日志列表
// @Description 获取死信队列失败日志列表
// @Tags 系统管理-死信队列日志
// @Accept application/json
// @Produce application/json
// @Param task_name query string false "任务名称"
// @Param status query int false "状态"
// @Param create_time_begin query string false "创建时间开始"
// @Param create_time_end query string false "创建时间结束"
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Router /admin_plat/monitor-queue-dlq/get [get]
func (c *QueueDlqController) Get() {
	baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	reqDataListDto, err2 := helper.GetReqDataListDto(&c.Controller)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	if err2 != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err2.Error(), nil)
		c.ServeJSON()
		return
	}

	dlqDto := page_dto.QueueDlqListReqDto{}
	dlqDto.BaseParamDto = baseParamDto
	dlqDto.ReqDataListDto = reqDataListDto
	dlqDto.TaskName = c.GetString("task_name")
	dlqDto.CreateTimeBegin = c.GetString("create_time_begin")
	dlqDto.CreateTimeEnd = c.GetString("create_time_end")
	dlqDto.GetTotal, _ = c.GetInt("getTotal", 0)

	status := -1
	statusStr := c.GetString("status")
	if statusStr != "" {
		status, _ = strconv.Atoi(statusStr)
	}
	dlqDto.Status = status

	data, err := c.QueueDlqService.GetList(dlqDto)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}

// @Summary 批量重新入队死信队列消息
// @Description 根据条件批量将死信队列失败日志中的消息重新添加到队列，支持按task_uuid单条或按条件批量
// @Tags 系统管理-死信队列日志
// @Accept application/json
// @Produce application/json
// @Param body body queue_dlq_dto.RequeueDto true "重新入队条件"
// @Success 200 {object} dto.Response "返回结果"
// @Router /admin_plat/monitor-queue-dlq/requeue [post]
func (c *QueueDlqController) Requeue() {
	reqDto, err1 := helper.GetReqBody[queue_dlq_dto.RequeueDto](c.Ctx)
	baseParamDto, err2 := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	if err1 != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err1)
		c.Data["json"] = helper.Response(500, err1.Error(), nil)
		c.ServeJSON()
		return
	}
	if err2 != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err2)
		c.Data["json"] = helper.Response(500, err2.Error(), nil)
		c.ServeJSON()
		return
	}

	count, err := c.QueueDlqService.Requeue(baseParamDto, reqDto)
	if err != nil {
		c.Ctx.Input.SetData(constant.CTX_ERROR_KEY, err)
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", map[string]int{"count": count})
	c.ServeJSON()
}
