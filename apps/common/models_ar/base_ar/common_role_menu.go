package base_ar

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"
	"fmt"
)

func GetUserMenu[MenuModel itf.MenuItf, RoleMenuModel itf.RoleMenuItf, UserRoleModel itf.UserRoleItf, RoleModel itf.RoleItf](moduleName string, unitId string, userId string, menuModel MenuModel, roleMenuModel RoleMenuModel, userRoleModel UserRoleModel, roleModel RoleModel) (menuAuthList []dto.RoleMenuDto, err error) {
	if unitId == "" || userId == "" {
		str := fmt.Sprintf("GetUserMenu():获取菜单权限必填参数, unit_id:%s, classifyName:%s", unitId, userId)
		global.Log.Error(str)
		return menuAuthList, errors.New(str)
	}
	isAdmin := helper.IsAdmin(moduleName, unitId, userId)

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
			Where(tableMenu+".visible = ?", 1).
			Where(tableMenu+".deleted = ?", 0).
			Group(tableMenu + ".id").
			Order(tableMenu + ".weight asc").
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
			Where(tableMenu+".visible = ?", 1).
			Where(tableMenu+".deleted = ?", 0).
			Where(tableRole+".deleted = ?", 0).
			Where(tableRole+".status = ?", 1).
			Where(tableUserRole+".user_id = ?", userId).
			Where(tableUserRole+".unit_id = ?", unitId).
			Where(tableUserRole+".deleted = ?", 0).
			Group(tableMenu + ".id").
			Order(tableMenu + ".weight asc").
			Scan(&menuAuthList).
			Error
	}
	if tmpError != nil && !helper.DbNotFound(tmpError) {
		return menuAuthList, tmpError
	}
	return
}

func GetUserPermissions[MenuModel itf.MenuItf, MenuPermsModel itf.MenuPermsItf, RoleMenuModel itf.RoleMenuItf, UserRoleModel itf.UserRoleItf, RoleModel itf.RoleItf](moduleName string, unitId string, userId string, menuModel MenuModel, menuPermsModel MenuPermsModel, roleMenuModel RoleMenuModel, userRoleModel UserRoleModel, roleModel RoleModel) (menuPermsList []base_model.UnitMenuPerms, err error) {
	if unitId == "" || userId == "" {
		str := fmt.Sprintf("GetUserMenu():获取菜单权限必填参数, unit_id:%s, classifyName:%s", unitId, userId)
		global.Log.Error(str)
		return menuPermsList, errors.New(str)
	}
	isAdmin := helper.IsAdmin(moduleName, unitId, userId)

	tableMenuPerms := menuPermsModel.TableName()
	tableRoleMenu := roleMenuModel.TableName()
	tableUserRole := userRoleModel.TableName()
	tableRole := roleModel.TableName()
	tableMenu := menuModel.TableName()
	tableStruct := struct {
		TableMenuPerms string
		TableRoleMenu  string
		TableUserRole  string
		TableRole      string
		TableMenu      string
	}{
		TableMenuPerms: tableMenuPerms,
		TableRoleMenu:  tableRoleMenu,
		TableUserRole:  tableUserRole,
		TableRole:      tableRole,
		TableMenu:      tableMenu,
	}

	selectStr, err := helper.ParseStringTpl(`{{.TableMenuPerms}}.*`, tableStruct)
	joinMenuPermsStr, err2 := helper.ParseStringTpl(`inner join {{.TableMenuPerms}} on {{.TableMenuPerms}}.id = {{.TableRoleMenu}}.menu_perms_id`, tableStruct)
	joinUserRoleStr, err3 := helper.ParseStringTpl(`inner join {{.TableUserRole}} on {{.TableUserRole}}.role_id = {{.TableRoleMenu}}.role_id`, tableStruct)
	joinRoleStr, err4 := helper.ParseStringTpl(`inner join {{.TableRole}} on {{.TableRole}}.id = {{.TableRoleMenu}}.role_id`, tableStruct)
	joinMenuStr, err5 := helper.ParseStringTpl(`inner join {{.TableMenu}} on {{.TableMenu}}.id = {{.TableRoleMenu}}.menu_id`, tableStruct)
	if err != nil {
		return menuPermsList, err
	}
	if err2 != nil {
		return menuPermsList, err2
	}
	if err3 != nil {
		return menuPermsList, err3
	}
	if err4 != nil {
		return menuPermsList, err4
	}
	if err5 != nil {
		return menuPermsList, err5
	}

	var tmpError error
	if isAdmin {
		tmpError = global.GetReadDb().
			Model(menuPermsModel).
			Select(selectStr).
			Joins(joinMenuStr).
			Where(tableMenuPerms+".deleted = 0").
			Where(tableMenu+".deleted = 0").
			Where(tableMenu+".visible = 1").
			Where(tableMenu+".id = ?", unitId).
			Group(tableMenuPerms + ".id").
			Scan(&menuPermsList).
			Error
	} else {
		tmpError = global.GetReadDb().
			Model(roleMenuModel).
			Select(selectStr).
			Joins(joinMenuPermsStr).
			Joins(joinRoleStr).
			Joins(joinUserRoleStr).
			Where(tableMenuPerms+".deleted = 0").
			Where(tableUserRole+".user_id = ?", userId).
			Where(tableUserRole+".unit_id = ?", unitId).
			Where(tableUserRole + ".deleted = 0").
			Where(tableRole + ".status = 1").
			Where(tableRole + ".deleted = 0").
			Group(tableMenuPerms + ".id").
			Scan(&menuPermsList).
			Error
	}
	if tmpError != nil && !helper.DbNotFound(tmpError) {
		return menuPermsList, tmpError
	}
	return
}
