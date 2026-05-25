package role_vo

import "WenBeego/apps/common/models/base_model"

type UnitRoleListVo struct {
	base_model.UnitRole
	UnitName         string `json:"unit_name" gorm:"->"`
	RoleClassifyName string `json:"role_classify_name" gorm:"->"`
}
