package base_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"
	"strings"

	"gorm.io/gorm"
)

/**
 * 克隆菜单
 */
func CloneMenu[UnitMenuModel itf.MenuItf, UnitMenuPermsModel itf.MenuPermsItf](tx *gorm.DB, oldUnitId string, newUnitId string) (err error) {
	var unitMenuModel UnitMenuModel
	var unitMenuPermsModel UnitMenuPermsModel
	var menuTableName = unitMenuModel.TableName()
	var menuPermsTableName = unitMenuPermsModel.TableName()
	var insertMenuData = []base_model.UnitMenu{}
	var insertMenuPermsData = []base_model.UnitMenuPerms{}
	menuOldIdMapNew := map[string]string{}
	if oldUnitId == "" {
		return errors.New("克隆菜单，原组织单位id 不能为空")
	}
	if newUnitId == "" {
		return errors.New("克隆菜单，新组织单位id 不能为空")
	}

	err = global.GetReadDb().Model(unitMenuModel).Where("clone = 1 and unit_id = ?", oldUnitId).Find(&insertMenuData).Error
	if err != nil {
		return err
	}
	err = global.GetReadDb().Model(unitMenuPermsModel).
		Select(menuPermsTableName+".*").
		Joins("inner join "+menuTableName+" on "+menuPermsTableName+".menu_id = "+menuTableName+".id").
		Where(menuTableName+".unit_id = ?", oldUnitId).
		Find(&insertMenuPermsData).Error
	if err != nil {
		return err
	}

	for _, item := range insertMenuData {
		newId, _ := helper.GetUuid()
		menuOldIdMapNew[item.Id] = newId
	}
	for key, item := range insertMenuData {
		newMenuId, _ := helper.GetUuid()
		insertMenuData[key].Id = newMenuId
		insertMenuData[key].UnitId = newUnitId
		insertMenuData[key].CreatedAt = helper.GetTime()
		if item.Pid != "" {
			insertMenuData[key].Pid = menuOldIdMapNew[item.Pid]

			newAllPids := make([]string, 0)
			allPids := strings.Split(item.AllPid, ",")
			for _, pid := range allPids {
				if pid == "" {
					continue
				}
				newAllPids = append(newAllPids, menuOldIdMapNew[pid])
			}
			insertMenuData[key].AllPid = strings.Join(newAllPids, ",")
		}
	}
	for key, item := range insertMenuPermsData {
		newId, _ := helper.GetUuid()
		insertMenuPermsData[key].Id = newId
		insertMenuPermsData[key].MenuId = menuOldIdMapNew[item.MenuId]
		insertMenuPermsData[key].CreatedAt = helper.GetTime()
		insertMenuPermsData[key].UpdatedAt = helper.GetTime()
	}
	// menu 组织单位菜单
	err = tx.Model(unitMenuModel).
		Create(&insertMenuData).Error
	if err != nil {
		return err
	}
	// menu_perms 组织单位菜单权限
	err = tx.Model(unitMenuPermsModel).
		Create(&insertMenuPermsData).Error
	if err != nil {
		return err
	}
	return
}
