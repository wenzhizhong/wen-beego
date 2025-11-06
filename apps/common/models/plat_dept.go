package models

import "WenBeego/apps/common/models/base_model"

type PlatDept struct {
	base_model.UnitDept
}

func (m *PlatDept) TableName() string {
	return "plat_dept"
}

func (m *PlatDept) GetId() string {
	return m.Id
}
func (m *PlatDept) GetPid() string {
	return m.Pid
}
func (m *PlatDept) GetUnitId() string {
	return m.UnitId
}
func (m *PlatDept) GetName() string {
	return m.Name
}
func (m *PlatDept) GetPrincipalId() string {
	return m.PrincipalId
}
func (m *PlatDept) GetPrincipal() string {
	return m.Principal
}
func (m *PlatDept) GetPhone() string {
	return m.Phone
}
func (m *PlatDept) GetEmail() string {
	return m.Email
}
func (m *PlatDept) GetSort() int {
	return m.Sort
}
func (m *PlatDept) GetStatus() int {
	return m.Status
}
func (m *PlatDept) GetDeleted() int {
	return m.Deleted
}
func (m *PlatDept) GetUpdatedAt() int64 {
	return m.UpdatedAt
}
func (m *PlatDept) GetDeletedAt() *int64 {
	return m.DeletedAt
}
func (m *PlatDept) GetRemark() string {
	return m.Remark
}
