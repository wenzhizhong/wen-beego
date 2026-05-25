package models_ar

import (
	"WenBeego/apps/common/dto_vo/cron_vo"
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
)

type PlatCronLogAr struct {
}

func (ar *PlatCronLogAr) GetList(pageReq page_dto.MonitorCronLogListReqDto) (data []cron_vo.UnitCronLogListVo, total int64, err error) {
	data = make([]cron_vo.UnitCronLogListVo, 0)

	query := global.GetReadDb().Model(&models.PlatCronLog{})
	if pageReq.NameEn != "" {
		query = query.Where("name_en = ?", pageReq.NameEn)
	}
	if pageReq.CreatedAt != nil {
		tmpBeginTime := pageReq.CreatedAt.Format("2006-01-02") + " 00:00:00"
		tmpEndTime := pageReq.CreatedAt.Format("2006-01-02") + " 23:59:59"
		beginTime := helper.GetTimestamp(tmpBeginTime)
		endTime := helper.GetTimestamp(tmpEndTime)
		query = query.Where("created_at BETWEEN ? AND ?", beginTime, endTime)
	}

	err = query.Count(&total).Error
	if err != nil {
		return
	}
	if total == 0 {
		return
	}
	err = query.Order("created_at desc").
		Offset(pageReq.Offset).
		Limit(pageReq.PageSize).
		Find(&data).Error
	if err != nil {
		return
	}
	return
}
