package models

import "time"

type Mchnt struct {
	Id          string    `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	Pid         string    `json:"pid" gorm:"type:bpchar(36);not null;default:'';comment:上级组织id"`
	Logo        string    `json:"logo" gorm:"type:varchar(512);default:'';comment:logo"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null;comment:单位名称"`
	Code        string    `json:"code" gorm:"type:varchar(100);not null;comment:组织机构代码"`
	Corporation string    `json:"corporation" gorm:"type:varchar(100);not null;comment:法人"`
	License     string    `json:"license" gorm:"type:varchar(512);not null;default:'';comment:营业执照"`
	Address     string    `json:"address" gorm:"type:varchar(255);default:'';comment:地址"`
	Status      int       `json:"status" gorm:"type:int4;not null;default:0;comment:0未审核，1审核通过，2审核不通过，3禁用"`
	Deleted     int       `json:"deleted" gorm:"type:int4;default:0;comment:是否删除：0否1是"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;comment:更新时间"`
}

var MCHNT_STATUS_UNREVIEWED = 0
var MCHNT_STATUS_PASSED = 1
var MCHNT_STATUS_UNPASSED = 2
var MCHNT_STATUS_DISABLED = 3
var MCHNT_STATUS_MAP = map[int]string{
	MCHNT_STATUS_UNREVIEWED: "未审核",
	MCHNT_STATUS_PASSED:     "审核通过",
	MCHNT_STATUS_UNPASSED:   "审核不通过",
	MCHNT_STATUS_DISABLED:   "已禁用",
}

func (m *Mchnt) TableName() string {
	return `mchnt`
}
