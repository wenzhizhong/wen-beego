package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"
)

var _ itf.MenuItf = (*MchntMenu)(nil)

type MchntMenu struct {
	base_model.UnitMenu
}

func (m *MchntMenu) TableName() string {
	return `mchnt_menu`
}

func (m *MchntMenu) GetId() string {
	return m.Id
}
func (m *MchntMenu) GetUnitId() string {
	return m.UnitId
}
func (m *MchntMenu) GetIcon() string {
	return m.Icon
}
func (m *MchntMenu) GetName() string {
	return m.Name
}
func (m *MchntMenu) GetApiPath() string {
	return m.ApiPath
}
func (m *MchntMenu) GetPagePath() string {
	return m.PagePath
}
func (m *MchntMenu) GetType() int {
	return m.Type
}
func (m *MchntMenu) GetPid() string {
	return m.Pid
}
func (m *MchntMenu) GetAllPid() string {
	return m.AllPid
}
func (m *MchntMenu) GetCreatedAt() time.Time {
	return m.CreatedAt
}
func (m *MchntMenu) GetWeight() int {
	return m.Weight
}
func (m *MchntMenu) GetVisible() int {
	return m.Visible
}
func (m *MchntMenu) GetDeleted() int {
	return m.Deleted
}
