package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.UnitItf = (*Mchnt)(nil)

type Mchnt struct {
	base_model.Unit
	PlatStatus int `json:"plat_status" gorm:"type:int4;not null;default:0;comment:平台审核状态（0未审核，1审核通过，2审核不通过，3禁用）"`
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
