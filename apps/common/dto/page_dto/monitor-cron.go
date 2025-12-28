package page_dto

import (
	"WenBeego/apps/common/dto"
)

// 系统管理-定时任务-请求参数
type MonitorCronListReqDto struct {
	dto.BaseParamDto
	dto.ReqDataListDto
	SelectUnitIds []string
}
