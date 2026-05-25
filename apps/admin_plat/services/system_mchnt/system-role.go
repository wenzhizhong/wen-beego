package system_mchnt

import (
	dto "WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/dto_vo/role_dto"
	CommonSystem "WenBeego/apps/common/services/system"
)

type RoleService struct {
	commonSystemRole CommonSystem.Role
}

// 获取角色列表
func (s *RoleService) GetUnitRoleList(roleDto page_dto.SystemRoleListReqDto) (resultDto dto.RespDataListDto, err error) {
	return s.commonSystemRole.GetUnitRoleListForAdminPlat(roleDto)
}

// 保存角色
func (s *RoleService) SaveUnitRole(baseParamDto dto.BaseParamDto, roleDto role_dto.UnitRoleDto) (id string, err error) {
	return s.commonSystemRole.SaveUnitRoleForAdminPlat(baseParamDto, roleDto)
}

// 删除角色
func (s *RoleService) Del(baseParamDto dto.BaseParamDto, roleDto role_dto.UnitRoleDto) (err error) {
	return s.commonSystemRole.DelUnitRoleForAdminPlat(baseParamDto, roleDto)
}

// 获取角色可选菜单
func (s *RoleService) GetRoleMenu(baseParamDto dto.BaseParamDto, selectUnitIds []string) (dataList interface{}, err error) {
	return s.commonSystemRole.GetRoleMenuForAdminPlat(baseParamDto, selectUnitIds)
}

// 获取角色已选菜单
func (s *RoleService) GetRoleMenuIds(baseParamDto dto.BaseParamDto, roleId string) (dataList interface{}, err error) {
	return s.commonSystemRole.GetRoleMenuIdsForAdminPlat(baseParamDto, roleId)
}

// 保存角色所选菜单
func (s *RoleService) RoleMenuSave(baseParamDto dto.BaseParamDto, roleMenuSaveDto role_dto.RoleMenuSaveDto) (err error) {
	return s.commonSystemRole.RoleMenuSaveForAdminPlat(baseParamDto, roleMenuSaveDto)
}
