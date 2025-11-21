package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
)

var _ itf.UserDeptItf = (*PlattUserDept)(nil)

type PlattUserDept struct {
	base_model.UnitUserDept
}

func (m *PlattUserDept) TableName() string {
	return `plat_user_dept`
}

func (m *PlattUserDept) GetId() string {
	return m.Id
}
func (m *PlattUserDept) GetUserId() string {
	return m.UserId
}
func (m *PlattUserDept) GetDeptId() string {
	return m.DeptId
}
func (m *PlattUserDept) GetDeleted() int {
	return m.Deleted
}
