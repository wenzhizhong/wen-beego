package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

type MchntDept struct {
	base_model.UnitDept
}

var _ itf.DeptItf = (*MchntDept)(nil)

func (m *MchntDept) TableName() string {
	return "mchnt_dept"
}

func (m *MchntDept) GetId() string {
	return m.Id
}
func (m *MchntDept) GetPid() string {
	return m.Pid
}
func (m *MchntDept) GetUnitId() string {
	return m.UnitId
}
func (m *MchntDept) GetName() string {
	return m.Name
}
func (m *MchntDept) GetPrincipalId() string {
	return m.PrincipalId
}
func (m *MchntDept) GetSort() int {
	return m.Sort
}
func (m *MchntDept) GetStatus() int {
	return m.Status
}
func (m *MchntDept) GetDeleted() int {
	return m.Deleted
}
func (m *MchntDept) GetUpdatedAt() int64 {
	return m.UpdatedAt
}
func (m *MchntDept) GetDeletedAt() *int64 {
	return m.DeletedAt
}
func (m *MchntDept) GetRemark() string {
	return m.Remark
}
