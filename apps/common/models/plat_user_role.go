package models

type PlatUserRole struct {
	Id      string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	UnitId  string `json:"unit_id" gorm:"type:bpchar(36);not null;comment:组织单位id"`
	UserId  string `json:"user_id" gorm:"type:bpchar(36);not null;comment:员工id"`
	RoleId  string `json:"role_id" gorm:"type:bpchar(36);not null;comment:员工角色id"`
	Deleted int64  `json:"deleted" gorm:"type:int4;default:0;comment:删除：0否,1是"`
}

func (m *PlatUserRole) TableName() string {
	return `plat_user_role`
}
