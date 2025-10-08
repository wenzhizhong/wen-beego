package page_dto

import "WenBeego/apps/common/dto"

// 系统管理-用户列表
type SystemUserListReqDto struct {
	dto.BaseParamDto
	dto.ReqDataListDto
	UserName string
	Phone    string
}
