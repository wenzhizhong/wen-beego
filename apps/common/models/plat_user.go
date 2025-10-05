package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.UnitUserItf = (*PlatUser)(nil)

type PlatUser struct {
	base_model.UnitUser
}

func (m *PlatUser) TableName() string {
	return `plat_user`
}

func (m *PlatUser) GetId() string {
	return m.Id
}
func (m *PlatUser) GetUnitId() string {
	return m.UnitId
}
func (m *PlatUser) GetIsDefault() int {
	return m.IsDefault
}
func (m *PlatUser) GetDefaultUnitId() string {
	return m.DefaultUnitId
}
func (m *PlatUser) GetUserId() string {
	return m.UserId
}
func (m *PlatUser) GetDeleted() int {
	return m.Deleted
}
func (m *PlatUser) GetPhone() string {
	return m.Phone
}
func (m *PlatUser) GetName() string {
	return m.Name
}
