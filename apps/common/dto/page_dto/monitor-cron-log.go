package page_dto

import (
	"WenBeego/apps/common/dto"
	"time"
)

type MonitorCronLogListReqDto struct {
	dto.ReqDataListDto

	NameEn    string
	CreatedAt *time.Time
}
type MonitorCronLogListRespDto struct {
}
