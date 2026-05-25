package queue_dlq_dto

type RequeueDto struct {
	TaskUUID        string   `json:"task_uuid"`
	TaskUUIDs       []string `json:"task_uuids"`
	TaskName        string   `json:"task_name"`
	CreateTimeBegin string   `json:"create_time_begin"`
	CreateTimeEnd   string   `json:"create_time_end"`
	GetTotal        int      `json:"getTotal"`
}
