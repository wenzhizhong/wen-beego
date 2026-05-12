package page_dto

import "WenBeego/apps/common/dto"

type GenerateCodeListReqDto struct {
	dto.BaseParamDto
	dto.ReqDataListDto
	Keyword string `json:"keyword"`
}
