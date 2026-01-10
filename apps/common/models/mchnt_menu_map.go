package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

type MchntMenuMap struct {
	base_model.UnitMenuMap
}

var _ itf.MenuMapItf = (*MchntMenuMap)(nil)

func (m *MchntMenuMap) TableName() string {
	return "mchnt_menu_map"
}

func (m *MchntMenuMap) GetId() string {
	return m.Id
}
func (m *MchntMenuMap) GetUnitId() string {
	return m.UnitId
}
func (m *MchntMenuMap) GetMenuId() string {
	return m.MenuId
}
func (m *MchntMenuMap) GetUpdatedAt() int64 {
	return m.UpdatedAt
}
func (m *MchntMenuMap) GetDeleted() int {
	return m.Deleted
}
