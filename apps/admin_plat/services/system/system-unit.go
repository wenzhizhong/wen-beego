package system

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/dto/unit_dto"
	CommonSystem "WenBeego/apps/common/services/system"
)

type UnitService struct {
	commonSystemUnit CommonSystem.Unit
}

// 系统管理- 获取用户列表
func (s *UnitService) GetUnitList(unitDto page_dto.SystemUnitListReqDto) (dto.RespDataListDto, error) {
	return s.commonSystemUnit.GetUnitList(unitDto)
}

// func (s *UnitService) GetUnitTree() (dto.RespDataListDto, error) {
// 	return s.commonSystemUnit.GetUnitTree()
// }

// func (s *UnitService) GetUnitInfo(unitId string) (dto.RespDataDto, error) {
// 	return s.commonSystemUnit.GetUnitInfo(unitId)
// }

func (s *UnitService) Save(baseParamDto dto.BaseParamDto, unitDto unit_dto.UnitDto) (map[string]string, error) {
	return s.commonSystemUnit.Save(baseParamDto, unitDto)
}
