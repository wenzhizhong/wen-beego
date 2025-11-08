package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.UserRoleItf = (*MchntUserRole)(nil)

type MchntUserRole struct {
	base_model.UnitUserRole
}

func (m *MchntUserRole) TableName() string {
	return `mchnt_user_role`
}

func (m *MchntUserRole) GetId() string {
	return m.Id
}
func (m *MchntUserRole) GetUserId() string {
	return m.UserId
}
func (m *MchntUserRole) GetRoleId() string {
	return m.RoleId
}
func (m *MchntUserRole) GetDeleted() int {
	return m.Deleted
}
