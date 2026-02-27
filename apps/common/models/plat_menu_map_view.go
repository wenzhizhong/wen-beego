package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

type PlatMenuMapView struct {
	base_model.UnitMenuMap
	MenuFrom string `json:"menuFrom" gorm:"->"`
}

var _ itf.MenuMapItf = (*PlatMenuMapView)(nil)

func (m *PlatMenuMapView) TableName() string {
	return "plat_menu_map_view"
}

func (m *PlatMenuMapView) GetId() string {
	return m.Id
}
func (m *PlatMenuMapView) GetUnitId() string {
	return m.UnitId
}
func (m *PlatMenuMapView) GetMenuId() string {
	return m.MenuId
}
func (m *PlatMenuMapView) GetUpdatedAt() int64 {
	return m.UpdatedAt
}
func (m *PlatMenuMapView) GetDeleted() int {
	return m.Deleted
}
