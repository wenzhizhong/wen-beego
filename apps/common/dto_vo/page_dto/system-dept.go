package page_dto

import "WenBeego/apps/common/dto_vo"

// 系统管理-内部部门列表
type SystemDeptListReqDto struct {
	dto_vo.BaseParamDto
	dto_vo.ReqDataListDto
	SelectUnitIds []string
	Name          string
}

// 系统管理-获取可用的内部部门负责人列表
type SystemDeptPrincipalReqDto struct {
	dto_vo.BaseParamDto
	dto_vo.ReqDataListDto
	SelectUnitIds []string
	Keyword       string
}
