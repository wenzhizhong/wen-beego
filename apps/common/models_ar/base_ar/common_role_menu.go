package base_ar

import (
	"WenBeego/apps/common/dto/auth_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

var operateMenuTypeArr = []int{base_model.UNIT_MENU_TYPE_BUTTON, base_model.UNIT_MENU_TYPE_OTHER_API}

func GetUserMenu[MenuModel itf.MenuItf, MenuMapModel itf.MenuMapItf, RoleMenuModel itf.RoleMenuItf, UserRoleModel itf.UserRoleItf, RoleModel itf.RoleItf](moduleName string, unitId string, unitUserId string, menuModel MenuModel, menuMapModel MenuMapModel, roleMenuModel RoleMenuModel, userRoleModel UserRoleModel, roleModel RoleModel) (menuAuthList []auth_dto.RoleMenuDto, err error) {
	if unitId == "" || unitUserId == "" {
		str := fmt.Sprintf("GetUserMenu():获取菜单权限必填参数, unit_id:%s, unitUserId:%s", unitId, unitUserId)
		global.Log.Error(str)
		return menuAuthList, errors.New(str)
	}
	isAdmin, err := helper.IsAdmin(moduleName, unitUserId)
	if err != nil {
		return menuAuthList, err
	}

	tableMenu := menuModel.TableName()
	tableMenuMap := menuMapModel.TableName()
	tableRoleMenu := roleMenuModel.TableName()
	tableUserRole := userRoleModel.TableName()
	tableRole := roleModel.TableName()
	tableStruct := struct {
		TableMenu     string
		TableMenuMap  string
		TableRoleMenu string
		TableUserRole string
		TableRole     string
	}{
		TableMenu:     tableMenu,
		TableMenuMap:  tableMenuMap,
		TableRoleMenu: tableRoleMenu,
		TableUserRole: tableUserRole,
		TableRole:     tableRole,
	}

	selectStr, err := helper.ParseStringTpl(`{{.TableMenu}}.*, ANY_VALUE({{.TableMenuMap}}.unit_id) AS unit_id`, tableStruct)
	joinMenuStr, err2 := helper.ParseStringTpl(`inner join {{.TableMenu}} on {{.TableMenu}}.id = {{.TableRoleMenu}}.menu_id`, tableStruct)
	joinMenuMapStr, err3 := helper.ParseStringTpl(`inner join {{.TableMenuMap}} on {{.TableMenuMap}}.menu_id = {{.TableMenu}}.id`, tableStruct)
	joinUserRoleStr, err4 := helper.ParseStringTpl(`inner join {{.TableUserRole}} on {{.TableUserRole}}.role_id = {{.TableRoleMenu}}.role_id`, tableStruct)
	joinRoleStr, err5 := helper.ParseStringTpl(`inner join {{.TableRole}} on {{.TableRole}}.id = {{.TableRoleMenu}}.role_id`, tableStruct)

	if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		err = helper.Ternary(err2 != nil, err2, err)
		err = helper.Ternary(err3 != nil, err3, err)
		err = helper.Ternary(err4 != nil, err4, err)
		err = helper.Ternary(err5 != nil, err5, err)
		return menuAuthList, err
	}

	var tmpError error
	if isAdmin {
		tmpError = global.GetReadDb().
			Model(menuModel).
			Select(selectStr).
			Joins(joinMenuMapStr).
			Where(tableMenuMap+".unit_id = ?", unitId).
			Where(tableMenuMap+".deleted = ?", 0).
			Where(tableMenu+".deleted = ?", 0).
			Where(tableMenu+".menu_type NOT IN ?", operateMenuTypeArr).
			Group(tableMenu + ".id").
			Order(tableMenu + ".rank asc").
			Scan(&menuAuthList).
			Error
	} else {
		tmpError = global.GetReadDb().
			Model(roleMenuModel).
			Select(selectStr).
			Joins(joinMenuStr).
			Joins(joinMenuMapStr).
			Joins(joinRoleStr).
			Joins(joinUserRoleStr).
			Where(tableMenuMap+".unit_id = ?", unitId).
			Where(tableMenuMap+".deleted = ?", 0).
			Where(tableMenu+".deleted = ?", 0).
			Where(tableMenu+".menu_type NOT IN ?", operateMenuTypeArr).
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

/**
 * 页面-获取角色可选权限列表
 */
func GetRoleMenu[MenuModel itf.MenuItf, MenuMapModel itf.MenuMapItf](unitIds []string, menuModel MenuModel, menuMapModel MenuMapModel) (dataList []base_model.UnitMenu, err error) {
	dataList = make([]base_model.UnitMenu, 0)
	if len(unitIds) == 0 {
		str := fmt.Sprintf("GetUserMenu():获取菜单权限必填参数, unit_id:%v", unitIds)
		global.Log.Error(str)
		return dataList, errors.New(str)
	}

	tableMenu := menuModel.TableName()
	tableMenuMap := menuMapModel.TableName()
	tableStruct := struct {
		TableMenu    string
		TableMenuMap string
	}{
		TableMenu:    tableMenu,
		TableMenuMap: tableMenuMap,
	}

	selectStr, err := helper.ParseStringTpl(`{{.TableMenu}}.id,{{.TableMenu}}.parent_id,{{.TableMenu}}.menu_type,{{.TableMenu}}.title`, tableStruct)
	if err != nil {
		return dataList, err
	}

	err = global.GetReadDb().
		Model(menuModel).
		Select(selectStr).
		Where(tableMenuMap+".unit_id in ?", unitIds).
		Where(tableMenuMap+".deleted = 0").
		Where(tableMenu+".deleted = ?", 0).
		Group(tableMenu + ".id").
		Order(tableMenu + ".rank asc").
		Find(&dataList).
		Error

	if err != nil {
		return dataList, err
	}
	return
}

func GetRoleMenuIds[MenuModel itf.MenuItf, MenuMapModel itf.MenuMapItf, RoleMenuModel itf.RoleMenuItf](roleId string, menuModel MenuModel, menuMapModel MenuMapModel, roleMenuModel RoleMenuModel) (dataList []base_model.UnitRoleMenu, err error) {
	if roleId == "" {
		str := fmt.Sprintf("GetRoleMenuIds():获取角色权限必填参数, role_id:%s", roleId)
		global.Log.Error(str)
		return dataList, errors.New(str)
	}
	tableMenuName := menuModel.TableName()
	tableMenuMapName := menuMapModel.TableName()
	tableRoleMenuName := roleMenuModel.TableName()

	err = global.GetReadDb().
		Model(roleMenuModel).
		Select(tableRoleMenuName+".*").
		Joins("inner join "+tableMenuName+" on "+tableMenuName+".id = "+tableRoleMenuName+".menu_id").
		Joins("inner join "+tableMenuMapName+" on "+tableMenuMapName+".menu_id = "+tableRoleMenuName+".menu_id").
		Where(tableRoleMenuName+".role_id = ?", roleId).
		// Where(tableMenuName+".menu_type IN ?", operateMenuTypeArr).
		Where(tableMenuName+".deleted = ?", 0).
		Where(tableMenuMapName+".deleted = ?", 0).
		Find(&dataList).Error
	return
}

func GetUserPermissions[MenuModel itf.MenuItf, MenuMapModel itf.MenuMapItf, RoleMenuModel itf.RoleMenuItf, UserRoleModel itf.UserRoleItf, RoleModel itf.RoleItf](moduleName string, unitId string, unitUserId string, menuModel MenuModel, menuMapModel MenuMapModel, roleMenuModel RoleMenuModel, userRoleModel UserRoleModel, roleModel RoleModel) (menuList []base_model.UnitMenu, err error) {
	if unitUserId == "" {
		str := fmt.Sprintf("GetUserMenu():获取菜单权限必填参数, unit_id:%s, classifyName:%s", unitId, unitUserId)
		global.Log.Error(str)
		return menuList, errors.New(str)
	}
	isAdmin, err := helper.IsAdmin(moduleName, unitUserId)
	if err != nil {
		return menuList, err
	}

	tableMenu := menuModel.TableName()
	tableMenuMap := menuMapModel.TableName()
	tableRoleMenu := roleMenuModel.TableName()
	tableUserRole := userRoleModel.TableName()
	tableRole := roleModel.TableName()
	tableStruct := struct {
		TableMenu     string
		TableMenuMap  string
		TableRoleMenu string
		TableUserRole string
		TableRole     string
	}{
		TableMenu:     tableMenu,
		TableMenuMap:  tableMenuMap,
		TableRoleMenu: tableRoleMenu,
		TableUserRole: tableUserRole,
		TableRole:     tableRole,
	}

	selectStr, err := helper.ParseStringTpl(`{{.TableMenu}}.*, ANY_VALUE({{.TableMenuMap}}.unit_id) AS unit_id`, tableStruct)
	joinMenuStr, err2 := helper.ParseStringTpl(`inner join {{.TableMenu}} on {{.TableMenu}}.id = {{.TableRoleMenu}}.menu_id`, tableStruct)
	joinMenuMapStr, err3 := helper.ParseStringTpl(`inner join {{.TableMenuMap}} on {{.TableMenuMap}}.menu_id = {{.TableMenu}}.id`, tableStruct)
	joinUserRoleStr, err4 := helper.ParseStringTpl(`inner join {{.TableUserRole}} on {{.TableUserRole}}.role_id = {{.TableRoleMenu}}.role_id`, tableStruct)
	joinRoleStr, err5 := helper.ParseStringTpl(`inner join {{.TableRole}} on {{.TableRole}}.id = {{.TableRoleMenu}}.role_id`, tableStruct)

	if err != nil || err2 != nil || err3 != nil || err4 != nil || err5 != nil {
		err = helper.Ternary(err2 != nil, err2, err)
		err = helper.Ternary(err3 != nil, err3, err)
		err = helper.Ternary(err4 != nil, err4, err)
		err = helper.Ternary(err5 != nil, err5, err)
		return menuList, err
	}

	var tmpError error
	if isAdmin {
		tmpError = global.GetReadDb().
			Model(menuModel).
			Select(selectStr).
			Joins(joinMenuMapStr).
			Where(tableMenu+".deleted = 0").
			Where(tableMenu+".menu_type IN ?", operateMenuTypeArr).
			Where(tableMenuMap+".unit_id = ?", unitId).
			Where(tableMenuMap + ".deleted = 0").
			Group(tableMenu + ".id").
			Scan(&menuList).
			Error
	} else {
		tmpError = global.GetReadDb().
			Model(roleMenuModel).
			Select(selectStr).
			Joins(joinMenuStr).
			Joins(joinMenuMapStr).
			Joins(joinRoleStr).
			Joins(joinUserRoleStr).
			Where(tableMenu+".deleted = 0").
			Where(tableMenu+".menu_type IN ?", operateMenuTypeArr).
			Where(tableUserRole+".user_id = ?", unitUserId).
			Where(tableUserRole + ".deleted = 0").
			Where(tableRole + ".status = 1").
			Where(tableRole + ".deleted = 0").
			Where(tableMenuMap + ".deleted = 0").
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
		Where(tableMenuName+".menu_type IN ?", operateMenuTypeArr).
		Find(&menuList).Error
	return
}

// 删除角色已选权限
func DeleteRoleMenuByRole[RoleMenuModel itf.RoleMenuItf](tx *gorm.DB, roleId string, roleMenuModel RoleMenuModel) (err error) {
	if roleId == "" {
		str := fmt.Sprintf("DeleteRoleMenuByRole():删除角色权限必填参数, role_id:%s", roleId)
		global.Log.Error(str)
		return errors.New(str)
	}
	return tx.Model(roleMenuModel).
		Where("role_id = ?", roleId).
		Delete(&roleMenuModel).Error
}

// 保存角色已选权限
func RoleMenuSave[RoleMenuModel itf.RoleMenuItf](tx *gorm.DB, roleId string, menuIds []string, roleMenuModel RoleMenuModel) (err error) {
	if len(menuIds) == 0 {
		return
	}
	if roleId == "" {
		str := fmt.Sprintf("RoleMenuSave():保存角色权限必填参数, role_id:%s", roleId)
		global.Log.Error(str)
		return errors.New(str)
	}

	data := make([]base_model.UnitRoleMenu, 0)
	for _, menuId := range menuIds {
		if menuId == "" {
			continue
		}
		roleMenu := base_model.UnitRoleMenu{}
		roleMenu.Id, _ = helper.GetUuid()
		roleMenu.RoleId = roleId
		roleMenu.MenuId = menuId
		data = append(data, roleMenu)
	}
	if len(data) == 0 {
		err = errors.New("RoleMenuSave():保存角色权限数据异常")
		return
	}

	return tx.Model(roleMenuModel).Create(&data).Error
}
