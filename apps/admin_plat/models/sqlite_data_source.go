package models

import "time"

type SqliteDataSource struct {
	Id          string     `gorm:"column:id;type:bpchar(36);primaryKey;comment:ID" json:"id"`
	Name        string     `gorm:"column:name;type:varchar(100);not null;comment:数据源名称" json:"name"`
	FilePath    string     `gorm:"column:file_path;type:varchar(500);not null;comment:数据文件路径" json:"file_path"`
	Password    string     `gorm:"column:password;type:varchar(500);comment:加密后的连接密码" json:"-"`
	Description string     `gorm:"column:description;type:varchar(500);comment:描述" json:"description"`
	CreateTime  *time.Time `gorm:"column:create_time;type:timestamptz(6);comment:创建时间" json:"create_time"`
	UpdateTime  *time.Time `gorm:"column:update_time;type:timestamptz(6);comment:更新时间" json:"update_time"`
	Deleted     int        `gorm:"column:deleted;type:int2;default:0;comment:是否删除" json:"-"`
}

func (m *SqliteDataSource) TableName() string {
	return `tools_sqlite_data_source`
}
