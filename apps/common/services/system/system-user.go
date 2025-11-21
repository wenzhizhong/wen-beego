package system

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models_ar/base_ar"
	"errors"
)

type User struct {
}

// 系统管理-获取用户列表
func (s *User) GetUserList(reqDto page_dto.SystemUserListReqDto) (resultDto dto.RespDataListDto, err error) {
	data := make([]page_dto.SystemUserListDto, 0)
	var count int64 = 0

	switch reqDto.ModuleName {
	case "admin_plat":
		data, count, err = base_ar.GetUserListOfUnitById(reqDto, &models.PlatUser{}, &models.PlatUserProfile{}, &models.PlatDept{}, &models.PlattUserDept{}, &models.PlatRole{}, &models.PlatUserRole{})
	case "mchnt_plat":
		data, count, err = base_ar.GetUserListOfUnitById(reqDto, &models.MchntUser{}, &models.MchntUserProfile{}, &models.MchntDept{}, &models.MchntUserDept{}, &models.MchntRole{}, &models.MchntUserRole{})
	default:
		err = errors.New("GetUserList：模块名称错误")
	}
	if err != nil {
		return
	}

	for k, v := range data {
		data[k].Avatar, _ = helper.LocalFileSign(reqDto.Host, v.Avatar)
	}
	resultDto, err = helper.GetRespDataListDto(reqDto.PageSize, reqDto.CurrentPage, count, data)
	return
}
