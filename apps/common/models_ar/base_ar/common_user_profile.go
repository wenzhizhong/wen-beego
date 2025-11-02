package base_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"

	"gorm.io/gorm"
)

/**
 * 获取用户信息
 * UnitUserProfileModel: models.PlatProfileUser, models.MchntProfileUser
 * @param userId
 * @param unitId
 * @return
 * @throws
 */
func GetUserProfileOfUnitById[UnitUserModel itf.UnitUserItf, UnitUserProfileModel itf.UserProfileItf](userId string, unitId string) (base_model.UnitUserProfile, error) {
	var unitUserModel UnitUserModel
	var unitUserProfileModel UnitUserProfileModel
	tableUnitUserName := unitUserModel.TableName()
	tableUnitUserProfileName := unitUserProfileModel.TableName()
	var data base_model.UnitUserProfile
	if userId == "" {
		return data, errors.New("userId 不能为空")
	}
	if unitId == "" {
		return data, errors.New("unitId 不能为空")
	}

	result := global.GetReadDb().
		Model(unitUserProfileModel).
		Select(tableUnitUserProfileName+".*").
		Joins("inner join "+tableUnitUserName+" on "+tableUnitUserName+".id = "+tableUnitUserProfileName+".id").
		Where(tableUnitUserName+".user_id = ?", userId).
		Where(tableUnitUserName+".unit_id = ?", unitId).
		Where(tableUnitUserName+".deleted = ?", 0).
		Take(&data)
	return data, result.Error
}

/**
 * 新增用户信息
 */
func InsertUnitUserProfile[UnitUserProfileModel itf.UserProfileItf](tx *gorm.DB, unitUserProfileModel base_model.UnitUserProfile) (err error) {
	if unitUserProfileModel.Id == "" {
		return errors.New("新增用户信息，id 不能为空")
	}
	var tmpUnitUserProfileModel UnitUserProfileModel
	err = global.GetReadDb().
		Model(tmpUnitUserProfileModel).
		Where("id = ?", unitUserProfileModel.Id).
		Take(&tmpUnitUserProfileModel).Error
	if err == nil && tmpUnitUserProfileModel.GetId() != "" {
		return nil
	}

	err = tx.Model(tmpUnitUserProfileModel).
		Create(&unitUserProfileModel).Error
	if err != nil {
		return err
	}
	return nil
}
