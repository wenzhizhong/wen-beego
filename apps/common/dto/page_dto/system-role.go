package page_dto

import "WenBeego/apps/common/dto"

type SystemRoleListReqDto struct {
	dto.BaseParamDto
	dto.ReqDataListDto
	SelectUnitIds    []string
	RoleName         string
	Status           int
	RoleClassifyName string
}
