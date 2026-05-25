package page_dto

import (
	"WenBeego/apps/common/dto_vo"
)

// 系统管理-内部用户列表-请求参数
type SystemUserListReqDto struct {
	dto_vo.BaseParamDto
	dto_vo.ReqDataListDto
	UserName      string
	Phone         string
	SelectUnitIds []string
	RoleIds       []string
	DeptIds       []string
}
