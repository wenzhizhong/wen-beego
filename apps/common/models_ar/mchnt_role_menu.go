package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type MchntRoleMenuAr struct {
	models.MchntRoleMenu
}

func (a *MchntRoleMenuAr) GetById(id string) (models.MchntRoleMenu, error) {
	data := models.MchntRoleMenu{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}
