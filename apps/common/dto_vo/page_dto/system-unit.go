package page_dto

import "WenBeego/apps/common/dto_vo"

// 系统管理-内部组织单位列表
type SystemUnitListReqDto struct {
	dto_vo.BaseParamDto
	dto_vo.ReqDataListDto
	SelectUnitIds []string
	Name          string
	Code          string
	Status        int
}
