package menu_dto

import "WenBeego/apps/common/models/base_model"

type MenuDto struct {
	base_model.UnitMenu `gorm:"embedded"`
	AsyncToAll          string `gorm:"-" json:"asyncToAll"`
	MenuFrom            string `json:"menuFrom" gorm:"column:menu_from;comment:系统菜单分类：admin_plat、admin_mchnt；admin_mchnt表示作为商户系统数据管理"`
}
