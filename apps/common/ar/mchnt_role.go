package ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type MchntRoleAr struct {
	models.MchntRole
}

func (a *MchntRoleAr) GetById(id string) (models.MchntRole, error) {
	data := models.MchntRole{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}
