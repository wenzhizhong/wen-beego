package models

import (
	"WenBeego/apps/common/global"
	"errors"
	"time"
)

type PlatMenu struct {
	Id        string    `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	UnitId    string    `json:"unit_id" gorm:"type:bpchar(36);not null;comment:plat/mchut id"`
	Icon      string    `json:"icon" gorm:"type:varchar(50);default:'';comment:图标"`
	Name      string    `json:"name" gorm:"type:varchar(20);default:'';comment:菜单名称"`
	ApiPath   string    `json:"api_path" gorm:"type:varchar(255);default:'';comment:API路径"`
	PagePath  string    `json:"page_path" gorm:"type:varchar(255);default:'';comment:页面路径"`
	Type      *int16    `json:"type" gorm:"type:int2;default:1;comment:类型"`
	Pid       string    `json:"pid" gorm:"type:varchar(36);default:'';comment:父级ID"`
	AllPid    string    `json:"all_pid" gorm:"type:varchar(1000);default:'';comment:所有父级ID"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;default:CURRENT_TIMESTAMP;comment:创建时间"`
	Weight    *int32    `json:"weight" gorm:"type:int4;default:0;comment:权重"`
	Visible   int32     `json:"visible" gorm:"type:int4;not null;default:1;comment:是否可见"`
	Deleted   *int32    `json:"deleted" gorm:"type:int4;default:0;comment:是否删除"`
}

func (m *PlatMenu) TableName() string {
	return `plat_menu`
}

func (m *PlatMenu) GetById(id string) (PlatMenu, error) {
	menu := PlatMenu{}
	if id == "" {
		return menu, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&menu)
	return menu, result.Error
}
