package base_ar

import (
	"WenBeego/apps/common/dto/user_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

/*
* 身份认证-获取用户角色
* @param moduleName 模块名称
* @param unitId 组织单位id
* @param userId 用户id
* @return rolesList 角色列表
 */
func GetUserRole[UserRoleModel itf.UserRoleItf, RoleModel itf.RoleItf](moduleName string, unitUserId string, userRoleModel UserRoleModel, roleModel RoleModel) (rolesList []base_model.UnitRole, err error) {
	if unitUserId == "" {
		str := fmt.Sprintf("GetUserMenu():获取菜单权限必填参数, unitUserId:%s", unitUserId)
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
		Where(tableUserRole+".user_id = ?", unitUserId).
		Where(tableUserRole + ".deleted = 0").
		Scan(&rolesList).
		Error
	if helper.DbNotFound(tmpError) {
		return rolesList, nil
	}
	return
}

/**
 * 新增组织单位-初始化新增组织单位用户角色
 */
func InsertUnitUserRole[UserRoleModel itf.UserRoleItf](tx *gorm.DB, userRoleModel base_model.UnitUserRole) (err error) {
	if userRoleModel.Id == "" {
		return errors.New("新增用户角色，用户角色id不能为空")
	}
	var tmpUnitUserRole UserRoleModel
	err = global.GetReadDb().
		Model(tmpUnitUserRole).
		Where("id = ?", userRoleModel.Id).
		Take(&tmpUnitUserRole).Error
	if err == nil && tmpUnitUserRole.GetId() != "" {
		return nil
	}

	err = tx.Model(tmpUnitUserRole).
		Create(&userRoleModel).Error
	return
}

func SaveUnitUserRole[UserRoleModel itf.UserRoleItf](tx *gorm.DB, userRoleData []user_dto.UnitUserRoleDto, userRoleModel UserRoleModel) (err error) {
	if len(userRoleData) == 0 {
		err = errors.New("新增用户角色，用户角色数据不能为空")
		return
	}
	userIds := make([]string, 0)
	roleIds := make([]string, 0)
	for k, v := range userRoleData {
		if v.Id == "" {
			userRoleData[k].Id, err = helper.GetUuid()
			if err != nil {
				return
			}
		}
		if v.UserId == "" || v.RoleId == "" {
			return errors.New("新增用户角色，用户id和角色id不能为空")
		}
		userIds = append(userIds, v.UserId)
		roleIds = append(roleIds, v.RoleId)
	}
	err = tx.Where("user_id IN ? AND role_id IN ? ", userIds, roleIds).Delete(&userRoleModel).Error
	if err != nil {
		return
	}
	err = tx.Model(userRoleModel).
		Create(&userRoleData).Error
	return
}
