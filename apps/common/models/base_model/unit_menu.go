package base_model

import (
	"time"
)

type UnitMenu struct {
	Id        string    `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	UnitId    string    `json:"unit_id" gorm:"type:bpchar(36);not null;comment:plat/mchut id"`
	Icon      string    `json:"icon" gorm:"type:varchar(50);default:'';comment:图标"`
	Name      string    `json:"name" gorm:"type:varchar(20);default:'';comment:菜单名称"`
	ApiPath   string    `json:"api_path" gorm:"type:varchar(255);default:'';comment:API路径"`
	PagePath  string    `json:"page_path" gorm:"type:varchar(255);default:'';comment:页面路径"`
	Type      int       `json:"type" gorm:"type:int2;default:0;comment:类型"` // 菜单类型（0代表菜单、1代表iframe、2代表外链、3代表按钮）
	Pid       string    `json:"pid" gorm:"type:varchar(36);default:'';comment:父级ID"`
	AllPid    string    `json:"all_pid" gorm:"type:varchar(1000);default:'';comment:所有父级ID"`
	CreatedAt time.Time `json:"created_at" gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	Weight    int       `json:"weight" gorm:"type:int4;default:0;comment:权重"`
	Visible   int       `json:"visible" gorm:"type:int4;not null;default:1;comment:是否可见"`
	Deleted   int       `json:"deleted" gorm:"type:int4;default:0;comment:是否删除"`
	Clone     int       `json:"clone" gorm:"type:int2;default:0;comment:允许克隆：0否1是"`
}
