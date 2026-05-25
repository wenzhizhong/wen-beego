package monitor

import (
	"WenBeego/apps/admin_plat/models_ar"
	dto "WenBeego/apps/common/dto_vo"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/models"
)

type CronLogService struct {
	CronLogModel  models.PlatCronLog
	platCronLogAr models_ar.PlatCronLogAr
}

// 获取定时任务日志列表
func (s *CronLogService) GetCronLogList(reqDto *page_dto.MonitorCronLogListReqDto) (*dto.RespDataListDto, error) {
	data, count, err := s.platCronLogAr.GetList(*reqDto)

	return &dto.RespDataListDto{List: data, Total: count, PageSize: reqDto.PageSize, CurrentPage: reqDto.CurrentPage}, err
}
