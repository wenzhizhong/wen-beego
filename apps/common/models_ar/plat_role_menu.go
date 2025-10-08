package models_ar

import (
	"WenBeego/apps/common/dto/auth_dto"
	"WenBeego/apps/common/global"
	"errors"
)

type PlatRoleMenu struct {
	Id          string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	RoleId      string `json:"role_id" gorm:"type:bpchar(36);not null;comment:角色ID"`
	MenuId      string `json:"menu_id" gorm:"type:bpchar(36);not null;comment:菜单权限ID"`
	MenuPermsId string `json:"menu_perms_id" gorm:"type:varchar(36);not null;comment:资源关联菜单"`
}

func (a *PlatRoleMenu) TableName() string {
	return `plat_role_menu`
}

func (a *PlatRoleMenu) GetById(id string) (PlatRoleMenu, error) {
	roleMenu := PlatRoleMenu{}
	if id == "" {
		return roleMenu, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&roleMenu)
	return roleMenu, result.Error
}

/**
 * 获取用户权限菜单
 */
func (a *PlatRoleMenu) GetUserMenu(unitId string, userId string) (roleMenu []auth_dto.RoleMenuDto, err error) {

	return
}
