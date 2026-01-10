package page_dto

import "WenBeego/apps/common/dto"

// 系统管理-菜单列表
type SystemMenuListReqDto struct {
	dto.BaseParamDto
	dto.ReqDataListDto
	SelectUnitIds []string
	Title         string
}
