package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"
)

var _ itf.MenuItf = (*PlatMenu)(nil)

type PlatMenu struct {
	base_model.UnitMenu
}

func (m *PlatMenu) TableName() string {
	return `plat_menu`
}

func (m *PlatMenu) GetId() string {
	return m.Id
}
func (m *PlatMenu) GetUnitId() string {
	return m.UnitId
}
func (m *PlatMenu) GetIcon() string {
	return m.Icon
}
func (m *PlatMenu) GetName() string {
	return m.Name
}
func (m *PlatMenu) GetApiPath() string {
	return m.ApiPath
}
func (m *PlatMenu) GetPagePath() string {
	return m.PagePath
}
func (m *PlatMenu) GetType() int {
	return m.Type
}
func (m *PlatMenu) GetPid() string {
	return m.Pid
}
func (m *PlatMenu) GetAllPid() string {
	return m.AllPid
}
func (m *PlatMenu) GetCreatedAt() time.Time {
	return m.CreatedAt
}
func (m *PlatMenu) GetWeight() int {
	return m.Weight
}
func (m *PlatMenu) GetVisible() int {
	return m.Visible
}
func (m *PlatMenu) GetDeleted() int {
	return m.Deleted
}
