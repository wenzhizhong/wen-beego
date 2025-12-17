package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"
)

var _ itf.UnitCronItf = (*PlatCron)(nil)

type PlatCron struct {
	base_model.UnitCron
}

func (m *PlatCron) TableName() string {
	return `plat_crontab`
}

func (m *PlatCron) GetID() string {
	return m.Id
}
func (m *PlatCron) GetUnitID() string {
	return m.UnitId
}
func (m *PlatCron) GetName() string {
	return m.Name
}
func (m *PlatCron) GetNameEn() string {
	return m.NameEn
}
func (m *PlatCron) GetGroup() string {
	return m.Group
}
func (m *PlatCron) GetCronExpr() string {
	return m.CronExpr
}
func (m *PlatCron) GetStatus() int {
	return m.Status
}
func (m *PlatCron) GetCreatedBy() string {
	return m.CreatedBy
}
func (m *PlatCron) GetCreatedAt() *time.Time {
	return m.CreatedAt
}
func (m *PlatCron) GetUpdatedBy() *string {
	return m.UpdatedBy
}
func (m *PlatCron) GetUpdatedAt() *time.Time {
	return m.UpdatedAt
}
func (m *PlatCron) GetDeleted() int {
	return m.Deleted
}

func (m *PlatCron) GetRemark() string {
	return m.Remark
}
