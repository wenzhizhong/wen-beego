package page_dto

import "WenBeego/apps/common/dto_vo"

type GenerateCodeListReqDto struct {
	dto_vo.BaseParamDto
	dto_vo.ReqDataListDto
	Keyword string `json:"keyword"`
}
