package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
)

type PlatCronAr struct {
}

// 启动项目恢复定时任务-获取计划任务列表
func (ar *PlatCronAr) RunProjectGetCronList() (data []models.PlatCron, err error) {
	data = make([]models.PlatCron, 0)

	platCronModel := &models.PlatCron{}
	tablePlatCronName := platCronModel.TableName()
	query := global.GetReadDb().
		Model(platCronModel).
		Where("status = 1 AND deleted = 0")

	err = query.Select(tablePlatCronName + ".*, '' AS created_by_name, '' AS updated_by_name").
		Find(&data).Error

	return
}
