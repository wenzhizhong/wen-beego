package page_dto

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/models/base_model"
)

// 系统管理-内部用户列表-请求参数
type SystemUserListReqDto struct {
	dto.BaseParamDto
	dto.ReqDataListDto
	UserName      string
	Phone         string
	SelectUnitIds []string
	RoleIds       []string
	DeptIds       []string
}

// 系统管理-内部用户列表
type SystemUserListDto struct {
	Id string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	base_model.UnitUser
	base_model.UnitUserProfile
	DeptIds   string `json:"dept_ids" gorm:"->"`
	DeptNames string `json:"dept_names" gorm:"->"`
	RoleIds   string `json:"role_ids" gorm:"->"`
	RoleNames string `json:"role_names" gorm:"->"`
}
