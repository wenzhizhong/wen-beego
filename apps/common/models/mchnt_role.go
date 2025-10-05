package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"
)

var _ itf.RoleItf = (*MchntRole)(nil)

type MchntRole struct {
	base_model.UnitRole
}

func (m *MchntRole) TableName() string {
	return `mchnt_role`
}

func (m *MchntRole) GetId() string {
	return m.Id
}
func (m *MchntRole) GetUnitId() string {
	return m.UnitId
}
func (m *MchntRole) GetRoleName() string {
	return m.RoleName
}
func (m *MchntRole) GetRoleSort() int {
	return m.RoleSort
}
func (m *MchntRole) GetStatus() int {
	return m.Status
}
func (m *MchntRole) GetDeleted() int {
	return m.Deleted
}
func (m *MchntRole) GetCreatedBy() string {
	return m.CreatedBy
}
func (m *MchntRole) GetCreatedAt() time.Time {
	return m.CreatedAt
}
func (m *MchntRole) GetUpdatedBy() string {
	return m.UpdatedBy
}
func (m *MchntRole) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}
func (m *MchntRole) GetRemark() string {
	return m.Remark
}
