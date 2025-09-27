package ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"errors"

	"gorm.io/gorm"
)

/**
 * 获取用户
 * UnitUserModel: models.PlatUser, models.MchntUser
 * @param userId
 * @param unitId
 * @return
 * @throws
 */
func GetUserOfUnitById[UnitUserModel models.ModelInterface](userId string, unitId string) (map[string]interface{}, error) {
	var unitUserModel UnitUserModel
	tableUnitUserName := unitUserModel.TableName()
	var userData map[string]interface{}
	if userId == "" {
		return userData, errors.New("userId 不能为空")
	}

	result := global.GetReadDb().
		Model(unitUserModel).
		Select("*, case is_default when 1 then unit_id else '' end AS default_unit_id").
		Where(tableUnitUserName+".user_id = ?", userId).
		Where(tableUnitUserName+".unit_id = ?", unitId).
		Where(tableUnitUserName+".deleted = ?", 0).
		Take(&userData)
	return userData, result.Error
}

/**
 * 获取用户默认组织单位
 * UnitModel: models.PlatUnit, models.MchntUnit
 * UnitUserModel: models.PlatUser, models.MchntUser
 * @param userId
 * @return
 * @throws
 */
func GetUserDefaultUnit[UnitModel models.ModelInterface, UnitUserModel models.ModelInterface](userId string) (map[string]interface{}, error) {
	var unitModel UnitModel
	var unitUserModel UnitUserModel
	tableUnitName := unitModel.TableName()
	tableUnitUserName := unitUserModel.TableName()
	var userData map[string]interface{}
	if userId == "" {
		return userData, errors.New("userId 不能为空")
	}

	tableStruct := struct {
		TableUnit     string
		TableUserUnit string
	}{
		TableUnit:     tableUnitName,
		TableUserUnit: tableUnitUserName,
	}
	joinUserUnitStr, err1 := helper.ParseStringTpl(`inner join {{.TableUserUnit}} on {{.TableUserUnit}}.unit_id = {{.TableUnit}}.id`, tableStruct)
	selectStr, err2 := helper.ParseStringTpl(`{{.TableUnit}}.*, case {{.TableUserUnit}}.is_default when 1 then {{.TableUserUnit}}.unit_id else '' end AS default_unit_id`, tableStruct)
	if err1 != nil {
		return userData, err1
	}
	if err2 != nil {
		return userData, err2
	}

	result := global.GetReadDb().
		Model(unitModel).
		Select(selectStr).
		Joins(joinUserUnitStr).
		Where(tableUnitUserName+".user_id = ?", userId).
		Where(tableUnitUserName+".is_default = ?", 1).
		Where(tableUnitUserName+".deleted = ?", 0).
		Take(&userData)
	return userData, result.Error
}

/**
 * 更新用户默认组织单位id
 * UnitUserModel: models.PlatUser, models.MchntUser
 * @param userId
 * @param unitId
 * @return
 * @throws
 */
func UpdateUserDefaultUnit[UnitUserModel models.ModelInterface](userId string, unitId string) error {
	var unitUserModel UnitUserModel
	_, err := GetUserOfUnitById[UnitUserModel](userId, unitId)
	if err != nil {
		return err
	}

	updateData := struct {
		IsDefault int `json:"is_default"`
	}{
		IsDefault: 0,
	}
	err = global.GetWriteDb().Transaction(func(tx *gorm.DB) error {
		result := global.GetWriteDb().Model(unitUserModel).Where("user_id = ?", userId).Updates(updateData)
		if result.Error != nil {
			return result.Error
		}

		updateData.IsDefault = 1
		result = global.GetWriteDb().Model(unitUserModel).Where("user_id = ?", userId).Where("unit_id = ?", unitId).Updates(updateData)
		if result.Error != nil {
			return result.Error
		}
		return nil
	})
	return err
}

/**
 * 新增组织单位的用户
 * UnitUserModel: models.PlatUser, models.MchntUser
 * @param userId
 * @param unitId
 * @param isDefault
 * @param isAdmin
 * @return
 * @throws
 */
func AddUserOfUnit[UnitUserModel models.ModelInterface](userId string, unitId string, isDefault int, isAdmin int) error {
	uuid, _ := helper.GetUuid()
	insertData := struct {
		Id        string `json:"id"`
		UserId    string `json:"user_id"`
		UnitId    string `json:"unit_id"`
		IsDefault int    `json:"is_default"`
		IsAdmin   int    `json:"is_admin"`
		Deleted   int    `json:"deleted"`
	}{
		Id:        uuid,
		UserId:    userId,
		UnitId:    unitId,
		IsDefault: isDefault,
		IsAdmin:   isAdmin,
		Deleted:   0,
	}

	var unitUserModel UnitUserModel
	result := global.GetWriteDb().Model(unitUserModel).Create(&insertData)
	return result.Error
}
