package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.RoleMenuItf = (*MchntRoleMenu)(nil)

type MchntRoleMenu struct {
	base_model.UnitRoleMenu
}

func (m *MchntRoleMenu) TableName() string {
	return `mchnt_role_menu`
}

func (m *MchntRoleMenu) GetId() string {
	return m.Id
}
func (m *MchntRoleMenu) GetRoleId() string {
	return m.RoleId
}
func (m *MchntRoleMenu) GetMenuId() string {
	return m.MenuId
}
