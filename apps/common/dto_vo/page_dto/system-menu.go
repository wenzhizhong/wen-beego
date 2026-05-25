package page_dto

import "WenBeego/apps/common/dto_vo"

// 系统管理-菜单列表
type SystemMenuListReqDto struct {
	dto_vo.BaseParamDto
	dto_vo.ReqDataListDto
	SelectUnitIds []string
	Title         string
}
