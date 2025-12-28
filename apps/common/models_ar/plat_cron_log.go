package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
	"WenBeego/apps/common/models/base_model"
	"fmt"
	"strings"
)

type PlatCronLogAr struct {
}

// 插入定时任务日志
func (ar *PlatCronLogAr) Insert(cronNameEn string, result bool, remark string) (err error) {
	if cronNameEn == "" {
		return
	}

	// 获取定时任务信息
	crontabInfo := &models.PlatCron{}
	err = global.GetReadDb().Model(&models.PlatCron{}).
		Select("id").
		Where("name_en = ?", cronNameEn).
		Find(crontabInfo).Error
	if err != nil {
		return
	}
	if crontabInfo.GetID() == "" {
		err = fmt.Errorf("定时任务[%s]不存在", cronNameEn)
		return
	}

	// 新增
	uuid, err := helper.GetUuid()
	if err != nil {
		return err
	}
	if remark != "" {
		tmpPath := strings.ReplaceAll(global.RootPath, "\\", "/")
		remark = strings.ReplaceAll(remark, tmpPath, "")
		remark = remark[:512]
	}

	createdAt := helper.GetTimestamp()
	platCronLog := base_model.UnitCronLog{
		Id:        uuid,
		CrontabId: crontabInfo.GetID(),
		Result:    result,
		Remark:    remark,
		CreatedAt: createdAt,
	}

	return global.GetWriteDb().Model(&models.PlatCronLog{}).Create(&platCronLog).Error
}
