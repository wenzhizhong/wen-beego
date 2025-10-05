package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"
)

var _ itf.RoleItf = (*PlatRole)(nil)

type PlatRole struct {
	base_model.UnitRole
}

func (m *PlatRole) TableName() string {
	return `plat_role`
}

func (m *PlatRole) GetId() string {
	return m.Id
}
func (m *PlatRole) GetUnitId() string {
	return m.UnitId
}
func (m *PlatRole) GetRoleName() string {
	return m.RoleName
}
func (m *PlatRole) GetRoleSort() int {
	return m.RoleSort
}
func (m *PlatRole) GetStatus() int {
	return m.Status
}
func (m *PlatRole) GetDeleted() int {
	return m.Deleted
}
func (m *PlatRole) GetCreatedBy() string {
	return m.CreatedBy
}
func (m *PlatRole) GetCreatedAt() time.Time {
	return m.CreatedAt
}
func (m *PlatRole) GetUpdatedBy() string {
	return m.UpdatedBy
}
func (m *PlatRole) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}
func (m *PlatRole) GetRemark() string {
	return m.Remark
}
