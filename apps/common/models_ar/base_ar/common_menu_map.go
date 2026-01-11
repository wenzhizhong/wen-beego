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
 * 创建新增组织单位的菜单
 */
func CreateNewUnitMenuMap[UnitMenuModel itf.MenuItf, UnitMenuMapModel itf.MenuMapItf](tx *gorm.DB, newUnitId string) (err error) {
	var unitMenuModel UnitMenuModel
	var unitMenuMapModel UnitMenuMapModel
	var data = []base_model.UnitMenu{}
	var insertMenuMapData = []base_model.UnitMenuMap{}

	if newUnitId == "" {
		return errors.New("克隆菜单，新组织单位id 不能为空")
	}

	err = global.GetReadDb().Model(unitMenuModel).Select("id, title").Where("clone = 1").Find(&data).Error
	if err != nil || len(data) <= 0 {
		return err
	}

	for _, item := range data {
		newInsertItem, err := generateMenuMapItem(newUnitId, item.Id)
		if err != nil {
			return err
		}
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

func GetAllMenuMap[UnitModel itf.UnitItf, UnitMenuModel itf.MenuItf, UnitMenuMapModel itf.MenuMapItf](unitIds []string, unitModel UnitModel, unitMenuModel UnitMenuModel, unitMenuMapModel UnitMenuMapModel) (data []base_model.UnitMenuMap, err error) {
	data = make([]base_model.UnitMenuMap, 0)

	query := global.GetReadDb().
		Model(unitMenuMapModel).
		Where(unitMenuMapModel.TableName() + ".deleted = 0")
	if len(unitIds) > 0 {
		query = query.Where(unitMenuMapModel.TableName()+".unit_id IN ?", unitIds)
	}
	result := query.Find(&data)
	return data, result.Error
}

// 刷新组织单位菜单
func RefreshUnitMenuMap[UnitModel itf.UnitItf, UnitMenuModel itf.MenuItf, UnitMenuMapModel itf.MenuMapItf](tx *gorm.DB, unitIds []string, menuId string, unitModel UnitModel, unitMenuModel UnitMenuModel, unitMenuMapModel UnitMenuMapModel) (err error) {

	if len(unitIds) <= 0 {
		return
	}
	// 是否存在数据
	menuMapList, err := GetAllMenuMap(unitIds, unitModel, unitMenuModel, unitMenuMapModel)
	if err != nil {
		return err
	}
	menuMapExists := make(map[string]bool)
	for _, item := range menuMapList {
		menuMapExists[item.UnitId+item.MenuId] = true
	}

	// 新增
	var insertMenuMapData = []base_model.UnitMenuMap{}
	for _, item := range unitIds {
		if menuMapExists[item+menuId] {
			continue
		}
		newInsertItem, err := generateMenuMapItem(item, menuId)
		if err != nil {
			return err
		}
		insertMenuMapData = append(insertMenuMapData, newInsertItem)
	}
	if len(insertMenuMapData) <= 0 {
		return
	}

	err = tx.Model(unitMenuMapModel).
		Create(&insertMenuMapData).Error
	return
}

func generateMenuMapItem(unitId string, menuId string) (base_model.UnitMenuMap, error) {
	newInsertItem := base_model.UnitMenuMap{}
	newMenuMapId, err := helper.GetUuid()
	if err != nil {
		return newInsertItem, err
	}

	newInsertItem.Id = newMenuMapId
	newInsertItem.UnitId = unitId
	newInsertItem.MenuId = menuId
	newInsertItem.UpdatedAt = helper.GetTimestamp()
	newInsertItem.Deleted = 0
	return newInsertItem, nil
}
