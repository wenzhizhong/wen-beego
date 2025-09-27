package services

import (
	"WenBeego/apps/common/ar"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"errors"
)

type CommonPlat struct {
	PlatAr ar.PlatAr
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
func changeUnit[UnitModel models.ModelInterface, UnitUserModel models.ModelInterface, UnitUserProfileModel models.ModelInterface](userId string, unitId string) error {
	if userId == "" {
		return errors.New("userId 不能为空")
	}
	_, err := ar.GetUserOfUnitById[UnitUserModel](userId, unitId)
	if err != nil && !helper.DbNotFound(err) {
		return err
	} else if helper.DbNotFound(err) {
		return errors.New("用户不存在")
	} else {
		userProfile, err2 := ar.GetUserProfileOfUnitById[UnitUserModel, UnitUserProfileModel](userId, unitId)
		if err2 != nil {
			return err2
		}

		status := -1
		switch t := any(userProfile).(type) {
		case *models.PlatUserProfile:
			status = t.Status
			if status != models.PLAT_USER_PROFILE_STATUS_NORMAL {
				err = errors.New("用户" + models.PLAT_USER_PROFILE_STATUS_MAP[status])
			}
		case *models.MchntUserProfile:
			status = t.Status
			if status != models.MCHNT_USER_PROFILE_STATUS_NORMAL {
				err = errors.New("用户" + models.MCHNT_USER_PROFILE_STATUS_MAP[status])
			}
		default:
			return errors.New("未知的单位用户信息类型")
		}
		if err != nil {
			return err
		}

		err = ar.UpdateUserDefaultUnit[UnitUserModel](userId, unitId)

	}
	if err != nil {
		return err
	}

	return nil
}
