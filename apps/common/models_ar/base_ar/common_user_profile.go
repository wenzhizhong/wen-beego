package base_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"errors"
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
