package system

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"WenBeego/apps/common/models_ar/base_ar"
	"errors"
)

type Unit struct {
}

// 系统管理-获取用户列表
func (s *Unit) GetUnitList(unitDto page_dto.SystemUnitListReqDto) (resultDto dto.RespDataListDto, err error) {
	data := make([]base_model.Unit, 0)
	var count int64 = 0

	switch unitDto.ModuleName {
	case "admin_plat":
		data, count, err = base_ar.GetUnitListById(unitDto, &models.Plat{}, &models.PlatUser{})
	case "mchnt_plat":
		data, count, err = base_ar.GetUnitListById(unitDto, &models.Mchnt{}, &models.MchntUser{})
	default:
		err = errors.New("模块名称错误")
	}
	if err != nil {
		return
	}
	resultDto, err = helper.GetRespDataListDto(unitDto.PageSize, unitDto.CurrentPage, count, data)
	return
}
