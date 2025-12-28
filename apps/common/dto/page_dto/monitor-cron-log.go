package page_dto

import (
	"WenBeego/apps/common/dto"
	"time"
)

type MonitorCronLogListReqDto struct {
	dto.ReqDataListDto

	CrontabId string
	CreatedAt *time.Time
}
type MonitorCronLogListRespDto struct {
}
