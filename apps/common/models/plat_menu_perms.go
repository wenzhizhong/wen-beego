package models

import (
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"time"
)

var _ itf.MenuPermsItf = (*PlatMenuPerms)(nil)

type PlatMenuPerms struct {
	base_model.UnitMenuPerms
}

func (m *PlatMenuPerms) TableName() string {
	return `plat_menu_perms`
}

func (m *PlatMenuPerms) GetId() string {
	return m.Id
}
func (m *PlatMenuPerms) GetMenuId() string {
	return m.MenuId
}
func (m *PlatMenuPerms) GetType() int16 {
	return m.Type
}
func (m *PlatMenuPerms) GetName() string {
	return m.Name
}
func (m *PlatMenuPerms) GetPermission() string {
	return m.Permission
}
func (m *PlatMenuPerms) GetUri() string {
	return m.Uri

}
func (m *PlatMenuPerms) GetMethod() int16 {
	return m.Method
}
func (m *PlatMenuPerms) GetDeleted() int16 {
	return m.Deleted
}
func (m *PlatMenuPerms) GetCreatedAt() time.Time {
	return m.CreatedAt
}
func (m *PlatMenuPerms) GetUpdatedAt() time.Time {
	return m.UpdatedAt
}
