package system

import (
	"WenBeego/apps/common/dto"
	"WenBeego/apps/common/dto/page_dto"
	CommonSystem "WenBeego/apps/common/services/system"
)

type UnitService struct {
	commonSystemUnit CommonSystem.Unit
}

// 系统管理- 获取用户列表
func (s *UnitService) GetUnitList(unitDto page_dto.SystemUnitListReqDto) (dto.RespDataListDto, error) {
	return s.commonSystemUnit.GetUnitList(unitDto)
}
