package base_model

import "time"

type UnitCron struct {
	Id        string     `gorm:"column:id;type:bpchar(36);primaryKey;comment:ID" json:"id"`
	UnitId    string     `gorm:"column:unit_id;type:bpchar(36);comment:组织单位id" json:"unit_id"`
	Name      string     `gorm:"column:name;type:varchar(100);comment:任务名称" json:"name"`
	NameEn    string     `gorm:"column:name_en;type:varchar(100);comment: 任务名称-英文" json:"name_en"`
	Group     string     `gorm:"column:group;type:varchar(100);default:'default';comment:分组名称" json:"group"`
	CronExpr  string     `gorm:"column:cron_expr;type:varchar(2048);comment:cron表达式" json:"cron_expr"`
	Status    int        `gorm:"column:status;type:int;default:0;comment:状态：0禁用1启用" json:"status"`
	CreatedBy string     `gorm:"column:created_by;type:bpchar(36);comment:创建人" json:"created_by"`
	CreatedAt *time.Time `gorm:"column:created_at;type:date;comment:创建时间" json:"created_at"`
	UpdatedBy *string    `gorm:"column:updated_by;type:varchar(36);comment:更新人" json:"updated_by"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:date;comment:更新时间" json:"updated_at"`
	Deleted   int        `gorm:"column:deleted;type:int;comment:是否删除：0否1是" json:"deleted"`
	Remark    string     `gorm:"column:remark;type:varchar(100);comment:备注" json:"remark"`

	CreatedByName string `gorm:"->" json:"created_by_name"`
	UpdatedByName string `gorm:"->" json:"updated_by_name"`
}

var UNIt_CRONTAB_STATUS_ENABLED = 1
var UNIt_CRONTAB_STATUS_DISABLED = 0
var UNIt_CRONTAB_STATUS_MAP = map[int]string{
	UNIt_CRONTAB_STATUS_ENABLED:  "启用",
	UNIt_CRONTAB_STATUS_DISABLED: "禁用",
}
