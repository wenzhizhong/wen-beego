package user_dto

import (
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
)

// 系统内部用户管理：用户保存
type UnitUserSaveDto struct {
	UserDto            UserDto
	UserProfileDto     UserProfileDto
	UnitUserDto        UnitUserDto
	UnitUserProfileDto UnitUserProfileDto
	UnitUserRoleDto    []UnitUserRoleDto
	UnitUserDeptDto    UnitUserDeptDto

	DeptId string   `json:"dept_id"`
	RoleId []string `json:"role_id"`
}

type UserDto struct {
	models.User
	Id string `json:"user_id"`
}
type UserProfileDto struct {
	models.UserProfile
	Id string `json:"user_id"`
}
type UnitUserDto struct {
	base_model.UnitUser
}

type UnitUserProfileDto struct {
	base_model.UnitUserProfile
}

// 用户全部字段
type UserAllDataDto struct {
	models.User
	models.UserProfile
	Deleted int    `json:"deleted" gorm:"not null;default:0;comment:是否删除"`
	Id      string `json:"id" gorm:"type:bpchar(36);primaryKey;comment:ID"`
}

type UnitUserRoleDto struct {
	base_model.UnitUserRole
}

type UnitUserDeptDto struct {
	base_model.UnitUserDept
}
