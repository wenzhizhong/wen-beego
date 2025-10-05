package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.UnitUserItf = (*MchntUser)(nil)

type MchntUser struct {
	base_model.UnitUser
}

func (m *MchntUser) TableName() string {
	return `mchnt_user`
}

func (m *MchntUser) GetId() string {
	return m.Id
}
func (m *MchntUser) GetUnitId() string {
	return m.UnitId
}
func (m *MchntUser) GetIsDefault() int {
	return m.IsDefault
}
func (m *MchntUser) GetDefaultUnitId() string {
	return m.DefaultUnitId
}
func (m *MchntUser) GetUserId() string {
	return m.UserId
}
func (m *MchntUser) GetDeleted() int {
	return m.Deleted
}
func (m *MchntUser) GetPhone() string {
	return m.Phone
}
func (m *MchntUser) GetName() string {
	return m.Name
}
