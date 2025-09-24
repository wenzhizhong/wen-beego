package plat

import (
	"WenBeego/apps/common/ar"
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/services"
	"errors"
)

type Plat struct {
	commonPlat services.CommonPlat
}

// 获取组织信息
func (s *Plat) GetUserUnitList(userId string) (data interface{}, err error) {
	if userId == "" {
		return nil, errors.New("获取组织信息失败，请先登录！")
	}
	dataList, err := ar.GetUserUnitList(userId, &models.Plat{}, &models.PlatUser{})
	data = struct {
		List interface{} `json:"list"`
	}{
		List: dataList,
	}
	return
}
func (s *Plat) ChangeUnit(moduleName string, userId string, changeUnitDto dto.ChangeUnitDto) (result interface{}, err error) {
	unitId := changeUnitDto.Id
	if userId == "" || unitId == "" {
		return nil, errors.New("切换组织失败，请先登录！")
	}

	result, err = s.commonPlat.ChangeUnit(moduleName, userId, unitId)
	if err != nil {
		return nil, err
	}

	helper.DelRefreshToken(changeUnitDto.BrancaToken, changeUnitDto.RefreshToken)
	return
}
