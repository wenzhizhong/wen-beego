package page_dto

import (
	"WenBeego/apps/common/dto_vo"
	"time"
)

type MonitorCronLogListReqDto struct {
	dto_vo.ReqDataListDto

	NameEn    string
	CreatedAt *time.Time
}
type MonitorCronLogListRespDto struct {
}
