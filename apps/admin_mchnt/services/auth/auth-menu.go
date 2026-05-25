package auth

import (
	"WenBeego/apps/common/dto_vo/auth_vo"
	"WenBeego/apps/common/services/auth"
)

type MenuService struct {
	commonMenu auth.CommonMenu
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
func (s *MenuService) GetAsyncRoutes(moduleName string, unitId string, unitUserId string) (menuAuthList []auth_vo.RoleMenuVo, err error) {
	return s.commonMenu.GetAsyncRoutes(moduleName, unitId, unitUserId)
}
