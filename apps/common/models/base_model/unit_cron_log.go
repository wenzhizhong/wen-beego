package base_model

// 定时任务日志
type UnitCronLog struct {
	Id        string `gorm:"column:id;type:bpchar(36);primaryKey;comment:ID" json:"id"`
	CrontabId string `gorm:"column:crontab_id;type:bpchar(36);comment:定时任务id" json:"crontab_id"`
	CreatedAt int64  `gorm:"column:created_at;type:int8;comment:创建时间" json:"created_at"`
	Result    bool   `gorm:"column:result;type:bool;comment:执行结果" json:"result"`
	Remark    string `gorm:"column:remark;type:varchar(2048);comment:备注" json:"remark"`
}
