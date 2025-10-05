package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.RoleClassifyItf = (*MchntRoleClassify)(nil)

type MchntRoleClassify struct {
	base_model.UnitRoleClassify
}

func (m *MchntRoleClassify) TableName() string {
	return `mchnt_role_classify`
}

func (m *MchntRoleClassify) GetId() string {
	return m.Id
}
func (m *MchntRoleClassify) GetRoleId() string {
	return m.RoleId
}
func (m *MchntRoleClassify) GetUnitId() string {
	return m.UnitId
}
func (m *MchntRoleClassify) GetName() string {
	return m.Name
}
func (m *MchntRoleClassify) GetDeleted() int {
	return m.Deleted
}
