package base_ar

import (
	"WenBeego/apps/common/dto/auth_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"
	"fmt"
)

func GetUserMenu[MenuModel itf.MenuItf, RoleMenuModel itf.RoleMenuItf, UserRoleModel itf.UserRoleItf, RoleModel itf.RoleItf](moduleName string, unitId string, unitUserId string, menuModel MenuModel, roleMenuModel RoleMenuModel, userRoleModel UserRoleModel, roleModel RoleModel) (menuAuthList []auth_dto.RoleMenuDto, err error) {
	if unitId == "" || unitUserId == "" {
		str := fmt.Sprintf("GetUserMenu():获取菜单权限必填参数, unit_id:%s, unitUserId:%s", unitId, unitUserId)
		global.Log.Error(str)
		return menuAuthList, errors.New(str)
	}
	isAdmin := helper.IsAdmin(moduleName, unitUserId)

	tableMenu := menuModel.TableName()
	tableRoleMenu := roleMenuModel.TableName()
	tableUserRole := userRoleModel.TableName()
	tableRole := roleModel.TableName()
	tableStruct := struct {
		TableMenu     string
		TableRoleMenu string
		TableUserRole string
		TableRole     string
	}{
		TableMenu:     tableMenu,
		TableRoleMenu: tableRoleMenu,
		TableUserRole: tableUserRole,
		TableRole:     tableRole,
	}

	selectStr, err := helper.ParseStringTpl(`{{.TableMenu}}.*`, tableStruct)
	joinMenuStr, err2 := helper.ParseStringTpl(`inner join {{.TableMenu}} on {{.TableMenu}}.id = {{.TableRoleMenu}}.menu_id`, tableStruct)
	joinUserRoleStr, err3 := helper.ParseStringTpl(`inner join {{.TableUserRole}} on {{.TableUserRole}}.role_id = {{.TableRoleMenu}}.role_id`, tableStruct)
	joinRoleStr, err4 := helper.ParseStringTpl(`inner join {{.TableRole}} on {{.TableRole}}.id = {{.TableRoleMenu}}.role_id`, tableStruct)
	if err != nil {
		return menuAuthList, err
	}
	if err2 != nil {
		return menuAuthList, err2
	}
	if err3 != nil {
		return menuAuthList, err3
	}
	if err4 != nil {
		return menuAuthList, err4
	}

	var tmpError error
	if isAdmin {
		tmpError = global.GetReadDb().
			Model(menuModel).
			Select(selectStr).
			Where(tableMenu+".unit_id = ?", unitId).
			Where(tableMenu+".deleted = ?", 0).
			Where(tableMenu+".menu_type <> ?", 3).
			Group(tableMenu + ".id").
			Order(tableMenu + ".rank asc").
			Scan(&menuAuthList).
			Error
	} else {
		tmpError = global.GetReadDb().
			Model(roleMenuModel).
			Select(selectStr).
			Joins(joinMenuStr).
			Joins(joinRoleStr).
			Joins(joinUserRoleStr).
			Where(tableMenu+".unit_id = ?", unitId).
			Where(tableMenu+".deleted = ?", 0).
			Where(tableMenu+".menu_type <> ?", 3).
			Where(tableRole+".deleted = ?", 0).
			Where(tableRole+".status = ?", 1).
			Where(tableUserRole+".user_id = ?", unitUserId).
			Where(tableUserRole+".deleted = ?", 0).
			Group(tableMenu + ".id").
			Order(tableMenu + ".rank asc").
			Scan(&menuAuthList).
			Error
	}
	if tmpError != nil && !helper.DbNotFound(tmpError) {
		return menuAuthList, tmpError
	}
	return
}

func GetUserPermissions[MenuModel itf.MenuItf, RoleMenuModel itf.RoleMenuItf, UserRoleModel itf.UserRoleItf, RoleModel itf.RoleItf](moduleName string, unitId string, unitUserId string, menuModel MenuModel, roleMenuModel RoleMenuModel, userRoleModel UserRoleModel, roleModel RoleModel) (menuList []base_model.UnitMenu, err error) {
	if unitUserId == "" {
		str := fmt.Sprintf("GetUserMenu():获取菜单权限必填参数, unit_id:%s, classifyName:%s", unitId, unitUserId)
		global.Log.Error(str)
		return menuList, errors.New(str)
	}
	isAdmin := helper.IsAdmin(moduleName, unitUserId)

	tableMenu := menuModel.TableName()
	tableRoleMenu := roleMenuModel.TableName()
	tableUserRole := userRoleModel.TableName()
	tableRole := roleModel.TableName()
	tableStruct := struct {
		TableMenu     string
		TableRoleMenu string
		TableUserRole string
		TableRole     string
	}{
		TableMenu:     tableMenu,
		TableRoleMenu: tableRoleMenu,
		TableUserRole: tableUserRole,
		TableRole:     tableRole,
	}

	selectStr, err := helper.ParseStringTpl(`{{.TableMenu}}.*`, tableStruct)
	joinMenuStr, err2 := helper.ParseStringTpl(`inner join {{.TableMenu}} on {{.TableMenu}}.id = {{.TableRoleMenu}}.menu_id`, tableStruct)
	joinUserRoleStr, err3 := helper.ParseStringTpl(`inner join {{.TableUserRole}} on {{.TableUserRole}}.role_id = {{.TableRoleMenu}}.role_id`, tableStruct)
	joinRoleStr, err4 := helper.ParseStringTpl(`inner join {{.TableRole}} on {{.TableRole}}.id = {{.TableRoleMenu}}.role_id`, tableStruct)

	if err != nil {
		return menuList, err
	}
	if err2 != nil {
		return menuList, err2
	}
	if err3 != nil {
		return menuList, err3
	}
	if err4 != nil {
		return menuList, err4
	}

	var tmpError error
	if isAdmin {
		tmpError = global.GetReadDb().
			Model(menuModel).
			Select(selectStr).
			Where(tableMenu+".deleted = 0").
			Where(tableMenu+".menu_type = ?", 3).
			Where(tableMenu+".unit_id = ?", unitId).
			Group(tableMenu + ".id").
			Scan(&menuList).
			Error
	} else {
		tmpError = global.GetReadDb().
			Model(roleMenuModel).
			Select(selectStr).
			Joins(joinMenuStr).
			Joins(joinRoleStr).
			Joins(joinUserRoleStr).
			Where(tableMenu+".deleted = 0").
			Where(tableMenu+".menu_type = ?", 3).
			Where(tableUserRole+".user_id = ?", unitUserId).
			Where(tableUserRole + ".deleted = 0").
			Where(tableRole + ".status = 1").
			Where(tableRole + ".deleted = 0").
			Group(tableMenu + ".id").
			Scan(&menuList).
			Error
	}
	if tmpError != nil && !helper.DbNotFound(tmpError) {
		return menuList, tmpError
	}
	return
}

// 获取单位权限列表
func GetUnitPermissions[MenuModel itf.MenuItf](unitId string, menuModel MenuModel) (menuList []base_model.UnitMenu, err error) {
	tableMenuName := menuModel.TableName()
	err = global.GetReadDb().
		Model(menuModel).
		Select(tableMenuName+".*").
		Where(tableMenuName+".unit_id = ?", unitId).
		Where(tableMenuName+".deleted = ?", 0).
		Where(tableMenuName+".menu_type = ?", 3).
		Find(&menuList).Error
	return
}
