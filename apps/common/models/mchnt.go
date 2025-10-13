package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.UnitItf = (*Mchnt)(nil)

type Mchnt struct {
	base_model.Unit
}

func (m *Mchnt) TableName() string {
	return `mchnt`
}

func (m *Mchnt) GetId() string {
	return m.Id
}
func (m *Mchnt) GetPid() string {
	return m.Pid
}
func (m *Mchnt) GetLogo() string {
	return m.Logo
}
func (m *Mchnt) GetName() string {
	return m.Name
}
func (m *Mchnt) GetCode() string {
	return m.Code
}
func (m *Mchnt) GetCorporation() string {
	return m.Corporation
}
func (m *Mchnt) GetLicense() string {
	return m.License
}
func (m *Mchnt) GetAddress() string {
	return m.Address
}
func (m *Mchnt) GetStatus() int {
	return m.Status
}
func (m *Mchnt) GetDeleted() int {
	return m.Deleted
}
func (m *Mchnt) GetCreatedAt() int64 {
	return m.CreatedAt
}
func (m *Mchnt) GetUpdatedAt() int64 {
	return m.UpdatedAt
}

func (m *Mchnt) GetDeletedAt() *int64 {
	return m.DeletedAt
}
