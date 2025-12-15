package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"
)

var _ itf.UnitCronItf = (*MchntCron)(nil)

type MchntCron struct {
	base_model.UnitCron
}

func (m *MchntCron) TableName() string {
	return `mchnt_crontab`
}

func (m *MchntCron) GetID() string {
	return m.Id
}
func (m *MchntCron) GetUnitID() string {
	return m.UnitId
}
func (m *MchntCron) GetName() string {
	return m.Name
}
func (m *MchntCron) GetNameEn() string {
	return m.NameEn
}
func (m *MchntCron) GetGroup() string {
	return m.Group
}
func (m *MchntCron) GetCronExpr() string {
	return m.CronExpr
}
func (m *MchntCron) GetStatus() int {
	return m.Status
}
func (m *MchntCron) GetCreatedBy() string {
	return m.CreatedBy
}
func (m *MchntCron) GetCreatedAt() *time.Time {
	return m.CreatedAt
}
func (m *MchntCron) GetUpdatedBy() *string {
	return m.UpdatedBy
}
func (m *MchntCron) GetUpdatedAt() *time.Time {
	return m.UpdatedAt
}
func (m *MchntCron) GetDeleted() int {
	return m.Deleted
}
