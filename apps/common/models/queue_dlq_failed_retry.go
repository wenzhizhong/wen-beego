package models

type QueueDlqFailedRetry struct {
	Id    string `json:"id" gorm:"column:id;type:varchar(64);not null;comment:原任务UUID"`
	NewId string `json:"new_id" gorm:"column:new_id;type:varchar(64);not null;comment:新任务UUID"`
}

func (m *QueueDlqFailedRetry) TableName() string {
	return "queue_dlq_failed_retry"
}
