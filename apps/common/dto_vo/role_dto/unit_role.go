package role_dto

import "WenBeego/apps/common/models/base_model"

type UnitRoleDto struct {
	base_model.UnitRole
}

type RoleMenuSaveDto struct {
	RoleId  string   `json:"id"`
	MenuIds []string `json:"menuIds"`
}
