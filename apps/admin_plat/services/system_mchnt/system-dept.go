package system_mchnt

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/dept_dto"
	"WenBeego/apps/common/dto/page_dto"
	CommonSystem "WenBeego/apps/common/services/system"
)

type DeptService struct {
	commonSystemDept CommonSystem.Dept
}

// 获取组织架构列表
func (s *DeptService) GetUnitDeptList(deptDto page_dto.SystemDeptListReqDto) (resultDto dto.RespDataListDto, err error) {
	return s.commonSystemDept.GetUnitDeptListForAdminPlat(deptDto)
}

// 获取组织架构树
func (s *DeptService) GetUnitDeptTree(baseParamDto dto.BaseParamDto, selectUnitIds []string) (data interface{}, err error) {
	return s.commonSystemDept.GetUnitDeptTreeForAdminPlat(baseParamDto, selectUnitIds)
}

// 获取可用的组织架构负责人
func (s *DeptService) GetUnitDeptPrincipal(baseParamDto dto.BaseParamDto, deptPrincipalDto page_dto.SystemDeptPrincipalReqDto) (data interface{}, err error) {
	return s.commonSystemDept.GetUnitDeptPrincipalForAdminPlat(baseParamDto, deptPrincipalDto)
}

// 保存组织架构
func (s *DeptService) SaveUnitDept(baseParamDto dto.BaseParamDto, deptDto dept_dto.UnitDeptDto) (id string, err error) {
	return s.commonSystemDept.SaveUnitDeptForAdminPlat(baseParamDto, deptDto)
}

// 删除组织架构
func (s *DeptService) Del(baseParamDto dto.BaseParamDto, deptDto dept_dto.UnitDeptDto) (err error) {
	return s.commonSystemDept.DelUnitDeptForAdminPlat(baseParamDto, deptDto)
}
