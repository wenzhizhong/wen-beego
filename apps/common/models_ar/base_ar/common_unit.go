package base_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/itf"
)

/**
 * 获取用户组织列表
 * @param {string} userId 用户ID
 * @param {models.ModelInterface} unitModel models.Plat, models.Mchnt
 * @param {models.ModelInterface} unitUserModel models.PlatUser, models.MchntUser
 * @return {[]models.ModelInterface} []models.Plat, []models.Mchnt
 * @return {error}
 */
func GetUserUnitList[UnitModel itf.UnitItf, UnitUserModel itf.UnitUserItf](userId string, unitModel UnitModel, unitUserModel UnitUserModel) ([]UnitModel, error) {
	tableUnitName := unitModel.TableName()
	tableUnitUserName := unitUserModel.TableName()

	tableStruct := struct {
		TableUnit     string
		TableUserUnit string
	}{
		TableUnit:     tableUnitName,
		TableUserUnit: tableUnitUserName,
	}

	selectStr, err := helper.ParseStringTpl(`{{.TableUnit}}.*`, tableStruct)
	joinStr, err2 := helper.ParseStringTpl(`inner join {{.TableUserUnit}} on {{.TableUserUnit}}.unit_id = {{.TableUnit}}.id`, tableStruct)
	if err != nil {
		return nil, err
	}
	if err2 != nil {
		return nil, err2
	}

	listData := []UnitModel{}
	result := global.GetReadDb().
		Model(unitModel).
		Select(selectStr).
		Joins(joinStr).
		Where(tableUnitUserName+".user_id = ?", userId).
		Where(tableUnitUserName + ".deleted = 0").
		Scan(&listData)
	return listData, result.Error
}

// 查找是否存在用户组织单位
func GetUserUnitById[UnitModel itf.UnitItf, UnitUserModel itf.UnitUserItf](userId string, unitId string, unitModel UnitModel, unitUserModel UnitUserModel) (UnitModel, error) {
	tableUnitName := unitModel.TableName()
	tableUnitUserName := unitUserModel.TableName()

	tableStruct := struct {
		TableUnit     string
		TableUserUnit string
	}{
		TableUnit:     tableUnitName,
		TableUserUnit: tableUnitUserName,
	}

	selectStr, err := helper.ParseStringTpl(`{{.TableUnit}}.*`, tableStruct)
	joinStr, err2 := helper.ParseStringTpl(`inner join {{.TableUserUnit}} on {{.TableUserUnit}}.unit_id = {{.TableUnit}}.id`, tableStruct)
	if err != nil {
		return unitModel, err
	}
	if err2 != nil {
		return unitModel, err2
	}

	result := global.GetReadDb().
		Model(unitModel).
		Select(selectStr).
		Joins(joinStr).
		Where(tableUnitUserName+".user_id = ?", userId).
		Where(tableUnitUserName+".unit_id = ?", unitId).
		Where(tableUnitUserName + ".deleted = 0").
		Take(&unitModel)
	return unitModel, result.Error
}
