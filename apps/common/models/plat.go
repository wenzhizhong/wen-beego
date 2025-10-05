package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"
)

var _ itf.UnitItf = (*Plat)(nil)

type Plat struct {
	base_model.Unit
}

func (m *Plat) TableName() string {
	return `plat`
}

func (m *Plat) GetId() string {
	return m.Id
}
func (m *Plat) GetPid() string {
	return m.Pid
}
func (m *Plat) GetLogo() string {
	return m.Logo
}
func (m *Plat) GetName() string {
	return m.Name
}
func (m *Plat) GetCode() string {
	return m.Code
}
func (m *Plat) GetCorporation() string {
	return m.Corporation
}
func (m *Plat) GetLicense() string {
	return m.License
}
func (m *Plat) GetAddress() string {
	return m.Address
}
func (m *Plat) GetStatus() int {
	return m.Status
}
func (m *Plat) GetDeleted() int {
	return m.Deleted
}
func (m *Plat) GetCreatedAt() time.Time {
	return m.CreatedAt
}
func (m *Plat) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}
