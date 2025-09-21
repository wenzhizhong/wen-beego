package ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type PlatUserRoleAr struct {
	models.PlatUserRole
}

func (a *PlatUserRoleAr) GetById(id string) (models.PlatUserRole, error) {
	data := models.PlatUserRole{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}
