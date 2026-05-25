package models_ar

import (
	"WenBeego/apps/common/dto_vo/page_dto"
	"WenBeego/apps/common/dto_vo/queue_dlq_vo"
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/helper"
	"WenBeego/apps/common/models"

	"gorm.io/gorm"
)

type QueueDlqFailedLogAR struct {
	models.QueueDlqFailedLog
}

func (ar *QueueDlqFailedLogAR) Insert(insertData *models.QueueDlqFailedLog) error {
	return global.GetWriteDb().Create(insertData).Error
}

func (ar *QueueDlqFailedLogAR) GetList(reqDto page_dto.QueueDlqListReqDto) (data []queue_dlq_vo.QueueDlqListVo, total int64, err error) {
	data = make([]queue_dlq_vo.QueueDlqListVo, 0)

	query := global.GetReadDb().Model(&models.QueueDlqFailedLog{}).Where("deleted = ?", "0")

	if reqDto.TaskName != "" {
		query = query.Where("task_name LIKE ?", "%"+reqDto.TaskName+"%")
	}
	if reqDto.Status >= 0 {
		query = query.Where("status = ?", reqDto.Status)
	}
	if reqDto.CreateTimeBegin != "" {
		query = query.Where("create_time >= ?", reqDto.CreateTimeBegin)
	}
	if reqDto.CreateTimeEnd != "" {
		reqDto.CreateTimeEnd = helper.TimeAddSeconds(reqDto.CreateTimeEnd, 1)
		query = query.Where("create_time < ?", reqDto.CreateTimeEnd)
	}

	err = query.Count(&total).Error
	if err != nil {
		return
	}
	if reqDto.GetTotal == 1 || total == 0 {
		return
	}

	err = query.Order("create_time DESC").Offset(reqDto.Offset).Limit(reqDto.PageSize).Find(&data).Error
	return
}

func (ar *QueueDlqFailedLogAR) GetPendingListByCondition(reqDto page_dto.QueueDlqRequeueReqDto) (data []models.QueueDlqFailedLog, err error) {
	data = make([]models.QueueDlqFailedLog, 0)

	query := global.GetReadDb().Model(&models.QueueDlqFailedLog{}).
		Where("deleted = ? AND status = ?", "0", page_dto.QUEUE_DLQ_STATUS_PENDING)

	if reqDto.TaskName != "" {
		query = query.Where("task_name LIKE ?", "%"+reqDto.TaskName+"%")
	}
	if reqDto.CreateTimeBegin != "" {
		query = query.Where("create_time >= ?", reqDto.CreateTimeBegin)
	}
	if reqDto.CreateTimeEnd != "" {
		reqDto.CreateTimeEnd = helper.TimeAddSeconds(reqDto.CreateTimeEnd, 1)
		query = query.Where("create_time < ?", reqDto.CreateTimeEnd)
	}

	err = query.Order("create_time ASC").Find(&data).Error
	return
}

func (ar *QueueDlqFailedLogAR) GetPendingListByUUIDs(taskUUIDs []string) (data []models.QueueDlqFailedLog, err error) {
	data = make([]models.QueueDlqFailedLog, 0)

	err = global.GetReadDb().Model(&models.QueueDlqFailedLog{}).
		Where("task_uuid IN ? AND deleted = ? AND status = ?", taskUUIDs, "0", page_dto.QUEUE_DLQ_STATUS_PENDING).
		Order("create_time ASC").
		Find(&data).Error
	return
}

func (ar *QueueDlqFailedLogAR) GetByTaskUUID(taskUUID string) (data *models.QueueDlqFailedLog, err error) {
	data = &models.QueueDlqFailedLog{}
	err = global.GetReadDb().Where("task_uuid = ? AND deleted = ? AND status = ?", taskUUID, "0", page_dto.QUEUE_DLQ_STATUS_PENDING).First(data).Error
	return
}

func (ar *QueueDlqFailedLogAR) UpdateStatus(tx *gorm.DB, taskUUID string, status int) error {
	return global.GetWriteDb().Model(&models.QueueDlqFailedLog{}).Where("task_uuid = ?", taskUUID).Update("status", status).Error
}
