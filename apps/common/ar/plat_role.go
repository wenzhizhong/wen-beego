package ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type PlatRoleAr struct {
	models.PlatRole
}

func (a *PlatRoleAr) GetById(id string) (models.PlatRole, error) {
	data := models.PlatRole{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}
