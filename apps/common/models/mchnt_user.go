package models

type MchntUser struct {
	Id      string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	UnitId  string `json:"unit_id" gorm:"type:bpchar(36);not null;comment:组织单位id"`
	UserId  string `json:"user_id" gorm:"type:bpchar(36);not null;comment:员工id"`
	IsAdmin *int32 `json:"is_admin" gorm:"type:int4;default:0;comment:是否管理员;0员工，1超级管理员，2普通管理员"`
	Status  *int32 `json:"status" gorm:"type:int4;default:1;comment:当前组织内状态;0禁用，1正常，2离职"`
	Deleted *int32 `json:"deleted" gorm:"type:int4;default:0;comment:删除：0否,1是"`
}

func (m *MchntUser) TableName() string {
	return `mchnt_user`
}
