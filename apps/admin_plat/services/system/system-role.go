package system

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/dto/role_dto"
	CommonSystem "WenBeego/apps/common/services/system"
)

type RoleService struct {
	commonSystemRole CommonSystem.Role
}

// 获取角色列表
func (s *RoleService) GetUnitRoleList(roleDto page_dto.SystemRoleListReqDto) (resultDto dto.RespDataListDto, err error) {
	return s.commonSystemRole.GetUnitRoleList(roleDto)
}

// 保存角色
func (s *RoleService) SaveUnitRole(baseParamDto dto.BaseParamDto, roleDto role_dto.UnitRoleDto) (id string, err error) {
	return s.commonSystemRole.SaveUnitRole(baseParamDto, roleDto)
}

// 删除角色
func (s *RoleService) Del(baseParamDto dto.BaseParamDto, roleDto role_dto.UnitRoleDto) (err error) {
	return s.commonSystemRole.DelUnitRole(baseParamDto, roleDto)
}

func (s *RoleService) GetRoleMenu(baseParamDto dto.BaseParamDto) (dataList interface{}, err error) {
	return s.commonSystemRole.GetRoleMenu(baseParamDto)
}
func (s *RoleService) GetRoleMenuIds(baseParamDto dto.BaseParamDto, ids string) (dataList interface{}, err error) {
	return s.commonSystemRole.GetRoleMenuIds(baseParamDto, ids)
}
