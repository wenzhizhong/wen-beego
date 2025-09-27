package models

type MchntUser struct {
	Id            string `json:"id" gorm:"type:bpchar(36);not null;primaryKey;comment:ID"`
	UnitId        string `json:"unit_id" gorm:"type:bpchar(36);not null;comment:组织单位id"`
	IsDefault     int    `json:"is_default" gorm:"type:int4;default:0;comment:是否默认：0否1是"`
	DefaultUnitId string `json:"default_unit_id" gorm:"-"`
	UserId        string `json:"user_id" gorm:"type:bpchar(36);not null;comment:员工id"`
	Deleted       int    `json:"deleted" gorm:"type:int4;default:0;comment:删除：0否,1是"`
	Phone         string `json:"phone" gorm:"not null;unique;size:11;comment:手机号"`
	Name          string `json:"name" gorm:"not null;size:20;comment:姓名"`
}

var MCHNT_IS_DEFAULT_NO = 0
var MCHNT_IS_DEFAULT_YES = 1

func (m *MchntUser) TableName() string {
	return `mchnt_user`
}
