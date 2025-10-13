package base_ar

import (
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
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
func GetUserOfUnitById[UnitUserModel itf.UnitUserItf](userId string, unitId string) (base_model.UnitUser, error) {
	var unitUserModel UnitUserModel
	tableUnitUserName := unitUserModel.TableName()
	var userData base_model.UnitUser
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
func GetUserListOfUnitById[UnitUserModel itf.UnitUserItf](reqDto page_dto.SystemUserListReqDto) (userData []base_model.UnitUser, count int64, err error) {
	var unitUserModel UnitUserModel
	tableUnitUserName := unitUserModel.TableName()
	userData = make([]base_model.UnitUser, 0)
	if reqDto.UnitId == "" {
		return userData, count, errors.New("GetUserListOfUnitById(): UnitId 不能为空")
	}
	selectStr := "*, case is_default when 1 then unit_id else '' end AS default_unit_id"
	query := global.GetReadDb().
		Model(unitUserModel).
		Where(tableUnitUserName+".unit_id = ?", reqDto.UnitId).
		Where(tableUnitUserName+".deleted = ?", 0)

	err = query.Select("id").Count(&count).Error
	if err != nil {
		return userData, count, err
	}
	result := query.
		Select(selectStr).
		Limit(reqDto.PageSize).
		Offset(reqDto.Offset).
		Find(&userData)

	return userData, count, result.Error
}

/**
 * 获取用户默认组织单位
 * UnitModel: models.PlatUnit, models.MchntUnit
 * UnitUserModel: models.PlatUser, models.MchntUser
 * @param userId
 * @return
 * @throws
 */
func GetUserDefaultUnit[UnitModel itf.UnitItf, UnitUserModel itf.UnitUserItf](userId string) (base_model.Unit, error) {
	var unitModel UnitModel
	var unitUserModel UnitUserModel
	tableUnitName := unitModel.TableName()
	tableUnitUserName := unitUserModel.TableName()
	var userData base_model.Unit
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
func UpdateUserDefaultUnit[UnitUserModel itf.UnitUserItf](userId string, unitId string) error {
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
func AddUserOfUnit[UnitUserModel itf.UnitUserItf](userId string, unitId string, isDefault int, isAdmin int) error {
	uuid, _ := helper.GetUuid()
	var insertData = base_model.UnitUser{
		Id:        uuid,
		UserId:    userId,
		UnitId:    unitId,
		IsDefault: isDefault,
		Deleted:   0,
	}

	var unitUserModel UnitUserModel
	result := global.GetWriteDb().Model(unitUserModel).Create(&insertData)
	return result.Error
}
