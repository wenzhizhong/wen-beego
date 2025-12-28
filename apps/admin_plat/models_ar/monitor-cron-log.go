package models_ar

import (
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
)

type PlatCronLogAr struct {
}

func (ar *PlatCronLogAr) GetList(pageReq page_dto.MonitorCronLogListReqDto) (data []models.PlatCronLog, total int64, err error) {
	data = make([]models.PlatCronLog, 0)

	query := global.GetReadDb().Model(&models.PlatCronLog{})
	if pageReq.CrontabId != "" {
		query.Where("crontab_id = ?", pageReq.CrontabId)
	}
	if pageReq.CreatedAt != nil {
		query.Where("created_at = ?", pageReq.CreatedAt.Unix())
	}

	err = query.Count(&total).Error
	if err != nil {
		return
	}
	if total == 0 {
		return
	}
	err = query.Offset(pageReq.Offset).Limit(pageReq.PageSize).Find(&data).Error
	if err != nil {
		return
	}
	return
}
