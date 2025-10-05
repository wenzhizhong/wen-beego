package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.RoleClassifyItf = (*PlatRoleClassify)(nil)

type PlatRoleClassify struct {
	base_model.UnitRoleClassify
}

func (m *PlatRoleClassify) TableName() string {
	return `plat_role_classify`
}

func (m *PlatRoleClassify) GetId() string {
	return m.Id
}
func (m *PlatRoleClassify) GetRoleId() string {
	return m.RoleId
}
func (m *PlatRoleClassify) GetUnitId() string {
	return m.UnitId
}
func (m *PlatRoleClassify) GetName() string {
	return m.Name
}
func (m *PlatRoleClassify) GetDeleted() int {
	return m.Deleted
}
