package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
)

type QueueDlxFailedLogAR struct {
	models.QueueDlxFailedLog
}

func (ar *QueueDlxFailedLogAR) Insert(insertData *models.QueueDlxFailedLog) error {
	return global.GetWriteDb().Create(insertData).Error
}
