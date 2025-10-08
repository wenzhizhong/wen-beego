package auth

import (
	"WenBeego/apps/common/dto/auth_dto"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models_ar"
	"WenBeego/apps/common/models_ar/base_ar"
	"errors"
)

type CommonMenu struct {
	PlatMenuAr models_ar.PlatMenuAr
}

/**
 * 获取用户菜单权限
 * @param moduleName 模块名称
 * @param unitId 单位ID
 * @param userId 用户ID
 * @param menuModel 菜单模型
 * @param roleMenuModel 角色菜单模型
 * @param userRoleModel 用户角色模型
 * @param roleModel 角色模型
 * @return
 */
func (s *CommonMenu) GetAsyncRoutes(moduleName string, unitId string, userId string) (menuAuthList []auth_dto.RoleMenuDto, err error) {

	var permissions []base_model.UnitMenuPerms
	var roleClassifies []base_model.UnitRoleClassify
	switch moduleName {
	case "admin_plat":
		menuAuthList, err = base_ar.GetUserMenu(moduleName, unitId, userId, &models.PlatMenu{}, &models.PlatRoleMenu{}, &models.PlatUserRole{}, &models.PlatRole{})
		if err != nil && !helper.DbNotFound(err) {
			return
		}
		permissions, err = base_ar.GetUserPermissions(moduleName, unitId, userId, &models.PlatMenu{}, &models.PlatMenuPerms{}, &models.PlatRoleMenu{}, &models.PlatUserRole{}, &models.PlatRole{})
		if err != nil && !helper.DbNotFound(err) {
			return
		}
		roleClassifies, err = base_ar.GetUserRoleClassifies(unitId, userId, &models.Plat{}, &models.PlatRole{}, &models.PlatRoleClassify{}, &models.PlatUserRole{})
	case "admin_mchnt":
		menuAuthList, err = base_ar.GetUserMenu(moduleName, unitId, userId, &models.MchntMenu{}, &models.MchntRoleMenu{}, &models.MchntUserRole{}, &models.MchntRole{})
		if err != nil && !helper.DbNotFound(err) {
			return
		}
		permissions, err = base_ar.GetUserPermissions(moduleName, unitId, userId, &models.MchntMenu{}, &models.MchntMenuPerms{}, &models.MchntRoleMenu{}, &models.MchntUserRole{}, &models.MchntRole{})
		if err != nil && !helper.DbNotFound(err) {
			return
		}
		roleClassifies, err = base_ar.GetUserRoleClassifies(unitId, userId, &models.Mchnt{}, &models.MchntRole{}, &models.MchntRoleClassify{}, &models.MchntUserRole{})
	default:
		err = errors.New("未知的模块名称")
	}

	if err != nil && !helper.DbNotFound(err) {
		return
	}

	if len(menuAuthList) > 0 {
		permissionsMap := make(map[string]map[string]string)
		roleClassifiesMap := make(map[string]string)
		for _, permission := range permissions {
			menuId := permission.MenuId
			perm := permission.Permission
			if _, exists := permissionsMap[menuId]; !exists {
				permissionsMap[menuId] = make(map[string]string)
			}
			permissionsMap[menuId][perm] = perm
		}
		for _, roleClassify := range roleClassifies {
			roleClassifiesMap[roleClassify.Name] = roleClassify.Name
		}
		roleClassifiesKeys := helper.GetMapKeys(roleClassifiesMap)
		for i, menu := range menuAuthList {
			menuAuthList[i].Auths = make([]string, 0)
			if permissionsMap[menu.Id] != nil {
				permissionsKeys := helper.GetMapKeys(permissionsMap[menu.Id])
				menuAuthList[i].Auths = append(menuAuthList[i].Auths, permissionsKeys...)
			}
			menuAuthList[i].Roles = append(menuAuthList[i].Roles, roleClassifiesKeys...)
		}
	}

	return
}
