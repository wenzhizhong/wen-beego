package system

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
	return s.commonSystemDept.GetUnitDeptList(deptDto)
}

func (s *DeptService) SaveUnitDept(baseParamDto dto.BaseParamDto, deptDto dept_dto.UnitDeptDto) (id string, err error) {
	return s.commonSystemDept.SaveUnitDept(baseParamDto, deptDto)
}

func (s *DeptService) Del(baseParamDto dto.BaseParamDto, deptDto dept_dto.UnitDeptDto) (err error) {
	return s.commonSystemDept.DelUnitDept(baseParamDto, deptDto)
}
