package models_ar

import (
	"WenBeego/apps/common/dto/page_dto"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
)

type PlatCronLogAr struct {
}

// 插入定时任务日志
func (ar *PlatCronLogAr) Insert(cronId string, result bool, remark string) (err error) {
	uuid, err := helper.GetUuid()
	if err != nil {
		return err
	}
	createdAt := helper.GetTimestamp()
	platCronLog := base_model.UnitCronLog{
		Id:        uuid,
		CrontabId: cronId,
		Result:    result,
		Remark:    remark,
		CreatedAt: createdAt,
	}

	return global.GetWriteDb().Model(&models.PlatCronLog{}).Create(&platCronLog).Error
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
