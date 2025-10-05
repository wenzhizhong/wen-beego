package services

import (
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models/itf"
	"WenBeego/apps/common/models_ar"
	"WenBeego/apps/common/models_ar/base_ar"
	"errors"
)

type CommonPlat struct {
	PlatAr models_ar.PlatAr
}

func (a *CommonPlat) ChangeUnit(moduleName string, userId string, unitId string) (result interface{}, err error) {
	if moduleName == "admin_plat" {
		err = changeUnit[*models.Plat, *models.PlatUser, *models.PlatUserProfile](userId, unitId)
	} else {
		err = changeUnit[*models.Mchnt, *models.MchntUser, *models.MchntUserProfile](userId, unitId)
	}
	if err != nil {
		return nil, err
	}
	result, err = (&CommonAuth{}).GetAdminLoginInfo(moduleName, userId)
	if err != nil {
		return nil, err
	}
	return result, nil
}
func changeUnit[UnitModel itf.UnitItf, UnitUserModel itf.UnitUserItf, UnitUserProfileModel itf.UserProfileItf](userId string, unitId string) error {
	if userId == "" {
		return errors.New("userId 不能为空")
	}
	_, err := base_ar.GetUserOfUnitById[UnitUserModel](userId, unitId)
	if err != nil && !helper.DbNotFound(err) {
		return err
	} else if helper.DbNotFound(err) {
		return errors.New("用户不存在")
	} else {
		userProfile, err2 := base_ar.GetUserProfileOfUnitById[UnitUserModel, UnitUserProfileModel](userId, unitId)
		if err2 != nil {
			return err2
		}

		if userProfile.Status != base_model.UNIT_USER_PROFILE_NORMAL {
			err = errors.New("用户" + base_model.UNIT_USER_PROFILE_MAP[userProfile.Status])
			return err
		}
		err = base_ar.UpdateUserDefaultUnit[UnitUserModel](userId, unitId)
	}
	if err != nil {
		return err
	}

	return nil
}
