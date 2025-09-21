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
		err = changeUnit[*models.Plat, *models.PlatUser](userId, unitId)
	} else {
		err = changeUnit[*models.Mchnt, *models.MchntUser](userId, unitId)
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
func changeUnit[UnitModel models.ModelInterface, UserUnitModel models.ModelInterface](userId string, unitId string) error {
	if userId == "" {
		return errors.New("userId 不能为空")
	}
	userData, err := ar.GetUserOfUnitById[UserUnitModel](userId, unitId)
	if err != nil && !helper.DbNotFound(err) {
		return err
	} else if helper.DbNotFound(err) {
		err = ar.AddUserOfUnit[UserUnitModel](userId, unitId, 1, 0)
	} else {
		status := userData["status"].(int)
		if status == 0 {
			return errors.New("用户已禁用")
		}
		if status == 2 {
			return errors.New("用户已离职")
		}
		err = ar.UpdateUserDefaultUnit[UserUnitModel](userId, unitId)

	}
	if err != nil {
		return err
	}

	return nil
}
