package system

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/page_dto"
	CommonSystem "WenBeego/apps/common/services/system"
)

type UserService struct {
	commonSystemUser CommonSystem.User
}

// 系统管理- 获取用户列表
func (s *UserService) GetUserList(reqDto page_dto.SystemUserListReqDto) (dto.RespDataListDto, error) {
	return s.commonSystemUser.GetUserList(reqDto)
}
