package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.RoleMenuItf = (*PlatRoleMenu)(nil)

type PlatRoleMenu struct {
	base_model.UnitRoleMenu
}

func (m *PlatRoleMenu) TableName() string {
	return `plat_role_menu`
}

func (m *PlatRoleMenu) GetId() string {
	return m.Id
}
func (m *PlatRoleMenu) GetRoleId() string {
	return m.RoleId
}
func (m *PlatRoleMenu) GetMenuId() string {
	return m.MenuId
}
