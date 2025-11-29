package base_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
 * 新增默认组织单位-新增其默认用户信息
 */
func InsertUnitUserProfileForCreateUnit[UnitUserProfileModel itf.UserProfileItf](tx *gorm.DB, baseUnitUserProfileModel base_model.UnitUserProfile, unitUserProfileModel UnitUserProfileModel) (err error) {
	if baseUnitUserProfileModel.Id == "" {
		return errors.New("新增用户信息，id 不能为空")
	}
	err = global.GetReadDb().
		Model(unitUserProfileModel).
		Where("id = ?", baseUnitUserProfileModel.Id).
		Take(&baseUnitUserProfileModel).Error
	if err == nil && baseUnitUserProfileModel.Id != "" {
		return nil
	}

	err = tx.Model(unitUserProfileModel).
		Create(&baseUnitUserProfileModel).Error
	if err != nil {
		return err
	}
	return nil
}

// 新增用户信息
func UpsertUnitUserProfile[UnitUserProfileModel itf.UserProfileItf](tx *gorm.DB, baseUnitUserProfileModel base_model.UnitUserProfile, unitUserProfileModel UnitUserProfileModel) (err error) {
	if baseUnitUserProfileModel.Id == "" {
		return errors.New("新增用户信息，id 不能为空")
	}

	err = tx.Model(unitUserProfileModel).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(&baseUnitUserProfileModel).Error
	if err != nil {
		return err
	}
	return nil
}
