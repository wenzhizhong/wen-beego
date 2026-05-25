package role_dto

import "WenBeego/apps/common/models/base_model"

type UnitRoleDto struct {
	base_model.UnitRole
	RoleClassifyName string `json:"role_classify_name" gorm:"->"`
}

type RoleMenuSaveDto struct {
	RoleId  string   `json:"id"`
	MenuIds []string `json:"menuIds"`
}
