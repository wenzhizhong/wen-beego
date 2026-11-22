package page_dto

import "WenBeego/apps/common/dto_vo"

type SqliteDataSourceListReqDto struct {
	dto_vo.BaseParamDto
	dto_vo.ReqDataListDto
	Keyword string `json:"keyword"`
}
