package auth_vo

import (
	"WenBeego/apps/common/models/base_model"
)

type RoleMenuVo struct {
	//table: menu
	base_model.UnitMenu

	//table: role_menu
	RoleId string `json:"role_id" gorm:"type:bpchar(36);not null;comment:角色ID"`
	MenuId string `json:"menu_id" gorm:"type:bpchar(36);not null;comment:菜单权限ID"`

	//other：api角色权限
	Roles []string `json:"roles" gorm:"-"`
	Auths []string `json:"auths" gorm:"-"`
}
