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

type User struct {
}

// 系统管理-获取用户列表
func (s *User) GetUserList(reqDto page_dto.SystemUserListReqDto) (resultDto dto.RespDataListDto, err error) {
	data := make([]base_model.UnitUser, 0)
	var count int64 = 0

	switch reqDto.ModuleName {
	case "admin_plat":
		data, count, err = base_ar.GetUserListOfUnitById[*models.PlatUser](reqDto)
	case "mchnt_plat":
		data, count, err = base_ar.GetUserListOfUnitById[*models.MchntUser](reqDto)
	default:
		err = errors.New("GetUserList：模块名称错误")
	}
	if err != nil {
		return
	}
	resultDto, err = helper.GetRespDataListDto(reqDto.PageSize, reqDto.CurrentPage, count, data)
	return
}
