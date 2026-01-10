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
func CloneMenuMap[UnitMenuModel itf.MenuItf, UnitMenuMapModel itf.MenuMapItf](tx *gorm.DB, oldUnitId string, newUnitId string) (err error) {
	var unitMenuModel UnitMenuModel
	var unitMenuMapModel UnitMenuMapModel
	var data = []base_model.UnitMenu{}
	var insertMenuMapData = []base_model.UnitMenuMap{}

	if oldUnitId == "" {
		return errors.New("克隆菜单，原组织单位id 不能为空")
	}
	if newUnitId == "" {
		return errors.New("克隆菜单，新组织单位id 不能为空")
	}

	err = global.GetReadDb().Model(unitMenuModel).Where("clone = 1 and unit_id = ?", oldUnitId).Find(&data).Error
	if err != nil || len(data) <= 0 {
		return err
	}

	for _, item := range data {
		newMenuId, err := helper.GetUuid()
		if err != nil {
			return err
		}
		newInsertItem := base_model.UnitMenuMap{}
		newInsertItem.Id = newMenuId
		newInsertItem.UnitId = newUnitId
		newInsertItem.MenuId = item.Id
		newInsertItem.UpdatedAt = helper.GetTimestamp()
		newInsertItem.Deleted = 0
		insertMenuMapData = append(insertMenuMapData, newInsertItem)

	}
	// menu 组织单位菜单
	err = tx.Model(unitMenuMapModel).
		Create(&insertMenuMapData).Error
	if err != nil {
		return err
	}
	return
}
