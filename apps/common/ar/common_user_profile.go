package ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
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
func GetUserProfileOfUnitById[UnitUserModel models.ModelInterface, UnitUserProfileModel models.ModelInterface](userId string, unitId string) (UnitUserProfileModel, error) {
	var unitUserModel UnitUserModel
	var unitUserProfileModel UnitUserProfileModel
	tableUnitUserName := unitUserModel.TableName()
	tableUnitUserProfileName := unitUserProfileModel.TableName()
	if userId == "" {
		return unitUserProfileModel, errors.New("userId 不能为空")
	}
	if unitId == "" {
		return unitUserProfileModel, errors.New("unitId 不能为空")
	}

	result := global.GetReadDb().
		Model(unitUserProfileModel).
		Select(tableUnitUserProfileName+".*").
		Joins("inner join "+tableUnitUserName+" on "+tableUnitUserName+".id = "+tableUnitUserProfileName+".id").
		Where(tableUnitUserName+".user_id = ?", userId).
		Where(tableUnitUserName+".unit_id = ?", unitId).
		Where(tableUnitUserName+".deleted = ?", 0).
		Take(&unitUserProfileModel)
	return unitUserProfileModel, result.Error
}
