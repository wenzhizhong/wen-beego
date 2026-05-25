package page_dto

import "WenBeego/apps/common/dto_vo"

type SystemRoleListReqDto struct {
	dto_vo.BaseParamDto
	dto_vo.ReqDataListDto
	SelectUnitIds    []string
	RoleName         string
	Status           int
	RoleClassifyName string
}
