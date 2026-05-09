package models

import "time"

type QueueDlqFailedLog struct {
	TaskUUID   string    `json:"task_uuid" gorm:"column:task_uuid;type:varchar(64);primaryKey;not null;comment:任务UUID"`
	TaskName   string    `json:"task_name" gorm:"column:task_name;type:varchar(128);comment:任务名称"`
	TaskArgs   string    `json:"task_args" gorm:"column:task_args;type:text;comment:任务参数JSON"`
	ErrorMsg   string    `json:"error_msg" gorm:"column:error_msg;type:text;comment:错误信息"`
	CreateTime time.Time `json:"create_time" gorm:"column:create_time;type:timestamp;default:now();comment:创建时间"`
	Status     int       `json:"status" gorm:"column:status;type:int4;default:0;comment:状态"`
	Deleted    string    `json:"deleted" gorm:"column:deleted;type:varchar;default:0;comment:软删除"`
}

func (m *QueueDlqFailedLog) TableName() string {
	return "queue_dlq_failed_log"
}
