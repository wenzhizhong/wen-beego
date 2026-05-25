package models_ar

import (
	"WenBeego/apps/common/dto_vo/auth_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type PlatRoleMenu struct {
	models.MchntRoleMenu
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
