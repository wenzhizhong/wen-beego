package base_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"

	"gorm.io/gorm"
)

/**
 * 克隆菜单
 */
func CloneMenu[UnitMenuModel itf.MenuItf](tx *gorm.DB, oldUnitId string, newUnitId string) (err error) {
	var unitMenuModel UnitMenuModel
	var insertMenuData = []base_model.UnitMenu{}
	menuOldIdMapNew := map[string]string{}
	if oldUnitId == "" {
		return errors.New("克隆菜单，原组织单位id 不能为空")
	}
	if newUnitId == "" {
		return errors.New("克隆菜单，新组织单位id 不能为空")
	}

	err = global.GetReadDb().Model(unitMenuModel).Where("deleted = 0 and clone = 1 and unit_id = ?", oldUnitId).Find(&insertMenuData).Error
	if err != nil {
		return err
	}

	for _, item := range insertMenuData {
		newId, err := helper.GetUuid()
		menuOldIdMapNew[item.Id] = newId

		if err != nil {
			return err
		}
	}
	for key, item := range insertMenuData {
		newMenuId, err := helper.GetUuid()
		if err != nil {
			return err
		}

		insertMenuData[key].Id = newMenuId
		insertMenuData[key].UnitId = newUnitId
		insertMenuData[key].CreatedAt = helper.GetTime()
		if item.ParentId != "" {
			insertMenuData[key].ParentId = menuOldIdMapNew[item.ParentId]
		}
	}
	// menu 组织单位菜单
	err = tx.Model(unitMenuModel).
		Create(&insertMenuData).Error
	if err != nil {
		return err
	}
	return
}

func GetMenuListByTitle[UnitMenuModel itf.MenuItf](unitIds []string, title string) (data []base_model.UnitMenu, err error) {
	var unitMenuModel UnitMenuModel
	data = make([]base_model.UnitMenu, 0)

	if len(unitIds) == 0 {
		return data, errors.New("GetMenuListByTitle():unitIds不能为空")
	}
	if title == "" {
		return data, errors.New("GetMenuListByTitle(): title不能为空")
	}

	err = global.GetReadDb().
		Model(unitMenuModel).
		Select("id").
		Where("deleted = 0 and title = ? and unit_id in ?", title, unitIds).
		Find(&data).Error
	if err != nil {
		return data, err
	}

	return
}
