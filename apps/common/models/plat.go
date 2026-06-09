package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.UnitItf = (*Plat)(nil)

type Plat struct {
	base_model.Unit
	IsOfficial bool `json:"isOfficial" gorm:"type:bool;default:false;comment:是否官方平台"`
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
func (m *Plat) GetCreatedAt() int64 {
	return m.CreatedAt
}
func (m *Plat) GetUpdatedAt() int64 {
	return m.UpdatedAt
}

func (m *Plat) GetDeletedAt() *int64 {
	return m.DeletedAt
}
