package models

import (
	"time"
)

type PlatRole struct {
	Id        string     `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:角色ID"`
	UnitId    string     `json:"unit_id" gorm:"type:bpchar(36);not null;comment:组织单位id"`
	RoleName  string     `json:"role_name" gorm:"type:varchar(50);not null;comment:角色名称"`
	RoleSort  int        `json:"role_sort" gorm:"type:int4;not null;comment:显示顺序"`
	Status    int        `json:"status" gorm:"type:int4;not null;default:1;comment:角色状态：0停用，1正常"`
	Deleted   int        `json:"deleted" gorm:"type:int4;default:0;comment:删除：0否,1是"`
	CreateBy  string     `json:"create_by" gorm:"type:varchar(64);default:'';comment:创建者"`
	CreatedAt *time.Time `json:"created_at" gorm:"type:timestamp;comment:创建时间"`
	UpdateBy  string     `json:"update_by" gorm:"type:varchar(64);default:'';comment:更新者"`
	UpdateAt  *time.Time `json:"updated_at" gorm:"type:timestamp;comment:更新时间"`
	Remark    string     `json:"remark" gorm:"type:varchar(500);comment:备注"`
}

var PLAT_ROLE_STATUS_NORMAL = 1
var PLAT_ROLE_STATUS_DISABLE = 0
var PLAT_ROLE_STATUS_MAP = map[int]string{
	PLAT_ROLE_STATUS_NORMAL:  "正常",
	PLAT_ROLE_STATUS_DISABLE: "已禁用",
}

func (m *PlatRole) TableName() string {
	return `plat_role`
}
