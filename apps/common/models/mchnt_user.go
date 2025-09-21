package models

type MchntUser struct {
	Id            string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	UnitId        string `json:"unit_id" gorm:"type:bpchar(36);not null;comment:组织单位id"`
	IsDefault     int    `json:"is_default" gorm:"type:int4;default:0;comment:是否默认：0否1是"`
	DefaultUnitId string `json:"default_unit_id" gorm:"not null; comment:默认组织ID"`
	UserId        string `json:"user_id" gorm:"type:bpchar(36);not null;comment:员工id"`
	Status        int    `json:"status" gorm:"type:int4;default:1;comment:当前组织内状态;0禁用，1正常，2离职"`
	Deleted       int    `json:"deleted" gorm:"type:int4;default:0;comment:删除：0否,1是"`
}

func (m *MchntUser) TableName() string {
	return `mchnt_user`
}
