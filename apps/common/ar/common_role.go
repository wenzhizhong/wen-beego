package ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
	"fmt"
)

func GetUserRole[UserRoleModel models.ModelInterface, RoleModel models.ModelInterface](moduleName string, unitId string, userId string, userRoleModel UserRoleModel, roleModel RoleModel) (rolesList []map[string]interface{}, err error) {
	if unitId == "" || userId == "" {
		str := fmt.Sprintf("GetUserMenu():获取菜单权限必填参数, unit_id:%s, classifyName:%s", unitId, userId)
		global.Log.Error(str)
		return rolesList, errors.New(str)
	}
	// TODO: 获取用户角色
	return
}
