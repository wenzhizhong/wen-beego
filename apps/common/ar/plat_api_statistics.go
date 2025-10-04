package ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"
)

type PlatApiStatisticsAr struct {
	models.PlatApiStatistics
}

// 插入
func (m *PlatApiStatisticsAr) Insert(data *models.PlatApiStatistics) error {
	if data.ID == "" {
		uuid, err := helper.GetUuid()
		data.ID = uuid
		if err != nil {
			return err
		}
	}
	result := global.GetWriteDb().Create(&data)
	return result.Error
}

// 批量插入
func (m *PlatApiStatisticsAr) InsertBatch(data []*models.PlatApiStatistics) error {
	return global.GetWriteDb().Create(&data).Error
}

// 获取今日数据
func (m *PlatApiStatisticsAr) GetTodayData() (data []*models.PlatApiStatistics, err error) {
	err = global.GetReadDb().
		Model(&models.PlatApiStatistics{}).
		Where("date = ?", helper.GetDateStamp()).Find(&data).Error
	return
}
