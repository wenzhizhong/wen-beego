package base_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"
	"fmt"
)

/*
* 获取用户角色
* @param moduleName 模块名称
* @param unitId 组织单位id
* @param userId 用户id
* @return rolesList 角色列表
 */
func GetUserRole[UserRoleModel itf.UserRoleItf, RoleModel itf.RoleItf](moduleName string, unitId string, userId string, userRoleModel UserRoleModel, roleModel RoleModel) (rolesList []base_model.UnitRole, err error) {
	if unitId == "" || userId == "" {
		str := fmt.Sprintf("GetUserMenu():获取菜单权限必填参数, unit_id:%s, classifyName:%s", unitId, userId)
		global.Log.Error(str)
		return rolesList, errors.New(str)
	}
	tableRole := roleModel.TableName()
	tableUserRole := userRoleModel.TableName()
	tableStruct := struct {
		TableRole     string
		TableUserRole string
	}{
		TableRole:     tableRole,
		TableUserRole: tableUserRole,
	}
	selectStr, err := helper.ParseStringTpl(`{{.TableRole}}.*`, tableStruct)
	joinUserRoleStr, err2 := helper.ParseStringTpl(`inner join {{.TableUserRole}} on {{.TableUserRole}}.role_id = {{.TableRole}}.id`, tableStruct)
	if err != nil {
		return rolesList, err
	}
	if err2 != nil {
		return rolesList, err2
	}

	tmpError := global.GetReadDb().
		Model(roleModel).
		Select(selectStr).
		Joins(joinUserRoleStr).
		Where(tableUserRole+".user_id = ?", userId).
		Where(tableUserRole+".unit_id = ?", unitId).
		Where(tableUserRole + ".deleted = 0").
		Scan(&rolesList).
		Error
	if helper.DbNotFound(tmpError) {
		return rolesList, nil
	}
	return
}
