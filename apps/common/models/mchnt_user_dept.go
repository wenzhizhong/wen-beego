package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.UserDeptItf = (*MchntUserDept)(nil)

type MchntUserDept struct {
	base_model.UnitUserDept
}

func (m *MchntUserDept) TableName() string {
	return `mchnt_user_dept`
}

func (m *MchntUserDept) GetId() string {
	return m.Id
}
func (m *MchntUserDept) GetUserId() string {
	return m.UserId
}
func (m *MchntUserDept) GetDeptId() string {
	return m.DeptId
}
func (m *MchntUserDept) GetDeleted() int {
	return m.Deleted
}
