package ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
)

/**
 * 获取用户组织列表
 * @param {string} userId 用户ID
 * @param {models.ModelInterface} unitModel models.Plat, models.Mchnt
 * @param {models.ModelInterface} unitUserModel models.PlatUser, models.MchntUser
 * @return {[]models.ModelInterface} []models.Plat, []models.Mchnt
 * @return {error}
 */
func GetUserUnitList[UnitModel models.ModelInterface, UnitUserModel models.ModelInterface](userId string, unitModel UnitModel, unitUserModel UnitUserModel) ([]UnitModel, error) {
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
