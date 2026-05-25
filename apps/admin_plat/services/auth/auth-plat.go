package auth

import (
	"WenBeego/apps/common/dto_vo/auth_dto"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models_ar/base_ar"
	"WenBeego/apps/common/services/auth"
	"errors"
)

type Plat struct {
	CommonUnit auth.CommonUnit
}

// 获取组织信息
func (s *Plat) GetUserUnitList(userId, host string) (data interface{}, err error) {
	if userId == "" {
		return nil, errors.New("获取组织信息失败，请先登录！")
	}
	dataList, err := base_ar.GetUserUnitList(userId, &models.Plat{}, &models.PlatUser{})
	for k, v := range dataList {
		tmpLogo, err1 := helper.LocalFileSign(host, v.Logo)
		if err1 != nil {
			continue
		}
		dataList[k].LogoLink = tmpLogo
	}

	data = struct {
		List interface{} `json:"list"`
	}{
		List: dataList,
	}
	return
}
func (s *Plat) ChangeUnit(moduleName string, userId string, changeUnitDto auth_dto.ChangeUnitDto) (result interface{}, err error) {
	unitId := changeUnitDto.Id
	if userId == "" || unitId == "" {
		return nil, errors.New("切换组织失败，请先登录！")
	}

	unitInfo, err := base_ar.GetUserUnitById(userId, unitId, &models.Plat{}, &models.PlatUser{})
	if err != nil {
		return nil, err
	}
	if unitInfo.Id == "" {
		return nil, errors.New("切换组织失败，请先添加组织！")
	}
	if unitInfo.Status != 1 {
		return nil, errors.New("切换组织失败，已禁用！")
	}

	result, err = s.CommonUnit.ChangeUnit(moduleName, userId, unitId)
	if err != nil {
		return nil, err
	}

	helper.DelRefreshToken(changeUnitDto.BrancaToken, changeUnitDto.RefreshToken)
	return
}
