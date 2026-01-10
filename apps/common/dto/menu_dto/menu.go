package menu_dto

import "WenBeego/apps/common/models/base_model"

type MenuDto struct {
	base_model.UnitMenu
	AsyncToAll string `json:"asyncToAll"`
}
