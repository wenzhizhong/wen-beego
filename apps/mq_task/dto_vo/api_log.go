package dto_vo

import "WenBeego/apps/common/dto_vo/mq_dto"

type ApiLogDataDto struct {
	ModuleMap map[string]*ApiLogDataUnitDto
}

type ApiLogDataUnitDto struct {
	UnitMap map[string]*ApiLogDataUriDto
}

type ApiLogDataUriDto struct {
	UriMap map[string][]mq_dto.ApiLogDto
}
