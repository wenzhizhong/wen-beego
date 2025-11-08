package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.UserRoleItf = (*PlatUserRole)(nil)

type PlatUserRole struct {
	base_model.UnitUserRole
}

func (m *PlatUserRole) TableName() string {
	return `plat_user_role`
}

func (m *PlatUserRole) GetId() string {
	return m.Id
}
func (m *PlatUserRole) GetUserId() string {
	return m.UserId
}
func (m *PlatUserRole) GetRoleId() string {
	return m.RoleId
}
func (m *PlatUserRole) GetDeleted() int {
	return m.Deleted
}
