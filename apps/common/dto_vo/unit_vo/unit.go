package unit_vo

import "WenBeego/apps/common/models/base_model"

type UnitListVo struct {
	base_model.Unit
	DefaultUnitId     string `json:"default_unit_id" gorm:"->"`
	DefaultUnitUserId string `json:"default_unit_user_id" gorm:"->"`
}
