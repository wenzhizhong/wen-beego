package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

type PlatMenuMap struct {
	base_model.UnitMenuMap
}

var _ itf.MenuMapItf = (*PlatMenuMap)(nil)

func (m *PlatMenuMap) TableName() string {
	return "plat_menu_map"
}

func (m *PlatMenuMap) GetId() string {
	return m.Id
}
func (m *PlatMenuMap) GetUnitId() string {
	return m.UnitId
}
func (m *PlatMenuMap) GetMenuId() string {
	return m.MenuId
}
func (m *PlatMenuMap) GetUpdatedAt() int64 {
	return m.UpdatedAt
}
func (m *PlatMenuMap) GetDeleted() int {
	return m.Deleted
}
