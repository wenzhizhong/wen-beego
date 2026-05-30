package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.ApiStatisticsItf = (*PlatApiStatistics)(nil)

type PlatApiStatistics struct {
	base_model.UnitApiStatistics
}

// TableName 指定表名
func (m *PlatApiStatistics) TableName() string {
	return "plat_api_statistics"
}

func (m *PlatApiStatistics) GetID() string {
	return m.ID
}
func (m *PlatApiStatistics) GetPermsID() string {
	return m.PermsID
}
func (m *PlatApiStatistics) GetURI() string {
	return m.URI
}
func (m *PlatApiStatistics) GetPV() int {
	return m.PV
}
func (m *PlatApiStatistics) GetUV() int {
	return m.UV
}
func (m *PlatApiStatistics) GetDate() int64 {
	return m.Date
}

func (m *PlatApiStatistics) GetUnitId() string {
	return m.UnitId
}

func (m *PlatApiStatistics) GetModuleName() string {
	return m.ModuleName
}
