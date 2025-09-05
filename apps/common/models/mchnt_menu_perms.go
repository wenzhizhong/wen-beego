package models

import "time"

type MchntMenuPerms struct {
	Id         string    `json:"id" gorm:"type:varchar(36);not null;primaryKey;comment:id"`
	MenuId     string    `json:"menu_id" gorm:"type:varchar(36);not null;comment:资源关联菜单"`
	Type       int16     `json:"type" gorm:"type:int2;not null;comment:类型:1通用，2非通用可选"`
	Name       string    `json:"name" gorm:"type:varchar(64);not null;comment:资源名称"`
	Permission string    `json:"permission" gorm:"type:varchar(255);not null;comment:权限对应的编码"`
	Uri        string    `json:"uri" gorm:"type:varchar(512);not null;comment:资源对应服务器路径"`
	Method     int16     `json:"method" gorm:"type:int2;not null;comment:请求方法 1.GET 2.POST 3.PUT 4.DELETE"`
	Deleted    int16     `json:"deleted" gorm:"type:int2;not null;default:0;comment:是否删除：0否1是"`
	CreatedAt  time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:更新时间"`
}

func (m *MchntMenuPerms) TableName() string {
	return `mchnt_menu_perms`
}
