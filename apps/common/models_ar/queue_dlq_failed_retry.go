package models_ar

import (
	"WenBeego/apps/common/models"

	"gorm.io/gorm"
)

type QueueDlqFailedRetryAR struct {
	models.QueueDlqFailedRetry
}

func (ar *QueueDlqFailedRetryAR) Insert(tx *gorm.DB, insertData *models.QueueDlqFailedRetry) error {
	return tx.Create(insertData).Error
}
