package page_dto

import "WenBeego/apps/common/dto"

// 系统管理-内部组织单位列表
type SystemUnitListReqDto struct {
	dto.BaseParamDto
	dto.ReqDataListDto
	SelectUnitIds []string
	Name          string
	Code          string
	Status        int
}
