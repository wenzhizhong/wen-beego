package auth

import (
	"WenBeego/apps/common/dto/auth_dto"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models_ar/base_ar"
	"WenBeego/apps/common/services/auth"
	"errors"
)

type Mchnt struct {
	CommonUnit auth.CommonUnit
}

// 获取组织信息
func (s *Mchnt) GetUserUnitList(userId, host string) (data interface{}, err error) {
	if userId == "" {
		return nil, errors.New("获取组织信息失败，请先登录！")
	}
	dataList, err := base_ar.GetUserUnitList(userId, &models.Mchnt{}, &models.MchntUser{})
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
func (s *Mchnt) ChangeUnit(moduleName string, userId string, changeUnitDto auth_dto.ChangeUnitDto) (result interface{}, err error) {
	unitId := changeUnitDto.Id
	if userId == "" || unitId == "" {
		return nil, errors.New("切换组织失败，请先登录！")
	}

	unitInfo, err := base_ar.GetUserUnitById(userId, unitId, &models.Mchnt{}, &models.MchntUser{})
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
