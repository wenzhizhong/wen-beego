package models_ar

import (
	"WenBeego/apps/common/global"
	"WenBeego/apps/common/models"
	"errors"
)

type MchntUserAr struct {
	models.MchntUser
}

func (a *MchntUserAr) GetById(id string) (models.MchntUser, error) {
	data := models.MchntUser{}
	if id == "" {
		return data, errors.New("id不能为空")
	}
	result := global.GetReadDb().Where("id = ?", id).Take(&data)
	return data, result.Error
}
