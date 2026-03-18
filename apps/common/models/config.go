package models

import "time"

type Config struct {
	Id          string    `json:"id" gorm:"not null;primaryKey;comment:ID"`
	Name        string    `json:"name" gorm:"not null;unique;size:255;comment:配置名称"`
	Value       string    `json:"value" gorm:"type:text;comment:配置值"`
	Description string    `json:"description" gorm:"type:text;comment:备注"`
	ValueType   string    `json:"value_type" gorm:"size:20;default:string;comment:数据类型"`
	Category    string    `json:"category" gorm:"size:50;comment:分类"`
	IsReadonly  bool      `json:"is_readonly" gorm:"default:false;comment:是否可修改"`
	Version     int32     `json:"version" gorm:"default:1;comment:版本"`
	CreatedAt   time.Time `json:"created_at" gorm:"default:now();comment:创建时间"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"default:now();comment:更新时间"`
	CreatedBy   string    `json:"created_by" gorm:"size:50;comment:创建者"`
	UpdatedBy   string    `json:"updated_by" gorm:"size:50;comment:更新人"`
	Deleted     bool      `json:"deleted" gorm:"default:false;comment:是否删除"`
}

func (m *Config) TableName() string {
	return `config`
}
