package page_dto

import (
	"WenBeego/apps/common/dto_vo"
)

// 系统管理-定时任务-请求参数
type MonitorCronListReqDto struct {
	dto_vo.BaseParamDto
	dto_vo.ReqDataListDto
	SelectUnitIds []string
}
