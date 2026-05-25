package unit_dto

import "WenBeego/apps/common/models/base_model"

type UnitUserDto struct {
	base_model.UnitUser
}
type UnitUserAllDataDto struct {
	base_model.UnitUser
	base_model.UnitUserProfile
	Deleted int    `json:"deleted" gorm:"not null;default:0;comment:是否删除"`
	Id      string `json:"id" gorm:"type:bpchar(36);primaryKey;comment:ID"`
}
