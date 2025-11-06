package page_dto

import "WenBeego/apps/common/dto"

// 系统管理-内部部门列表
type SystemDeptListReqDto struct {
	dto.BaseParamDto
	dto.ReqDataListDto
	SelectUnitIds []string
	Name          string
}
