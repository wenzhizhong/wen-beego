package models_ar

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/auth_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/global/constant"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models_ar/base_ar"
	"errors"
	"fmt"
)

type PlatMenuViewAr struct {
	models.PlatMenuView
}

func (a *PlatMenuViewAr) GetById(id string) (models.PlatMenu, error) {
	data := models.PlatMenu{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}

func (a *PlatMenuViewAr) GetUserMenu(moduleName string, unitId string, unitUserId string, menuViewModel *models.PlatMenuView, menuMapModel *models.PlatMenuMapView, roleMenuModel *models.PlatRoleMenu, userRoleModel *models.PlatUserRole, roleModel *models.PlatRole) (menuAuthList []auth_dto.RoleMenuDto, err error) {
	if unitId == "" || unitUserId == "" {
		str := fmt.Sprintf("GetUserMenu():获取菜单权限必填参数, unit_id:%s, unitUserId:%s", unitId, unitUserId)
		global.Log.Error(str)
		return menuAuthList, errors.New(str)
	}
	isAdmin, err := helper.IsAdmin(moduleName, unitUserId)
	if err != nil {
		return menuAuthList, err
	}

	tableMenu := menuViewModel.TableName()
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

	selectStr, err := helper.ParseStringTpl(`DISTINCT ON ({{.TableMenu}}.id) {{.TableMenu}}.*, {{.TableMenuMap}}.unit_id`, tableStruct)
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
		subQuery := global.GetReadDb().
			Model(menuViewModel).
			Select(selectStr).
			Joins(joinMenuMapStr).
			Where(tableMenuMap+".unit_id = ?", unitId).
			Where(tableMenuMap+".deleted = ?", 0).
			Where(tableMenu+".deleted = ?", 0).
			Where(tableMenu+".menu_type NOT IN ?", base_ar.OperateMenuTypeArr)
		tmpError = global.GetReadDb().
			Table("(?) tmp", subQuery).
			Order("rank asc").
			Scan(&menuAuthList).
			Error
	} else {
		subQuery := global.GetReadDb().
			Model(roleMenuModel).
			Select(selectStr).
			Joins(joinMenuStr).
			Joins(joinMenuMapStr).
			Joins(joinRoleStr).
			Joins(joinUserRoleStr).
			Where(tableMenuMap+".unit_id = ?", unitId).
			Where(tableMenuMap+".deleted = ?", 0).
			Where(tableMenu+".deleted = ?", 0).
			Where(tableMenu+".menu_type NOT IN ?", base_ar.OperateMenuTypeArr).
			Where(tableRole+".deleted = ?", 0).
			Where(tableRole+".status = ?", 1).
			Where(tableUserRole+".user_id = ?", unitUserId).
			Where(tableUserRole+".deleted = ?", 0)
		tmpError = global.GetReadDb().
			Table("(?) tmp", subQuery).
			Order("rank asc").
			Scan(&menuAuthList).
			Error
	}
	if tmpError != nil && !helper.DbNotFound(tmpError) {
		return menuAuthList, tmpError
	}
	return
}

func (a *PlatMenuViewAr) GetRoleMenuIds(baseParamDto dto.BaseParamDto, roleId string, menuModel models.PlatMenuView, menuMapModel models.PlatMenuMapView, roleMenuModel models.PlatRoleMenu) (dataList []base_model.UnitRoleMenu, err error) {
	if roleId == "" {
		str := fmt.Sprintf("GetRoleMenuIds():获取角色权限必填参数, role_id:%s", roleId)
		global.Log.Error(str)
		return dataList, errors.New(str)
	}

	tableMenuName := menuModel.TableName()
	tableMenuMapName := menuMapModel.TableName()
	tableRoleMenuName := roleMenuModel.TableName()

	// if baseParamDto.IsOfficial {
	// 	err = global.GetReadDb().
	// 		Model(roleMenuModel).
	// 		Select(tableRoleMenuName+".*").
	// 		Where(tableRoleMenuName+".role_id = ?", roleId).
	// 		Find(&dataList).Error
	// } else {

	err = global.GetReadDb().
		Model(roleMenuModel).
		Select(tableRoleMenuName+".*").
		Joins("inner join "+tableMenuName+" on "+tableMenuName+".id = "+tableRoleMenuName+".menu_id").
		Joins("inner join "+tableMenuMapName+" on "+tableMenuMapName+".menu_id = "+tableRoleMenuName+".menu_id").
		Where(tableRoleMenuName+".role_id = ?", roleId).
		// Where(tableMenuName+".menu_type IN ?", OperateMenuTypeArr).
		Where(tableMenuName+".deleted = ?", 0).
		Where(tableMenuMapName+".deleted = ?", 0).
		Find(&dataList).Error

	// }
	return
}
func (a *PlatMenuViewAr) GetUserPermissions(moduleName string, unitId string, unitUserId string, menuModel models.PlatMenuView, menuMapModel models.PlatMenuMapView, roleMenuModel models.PlatRoleMenu, userRoleModel models.PlatUserRole, roleModel models.PlatRole) (menuList []base_model.UnitMenu, err error) {
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

	selectStr, err := helper.ParseStringTpl(`DISTINCT ON ({{.TableMenu}}.id) {{.TableMenu}}.*, {{.TableMenuMap}}.unit_id`, tableStruct)
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
			Where(tableMenu+".menu_type IN ?", base_ar.OperateMenuTypeArr).
			Where(tableMenuMap+".unit_id = ?", unitId).
			Where(tableMenuMap + ".deleted = 0").
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
			Where(tableMenu+".menu_type IN ?", base_ar.OperateMenuTypeArr).
			Where(tableUserRole+".user_id = ?", unitUserId).
			Where(tableUserRole + ".deleted = 0").
			Where(tableRole + ".status = 1").
			Where(tableRole + ".deleted = 0").
			Where(tableMenuMap + ".deleted = 0").
			Scan(&menuList).
			Error
	}
	if tmpError != nil && !helper.DbNotFound(tmpError) {
		return menuList, tmpError
	}
	return
}

func (a *PlatMenuViewAr) GetRoleMenu(unitIds []string, menuModel models.PlatMenuView, menuMapModel models.PlatMenuMap) (dataList []base_model.UnitMenu, err error) {
	dataList = make([]base_model.UnitMenu, 0)
	if len(unitIds) == 0 {
		str := fmt.Sprintf("GetUserMenu():获取菜单权限必填参数, unit_id:%v", unitIds)
		global.Log.Error(str)
		return dataList, errors.New(str)
	}

	isOfficial, err1 := helper.IsOfficial(constant.ADMIN_PLAT, unitIds[0])
	if err1 != nil {
		err = err1
		return
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

	selectStr, err := helper.ParseStringTpl(`DISTINCT ON ({{.TableMenu}}.id) {{.TableMenu}}.id,{{.TableMenu}}.parent_id,{{.TableMenu}}.menu_type,{{.TableMenu}}.title,{{.TableMenu}}.rank`, tableStruct)
	selectStr2, err2 := helper.ParseStringTpl(`{{.TableMenu}}.id,{{.TableMenu}}.parent_id,{{.TableMenu}}.menu_type,{{.TableMenu}}.title,{{.TableMenu}}.rank`, tableStruct)
	joinStr, err3 := helper.ParseStringTpl(`INNER JOIN {{.TableMenuMap}} ON {{.TableMenuMap}}.menu_id={{.TableMenu}}.id`, tableStruct)
	if err != nil || err2 != nil || err3 != nil {
		err = helper.Ternary(err2 != nil, err2, err)
		err = helper.Ternary(err3 != nil, err3, err)
		return dataList, err
	}

	if isOfficial {
		err = global.GetReadDb().
			Model(menuModel).
			Select(selectStr2).
			Where("deleted = 0").
			Order("rank asc").
			Scan(&dataList).
			Error
	} else {
		subQuery := global.GetReadDb().
			Model(menuModel).
			Select(selectStr).
			Joins(joinStr).
			Where(tableMenuMap+".unit_id in ?", unitIds).
			Where(tableMenuMap+".deleted = 0").
			Where(tableMenu+".deleted = ?", 0)
		err = global.GetReadDb().
			Table("(?) tmp", subQuery).
			Order("rank asc").
			Find(&dataList).
			Error
	}

	if err != nil {
		return dataList, err
	}
	return
}
