package user_vo

import (
	"WenBeego/apps/common/models/base_model"
)

// 系统管理-内部用户列表
type SystemUserListVo struct {
	Id string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	base_model.UnitUser
	base_model.UnitUserProfile
	UserId    string `json:"user_id" gorm:"->"`
	DeptIds   string `json:"dept_ids" gorm:"->"`
	DeptNames string `json:"dept_names" gorm:"->"`
	RoleIds   string `json:"role_ids" gorm:"->"`
	RoleNames string `json:"role_names" gorm:"->"`
}
