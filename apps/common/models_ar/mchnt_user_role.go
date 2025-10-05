package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type MchntUserRoleAr struct {
	models.MchntUserRole
}

func (a *MchntUserRoleAr) GetById(id string) (models.MchntUserRole, error) {
	data := models.MchntUserRole{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}
