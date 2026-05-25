package dept_vo

import "WenBeego/apps/common/models/base_model"

type UnitDeptListVo struct {
	base_model.UnitDept
	Principal string `json:"principal" gorm:"->"`
	Phone     string `json:"phone" gorm:"->"`
	Email     string `json:"email" gorm:"->"`
	UnitName  string `json:"unit_name" gorm:"->"`
}
