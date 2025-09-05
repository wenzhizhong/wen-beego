package models

import (
	"time"
)

type MchntRole struct {
	Id        string     `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:角色ID"`
	UnitId    string     `json:"unit_id" gorm:"type:bpchar(36);not null;comment:组织单位id"`
	RoleName  string     `json:"role_name" gorm:"type:varchar(50);not null;comment:角色名称"`
	RoleSort  int32      `json:"role_sort" gorm:"type:int4;not null;comment:显示顺序"`
	Status    int32      `json:"status" gorm:"type:int4;not null;default:1;comment:角色状态：0停用，1正常"`
	Deleted   *int32     `json:"deleted" gorm:"type:int4;default:0;comment:删除：0否,1是"`
	CreatedBy string     `json:"created_by" gorm:"type:varchar(64);default:'';comment:创建者"`
	CreatedAt *time.Time `json:"created_at" gorm:"type:timestamp;comment:创建时间"`
	UpdatedBy string     `json:"updated_by" gorm:"type:varchar(64);default:'';comment:更新者"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"type:timestamp;comment:更新时间"`
	Remark    *string    `json:"remark" gorm:"type:varchar(500);comment:备注"`
}

func (m *MchntRole) TableName() string {
	return `mchnt_role`
}
