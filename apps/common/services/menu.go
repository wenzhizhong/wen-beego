package services

import (
	"WenBeego/apps/common/ar"
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
)

type CommonMenu struct {
	PlatMenuAr ar.PlatMenuAr
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
func (s *CommonMenu) GetAsyncRoutes(moduleName string, unitId string, userId string) (menuAuthList []dto.RoleMenuDto, err error) {

	var permissions []map[string]interface{}
	var roleClassifies []map[string]interface{}
	if moduleName == "admin_plat" {
		menuAuthList, err = ar.GetUserMenu(moduleName, unitId, userId, &models.PlatMenu{}, &models.PlatRoleMenu{}, &models.PlatUserRole{}, &models.PlatRole{})
		if err != nil && !helper.DbNotFound(err) {
			return
		}
		permissions, err = ar.GetUserPermissions(moduleName, unitId, userId, &models.PlatMenu{}, &models.PlatMenuPerms{}, &models.PlatRoleMenu{}, &models.PlatUserRole{}, &models.PlatRole{})
		if err != nil && !helper.DbNotFound(err) {
			return
		}
		roleClassifies, err = ar.GetUserRoleClassifies(unitId, userId, &models.Plat{}, &models.PlatRole{}, &models.PlatRoleClassify{}, &models.PlatUserRole{})
	} else {
		menuAuthList, err = ar.GetUserMenu(moduleName, unitId, userId, &models.MchntMenu{}, &models.MchntRoleMenu{}, &models.MchntUserRole{}, &models.MchntRole{})
		if err != nil && !helper.DbNotFound(err) {
			return
		}
		permissions, err = ar.GetUserPermissions(moduleName, unitId, userId, &models.MchntMenu{}, &models.MchntMenuPerms{}, &models.MchntRoleMenu{}, &models.MchntUserRole{}, &models.MchntRole{})
		if err != nil && !helper.DbNotFound(err) {
			return
		}
		roleClassifies, err = ar.GetUserRoleClassifies(unitId, userId, &models.Mchnt{}, &models.MchntRole{}, &models.MchntRoleClassify{}, &models.MchntUserRole{})
	}

	if err != nil && !helper.DbNotFound(err) {
		return
	}

	if len(menuAuthList) > 0 {
		permissionsMap := make(map[string]map[string]string)
		roleClassifiesMap := make(map[string]string)
		for _, permission := range permissions {
			menuId := permission["menu_id"].(string)
			perm := permission["permission"].(string)
			if _, exists := permissionsMap[menuId]; !exists {
				permissionsMap[menuId] = make(map[string]string)
			}
			permissionsMap[menuId][perm] = perm
		}
		for _, roleClassify := range roleClassifies {
			roleClassifiesMap[roleClassify["name"].(string)] = roleClassify["name"].(string)
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
