package monitor

import (
	platMonitorService "WenBeego/apps/admin_plat/services/monitor"
	commonControllers "WenBeego/apps/common/controller"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/helper"
	"time"
)

type CronLogController struct {
	commonControllers.AdminBaseController
	cronLogService platMonitorService.CronLogService
}

// 获取定时任务日志列表
// @Summary 获取定时任务日志列表
// @Description 获取定时任务日志列表
// @Tags 定时任务日志
// @Accept json
// @Produce json
// @Param req body page_dto.MonitorCronLogListReqDto true "请求参数"
// @Success 200 {object} dto.RespDataListDto "返回结果"
// @Router /plat_admin/monitor-cron-log/get [get]
func (c *CronLogController) Get() {
	// baseParamDto, err := helper.GetBaseParamDto(c.Ctx, c.ModuleName)
	reqDataListDto, err := helper.GetReqDataListDto(&c.Controller)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	createdAtStr := c.GetString("created_at")
	createdAt := time.Time{}
	if createdAtStr != "" {
		createdAt, err = time.Parse("2006-01-02 15:04:05", createdAtStr)
		if err != nil {
			c.Data["json"] = helper.Response(500, err.Error(), nil)
			c.ServeJSON()
			return
		}
	}
	req := &page_dto.MonitorCronLogListReqDto{
		ReqDataListDto: reqDataListDto,
		NameEn:         c.GetString("name_en"),
		CreatedAt:      helper.Ternary(createdAtStr != "", &createdAt, nil),
	}

	data, err := c.cronLogService.GetCronLogList(req)
	if err != nil {
		c.Data["json"] = helper.Response(500, err.Error(), nil)
		c.ServeJSON()
		return
	}
	c.Data["json"] = helper.Response(200, "success", data)
	c.ServeJSON()
}
