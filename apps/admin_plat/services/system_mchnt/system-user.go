package system_mchnt

import (
	dto "WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/dto_vo/user_dto"
	CommonSystem "WenBeego/apps/common/services/system"
)

type UserService struct {
	commonSystemUser CommonSystem.User
}

// 系统管理- 获取用户列表
func (s *UserService) GetUserList(reqDto page_dto.SystemUserListReqDto) (dto.RespDataListDto, error) {
	return s.commonSystemUser.GetUserListForAdminPlat(reqDto)
}

// 系统管理- 获取可选角色树
func (s *UserService) GetUnitRoleTree(baseParamDto dto.BaseParamDto, selectUnitIds []string) (interface{}, error) {
	return s.commonSystemUser.GetUnitRoleTreeForAdminPlat(baseParamDto, selectUnitIds)
}

// SaveUser
func (s *UserService) SaveUser(baseParamDto dto.BaseParamDto, unitUserSaveDto *user_dto.UnitUserSaveDto) (interface{}, error) {
	return s.commonSystemUser.SaveUserForAdminPlat(baseParamDto, unitUserSaveDto)
}

func (s *UserService) DelUnitUser(baseParamDto dto.BaseParamDto, ids []string) error {
	return s.commonSystemUser.DelUnitUserForAdminPlat(baseParamDto, ids)
}
