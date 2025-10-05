package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"
)

var _ itf.MenuPermsItf = (*MchntMenuPerms)(nil)

type MchntMenuPerms struct {
	base_model.UnitMenuPerms
}

func (m *MchntMenuPerms) TableName() string {
	return `mchnt_menu_perms`
}
func (m *MchntMenuPerms) GetId() string {
	return m.Id
}
func (m *MchntMenuPerms) GetMenuId() string {
	return m.MenuId
}
func (m *MchntMenuPerms) GetType() int16 {
	return m.Type
}
func (m *MchntMenuPerms) GetName() string {
	return m.Name
}
func (m *MchntMenuPerms) GetPermission() string {
	return m.Permission
}
func (m *MchntMenuPerms) GetUri() string {
	return m.Uri

}
func (m *MchntMenuPerms) GetMethod() int16 {
	return m.Method
}
func (m *MchntMenuPerms) GetDeleted() int16 {
	return m.Deleted
}
func (m *MchntMenuPerms) GetCreatedAt() time.Time {
	return m.CreatedAt
}
func (m *MchntMenuPerms) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}
