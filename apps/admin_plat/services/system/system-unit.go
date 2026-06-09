package system

import (
	dto "WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/dto_vo/unit_dto"
	CommonSystem "WenBeego/apps/common/services/system"
)

type UnitService struct {
	commonSystemUnit CommonSystem.Unit
}

// 系统管理- 获取内部组织列表
func (s *UnitService) GetUnitList(unitDto page_dto.SystemUnitListReqDto) (dto.RespDataListDto, error) {
	return s.commonSystemUnit.GetUnitList(unitDto)
}

// 系统管理-内部组织保存
func (s *UnitService) Save(baseParamDto dto.BaseParamDto, unitDto unit_dto.UnitDto) (map[string]string, error) {
	return s.commonSystemUnit.Save(baseParamDto, unitDto)
}

// 系统管理- 保存
func (s *UnitService) Del(baseParamDto dto.BaseParamDto, unitDto unit_dto.UnitDto) error {
	// TODO: 删除组织单位前置条件

	// 删除组织单位
	return s.commonSystemUnit.Del(baseParamDto, unitDto)
}

// 系统管理- 变更组织状态
func (s *UnitService) ChangeStatus(baseParamDto dto.BaseParamDto, unitDto unit_dto.UnitDto) error {
	return s.commonSystemUnit.ChangeStatus(baseParamDto, unitDto)
}
