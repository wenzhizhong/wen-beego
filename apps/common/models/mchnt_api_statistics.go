package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.ApiStatisticsItf = (*MchntApiStatistics)(nil)

type MchntApiStatistics struct {
	base_model.UnitApiStatistics
}

// TableName 指定表名
func (m *MchntApiStatistics) TableName() string {
	return "mchnt_api_statistics"
}

func (m *MchntApiStatistics) GetID() string {
	return m.ID
}
func (m *MchntApiStatistics) GetPermsID() string {
	return m.PermsID
}
func (m *MchntApiStatistics) GetURI() string {
	return m.URI
}
func (m *MchntApiStatistics) GetPV() int {
	return m.PV
}
func (m *MchntApiStatistics) GetUV() int {
	return m.UV
}
func (m *MchntApiStatistics) GetDate() int64 {
	return m.Date
}
func (m *MchntApiStatistics) GetUnitId() string {
	return m.UnitId
}
func (m *MchntApiStatistics) GetModuleName() string {
	return m.ModuleName
}
