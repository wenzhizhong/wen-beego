package base_model

import "time"

type GenerateCode struct {
	Id         string     `gorm:"column:id;type:bpchar(36);primaryKey;comment:ID" json:"id"`
	TableName  string     `gorm:"column:table_name;type:varchar(100);comment:表名" json:"table_name"`
	Data       string     `gorm:"column:data;type:text;comment:数据" json:"data"`
	CreateTime *time.Time `gorm:"column:create_time;type:timestamptz;comment:创建时间" json:"create_time"`
	Deleted    int        `gorm:"column:deleted;type:int2;default:0;comment:是否删除" json:"deleted"`
}
