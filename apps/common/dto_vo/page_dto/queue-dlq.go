package page_dto

import "WenBeego/apps/common/dto_vo"

var QUEUE_DLQ_STATUS_PENDING = 0
var QUEUE_DLQ_STATUS_REQUEUED = 1
var QUEUE_DLQ_STATUS_MAP = map[int]string{
	QUEUE_DLQ_STATUS_PENDING:  "待处理",
	QUEUE_DLQ_STATUS_REQUEUED: "已重新入队",
}

type QueueDlqListReqDto struct {
	dto_vo.BaseParamDto
	dto_vo.ReqDataListDto
	TaskName        string
	Status          int
	CreateTimeBegin string
	CreateTimeEnd   string
	GetTotal        int
}

type QueueDlqRequeueReqDto struct {
	TaskName        string `json:"task_name"`
	CreateTimeBegin string `json:"create_time_begin"`
	CreateTimeEnd   string `json:"create_time_end"`
}
